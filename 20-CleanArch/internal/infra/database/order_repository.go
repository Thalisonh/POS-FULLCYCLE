package database

import (
	"github.com/thalisonh/20-CleanArch/internal/entity"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	err := r.db.Create(order).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) FindAll() ([]entity.Order, error) {
	orders := []entity.Order{}

	err := r.db.Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}
