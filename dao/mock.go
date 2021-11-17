package dao

import (
	"api/customer"

	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
)

type (
	CustomerDaoMock struct {
		mock.Mock
	}
)

func (pg *CustomerDaoMock) Create(c customer.Customer) error {
	args := pg.Called(c)
	return args.Error(0)
}

func (pg *CustomerDaoMock) MigrateModels() error {
	args := pg.Called()
	return args.Error(0)
}

func (pg *CustomerDaoMock) Delete(c *customer.Customer, id int64) error {
	args := pg.Called(c, id)
	return args.Error(0)
}

func (pg *CustomerDaoMock) First(id int64) (*customer.Customer, error) {
	args := pg.Called(id)
	err := args.Error(0)
	if err != nil {
		return nil, err
	}
	return &customer.Customer{Model: gorm.Model{ID: uint(id)}}, nil
}
