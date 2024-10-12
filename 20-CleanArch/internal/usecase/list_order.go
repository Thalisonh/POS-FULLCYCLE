package usecase

import (
	"github.com/thalisonh/20-CleanArch/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrderUseCase) Execute() ([]OrderOutputDTO, error) {
	response := []OrderOutputDTO{}
	orders, err := l.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, item := range orders {
		response = append(response, OrderOutputDTO{
			ID:         item.ID,
			Price:      item.Price,
			Tax:        item.Tax,
			FinalPrice: item.FinalPrice,
		})
	}

	return response, nil
}
