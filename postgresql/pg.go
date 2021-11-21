package postgresql

import (
	"api/customer"
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	// Db is the closest representation to the database level.
	// Used to mock dao calls.
	Db interface {
		// Create handles calls to &gorm.DB.Create()
		Create(*customer.Customer) *gorm.DB
		// Migrate handles calls to &gorm.DB.Automigrate()
		Migrate(*customer.Customer) error
		// Delete does soft delete
		Delete(*customer.Customer, int64) *gorm.DB
		// First handles calls to &gorm.DB.First()
		First(int64) (customer.Customer, *gorm.DB)
		// Find finds customers
		Find() ([]customer.Customer, *gorm.DB)
		// DeleteOld removes old entries from database (soft delete)
		DeleteOld(int) *gorm.DB
		// DeleteByMailingID removes entries from database with the given mailingID (soft delete)
		DeleteByMailingID(int64) *gorm.DB
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

func (d *DBase) Delete(c *customer.Customer, id int64) *gorm.DB {
	return d.Tx.Delete(c, id)
}

func (d *DBase) First(id int64) (c customer.Customer, tx *gorm.DB) {
	tx = d.Tx.First(c, id)
	return
}

func (d *DBase) Find() (cs []customer.Customer, tx *gorm.DB) {
	tx = d.Tx.Find(&cs)
	return
}

func (d *DBase) DeleteOld(seconds int) *gorm.DB {
	return d.Tx.Where(fmt.Sprintf("CreatedAt < NOW() - %d", seconds)).Delete(&customer.Customer{})
}

func (d *DBase) DeleteByMailingID(mailingID int64) *gorm.DB {
	return d.Tx.Where("mailingID = ?", strconv.FormatInt(mailingID, 10)).Delete(customer.Customer{})
}
