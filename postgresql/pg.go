package postgresql

import (
	"api/customer"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	// Db is the closest representation to the database level.
	// Used to mock dao calls.
	Db interface {
		// Create handles calls to &gorm.DB.Create(*customer.Customer)
		Create(*customer.Customer) *gorm.DB
		// Migrate handles calls to &gorm.DB.Automigrate(*customer.Customer)
		Migrate(*customer.Customer) error
	}
	DBase struct {
		Tx *gorm.DB
	}
)

var (
	DB *DBase
)

func init() {
	// Establish connection to local PostgreSQL
	var err error
	dsn := "host=localhost user=postgres password=example dbname=customer port=5432 sslmode=disable TimeZone=Europe/Warsaw"
	pg, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = &DBase{Tx: pg}
}

func (d *DBase) Migrate(customer *customer.Customer) error {
	return d.Tx.AutoMigrate(customer)
}

func (d *DBase) Create(customer *customer.Customer) *gorm.DB {
	return d.Tx.Create(customer)
}
