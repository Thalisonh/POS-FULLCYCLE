package database

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/thalisonh/20-CleanArch/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *gorm.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:?_pragma=foreign_keys(1)"), &gorm.Config{})
	suite.NoError(err)
	db.AutoMigrate(&entity.Order{})
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	sqlDB, _ := suite.Db.DB()
	sqlDB.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.Where("id = ?", order.ID).First(&orderResult).Error

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}
