//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/thalisonh/20-CleanArch/internal/entity"
	"github.com/thalisonh/20-CleanArch/internal/infra/database"
	"github.com/thalisonh/20-CleanArch/internal/infra/web"
	"github.com/thalisonh/20-CleanArch/internal/usecase"
	"gorm.io/gorm"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

func NewCreateOrderUseCase(db *gorm.DB) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrderUseCase(db *gorm.DB) *usecase.ListOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrderUseCase,
	)

	return &usecase.ListOrderUseCase{}
}

func NewWebOrderHandler(db *gorm.DB) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
