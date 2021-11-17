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
		// Create creates a new customer.Customer in the database. It may return ErrPgIndex or a generic ErrPg.
		Create(customer.Customer) error

		// Delete deletes a customer.Customer from the database. It may return ErrPg.
		Delete(*customer.Customer, int64) error

		// MigrateModels applies any possible modifications to the underlying database schema.
		MigrateModels() error
	}

	CustomerDAO struct {
		Db postgresql.Db
	}
)

var (
	DAO        = &CustomerDAO{Db: postgresql.DB}
	ErrPgIndex = errors.New("duplicate key value for Tx index")
	ErrPg      = errors.New("database error")
)

func (pg *CustomerDAO) MigrateModels() error {
	return pg.Db.Migrate(&customer.Customer{})
}

func (pg *CustomerDAO) Create(c customer.Customer) error {
	if tx := pg.Db.Create(&c); tx.Error != nil {
		if strings.Contains(tx.Error.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("%w: %s", ErrPgIndex, tx.Error.Error())
		}
		return fmt.Errorf("%w: %s", ErrPg, tx.Error.Error())
	}
	return nil
}

func (pg *CustomerDAO) Delete(c *customer.Customer, id int64) error {
	if tx := pg.Db.Delete(c, id); tx.Error != nil {
		return fmt.Errorf("%w: delete: %s", ErrPg, tx.Error.Error())
	}
	return nil
}
