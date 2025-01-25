//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/thalisonh/20-CleanArch/internal/entity"
	"github.com/thalisonh/20-CleanArch/internal/event"
	"github.com/thalisonh/20-CleanArch/internal/infra/database"
	"github.com/thalisonh/20-CleanArch/internal/infra/web"
	"github.com/thalisonh/20-CleanArch/internal/usecase"
	"github.com/thalisonh/20-CleanArch/pkg/events"
	"gorm.io/gorm"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *gorm.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
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

func NewWebOrderHandler(db *gorm.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
