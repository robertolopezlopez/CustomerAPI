package postgresql

import (
	"api/customer"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	// Connection is the closest representation to the database level.
	// Used to mock db calls.
	Connection interface {
		// Create handles calls to &gorm.DB.Create(*customer.Customer)
		Create(*customer.Customer) *gorm.DB
		// Migrate handles calls to &gorm.DB.Automigrate(*customer.Customer)
		Migrate(*customer.Customer) error
	}

	DataBase struct {
		db *gorm.DB
	}
)

var (
	DB *DataBase
)

func init() {
	// Establish connection to local PostgreSQL
	var err error
	dsn := "host=localhost user=postgres password=example dbname=customer port=5432 sslmode=disable TimeZone=Europe/Warsaw"
	pg, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = &DataBase{db: pg}
}

func (d *DataBase) Create(customer *customer.Customer) *gorm.DB {
	return d.db.Create(customer)
}

func (d *DataBase) Migrate(customer *customer.Customer) error {
	return d.db.AutoMigrate(customer)
}
