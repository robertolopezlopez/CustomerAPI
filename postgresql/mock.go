package postgresql

import (
	"api/customer"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type DataBaseMock struct {
	mock.Mock
}

func (d *DataBaseMock) DeleteByMailingID(mailingID int64) *gorm.DB {
	args := d.Called(mailingID)
	return &gorm.DB{
		RowsAffected: args.Get(0).(int64),
		Error:        args.Error(1),
	}
}

func (d *DataBaseMock) Create(customer *customer.Customer) *gorm.DB {
	args := d.Called(customer)
	return &gorm.DB{Error: args.Error(0)}
}

func (d *DataBaseMock) Migrate(customer *customer.Customer) error {
	args := d.Called(customer)
	return args.Error(0)
}

func (d *DataBaseMock) Delete(customer *customer.Customer, id int64) *gorm.DB {
	args := d.Called(customer, id)
	return &gorm.DB{Error: args.Error(0)}
}

func (d *DataBaseMock) First(id int64) (c customer.Customer, tx *gorm.DB) {
	args := d.Called(id)
	tx = &gorm.DB{Error: args.Error(1)}
	if tx.Error == nil {
		c = args.Get(0).(customer.Customer)
	}
	return
}

func (d *DataBaseMock) Find() ([]customer.Customer, *gorm.DB) {
	args := d.Called()
	return args.Get(0).([]customer.Customer), &gorm.DB{
		Error: args.Error(1),
	}
}

func (d *DataBaseMock) DeleteOld(seconds int) *gorm.DB {
	args := d.Called(seconds)
	return &gorm.DB{
		RowsAffected: args.Get(0).(int64),
		Error:        args.Error(1),
	}
}
