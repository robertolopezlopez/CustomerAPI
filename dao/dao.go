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
		Create(*customer.Customer) error

		// Delete deletes a customer.Customer from the database. It may return ErrPg.
		Delete(*customer.Customer, int64) error

		// MigrateModels applies any possible modifications to the underlying database schema. It may return ErrPg.
		MigrateModels() error

		// First retrieves customer.Customer by primary key. It may return ErrPg.
		First(int64) (*customer.Customer, error)

		// Find retrieves all customer.Customer. It may return ErrPg
		Find() ([]customer.Customer, error)

		// DeleteOld handles removal of database entries older than 5 minutes
		DeleteOld(int) (int64, error)
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

func (dao *CustomerDAO) MigrateModels() error {
	return dao.Db.Migrate(&customer.Customer{})
}

func (dao *CustomerDAO) Create(c *customer.Customer) error {
	if tx := dao.Db.Create(c); tx.Error != nil {
		if strings.Contains(tx.Error.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("%w: %s", ErrPgIndex, tx.Error.Error())
		}
		return fmt.Errorf("%w: %s", ErrPg, tx.Error.Error())
	}
	return nil
}

func (dao *CustomerDAO) Delete(c *customer.Customer, id int64) error {
	if tx := dao.Db.Delete(c, id); tx.Error != nil {
		return fmt.Errorf("%w: delete: %s", ErrPg, tx.Error.Error())
	}
	return nil
}

func (dao *CustomerDAO) First(id int64) (*customer.Customer, error) {
	c, tx := dao.Db.First(id)
	if tx.Error != nil {
		return nil, fmt.Errorf("%w: first: %s", ErrPg, tx.Error.Error())
	}
	return &c, nil
}

func (dao *CustomerDAO) Find() ([]customer.Customer, error) {
	customers, tx := dao.Db.Find()
	if tx.Error != nil {
		return nil, fmt.Errorf("%w: find: %s", ErrPg, tx.Error.Error())
	}
	return customers, nil
}

func (dao *CustomerDAO) DeleteOld(seconds int) (int64, error) {
	tx := dao.Db.DeleteOld(seconds)
	if tx.Error != nil {
		return tx.RowsAffected, fmt.Errorf("%w: delete old: %s", ErrPg, tx.Error.Error())
	}
	return tx.RowsAffected, tx.Error
}
