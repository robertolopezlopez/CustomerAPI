package postgresql

import (
	"api/customer"

	"gorm.io/gorm"
)

func (d *DataBaseMock) Create(customer *customer.Customer) *gorm.DB {
	args := d.Called(customer)
	return &gorm.DB{Error: args.Error(0)}
}

func (d *DataBaseMock) Migrate(customer *customer.Customer) error {
	args := d.Called(customer)
	return args.Error(0)
}
