package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/thalisonh/20-CleanArch/configs"
	"github.com/thalisonh/20-CleanArch/internal/entity"
	"github.com/thalisonh/20-CleanArch/internal/infra/graph"
	"github.com/thalisonh/20-CleanArch/internal/infra/grpc/pb"
	"github.com/thalisonh/20-CleanArch/internal/infra/grpc/service"
	"github.com/thalisonh/20-CleanArch/internal/infra/web/webserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configs.DBUser,
		configs.DBPassword,
		configs.DBHost,
		configs.DBPort,
		configs.DBName,
	)
	// dsn := "root:root@tcp(mysql:3306)/orders?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println(dsn)

	log.Printf("\nConnecting to MYSQL database...")

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	database.AutoMigrate(&entity.Order{})

	createOrderUseCase := NewCreateOrderUseCase(database)
	listOrderUseCase := NewListOrderUseCase(database)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(database)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/order/list", webOrderHandler.FindAll)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}
