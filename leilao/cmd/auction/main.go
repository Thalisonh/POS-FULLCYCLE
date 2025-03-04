package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thalisonh/auction/configuration/database/mongodb"
	"github.com/thalisonh/auction/internal/infra/api/web/controller/auction_controller"
	"github.com/thalisonh/auction/internal/infra/api/web/controller/bid_controller"
	"github.com/thalisonh/auction/internal/infra/api/web/controller/user_controller"
	auctionRepository "github.com/thalisonh/auction/internal/infra/database/auction"
	bidRepository "github.com/thalisonh/auction/internal/infra/database/bid"
	userRepository "github.com/thalisonh/auction/internal/infra/database/user"
	"github.com/thalisonh/auction/internal/usecase/auction_usecase"
	"github.com/thalisonh/auction/internal/usecase/bid_usecase"
	"github.com/thalisonh/auction/internal/usecase/user_usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("error trying to load env variables")
		return
	}

	database, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())

		return
	}

	router := gin.Default()

	userController, bidController, auctionController := initDependencies(database)

	router.GET("/auctions", auctionController.FindAuctions)
	router.GET("/auctions/:auctionId", auctionController.FindAuctionById)
	router.POST("/auctions", auctionController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionById)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController,
) {
	auctionRepository := auctionRepository.NewAuctionRepository(database)
	userRepository := userRepository.NewUserRepository(database)
	bidRepository := bidRepository.NewBidRepository(database, auctionRepository)

	userUseCase := user_usecase.NewUserUseCase(userRepository)
	auctionUseCase := auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository)
	bidUseCase := bid_usecase.NewBidUsecase(bidRepository)

	bidController = bid_controller.NewBidController(bidUseCase)
	userController = user_controller.NewUserController(userUseCase)
	auctionController = auction_controller.NewAuctionController(auctionUseCase)

	return
}
