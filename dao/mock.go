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

func (pg *CustomerDaoMock) Create(c customer.Customer) error {
	args := pg.Called(c)
	return args.Error(0)
}

func (pg *CustomerDaoMock) MigrateModels() error {
	args := pg.Called()
	return args.Error(0)
}
