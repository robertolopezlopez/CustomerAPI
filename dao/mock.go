package dao

import (
	"api/customer"

	"github.com/stretchr/testify/mock"
)

type (
	CustomerDaoMock struct {
		mock.Mock
	}
)

func (dao *CustomerDaoMock) Create(c *customer.Customer) error {
	args := dao.Called(c)
	return args.Error(0)
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
	first := args.Get(0)
	if first == nil {
		return nil, args.Error(1)
	}
	return first.(*customer.Customer), args.Error(1)
}

func (dao *CustomerDaoMock) Find() ([]customer.Customer, error) {
	args := dao.Called()
	return args.Get(0).([]customer.Customer), args.Error(1)
}

func (dao *CustomerDaoMock) DeleteOld(seconds int) (int64, error) {
	args := dao.Called(seconds)
	return args.Get(0).(int64), args.Error(1)
}

func (dao *CustomerDaoMock) DeleteByMailingID(mailingID int64) (int64, error) {
	args := dao.Called(mailingID)
	return args.Get(0).(int64), args.Error(1)
}
