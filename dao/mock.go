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

func (dao *CustomerDaoMock) Create(c *customer.Customer) error {
	args := dao.Called(c)
	err := args.Error(0)
	if err == nil {
		c = args.Get(0).(*customer.Customer)
	}
	return err
}

func (dao *CustomerDaoMock) MigrateModels() error {
	args := dao.Called()
	return args.Error(0)
}

func (dao *CustomerDaoMock) Delete(c *customer.Customer, id int64) error {
	args := dao.Called(c, id)
	return args.Error(0)
}

func (dao *CustomerDaoMock) First(id int64) (*customer.Customer, error) {
	args := dao.Called(id)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	return &customer.Customer{Model: gorm.Model{ID: uint(id)}}, nil
}

func (dao *CustomerDaoMock) Find() ([]customer.Customer, error) {
	args := dao.Called()
	return args.Get(0).([]customer.Customer), args.Error(1)
}

func (dao *CustomerDaoMock) DeleteOld(seconds int) (int64, error) {
	args := dao.Called(seconds)
	return args.Get(0).(int64), args.Error(1)
}
