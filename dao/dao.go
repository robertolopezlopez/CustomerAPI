package dao

import (
	"api/customer"
	"api/postgresql"
	"errors"
	"fmt"
	"strings"
)

type (
	customerDao interface {
		// Create creates a new customer.Customer in the database. It may return ErrPgIndex or a generic ErrPg
		Create(customer.Customer) error

		// MigrateModels applies any possible modifications to the underlying database schema.
		MigrateModels() error
	}

	CustomerDao struct {
		Db postgresql.Db
	}
)

var (
	ORM        = &CustomerDao{Db: postgresql.DB}
	ErrPgIndex = errors.New("duplicate key value for Tx index")
	ErrPg      = errors.New("database error")
)

func (pg *CustomerDao) MigrateModels() error {
	return pg.Db.Migrate(&customer.Customer{})
}

func (pg *CustomerDao) Create(c customer.Customer) error {
	if tx := pg.Db.Create(&c); tx.Error != nil {
		if strings.Contains(tx.Error.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("%w: %s", ErrPgIndex, tx.Error.Error())
		}
		return fmt.Errorf("%w: %s", ErrPg, tx.Error.Error())
	}
	return nil
}
