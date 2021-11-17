package dao

import (
	"api/customer"
	"api/postgresql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDao_Create(t *testing.T) {
	tests := map[string]struct {
		db        *postgresql.DataBaseMock
		input     customer.Customer
		withError error
	}{
		"OK": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("Create", mock.Anything).Return(nil)
				return &m
			}(),
		},
		"index error": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("Create", mock.Anything).Return(fmt.Errorf("duplicate key value violates unique constraint"))
				return &m
			}(),
			withError: ErrPgIndex,
		},
		"other error": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("Create", mock.Anything).Return(fmt.Errorf("other error"))
				return &m
			}(),
			withError: ErrPg,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao := CustomerDAO{Db: test.db}
			err := dao.Create(test.input)
			if test.withError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, test.withError))
				return
			}
			require.NoError(t, err)
			test.db.AssertExpectations(t)
		})
	}
}

func TestDao_MigrateModels(t *testing.T) {
	tests := map[string]struct {
		db        *postgresql.DataBaseMock
		withError *regexp.Regexp
	}{
		"OK": {
			db: func() *postgresql.DataBaseMock {
				m := &postgresql.DataBaseMock{}
				m.On("Migrate", mock.Anything).Return(nil)
				return m
			}(),
		},
		"not OK": {
			db: func() *postgresql.DataBaseMock {
				m := &postgresql.DataBaseMock{}
				m.On("Migrate", mock.Anything).Return(fmt.Errorf("an error"))
				return m
			}(),
			withError: regexp.MustCompile("an error"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao := CustomerDAO{Db: test.db}
			err := dao.MigrateModels()
			if test.withError != nil {
				assert.Regexp(t, err, test.withError)
				return
			}
			require.NoError(t, err)
			test.db.AssertExpectations(t)
		})
	}
}

func TestCustomerDAO_Delete(t *testing.T) {
	tests := map[string]struct {
		db        *postgresql.DataBaseMock
		withError error
	}{
		"OK deletion": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("Delete", mock.Anything, mock.Anything).Return(nil)
				return &m
			}(),
		},
		"NOK deletion": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("Delete", mock.Anything, mock.Anything).Return(fmt.Errorf("an error"))
				return &m
			}(),
			withError: ErrPg,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao := CustomerDAO{Db: test.db}
			err := dao.Delete(&customer.Customer{}, 1)
			if test.withError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, test.withError))
				return
			}
			require.NoError(t, err)
			test.db.AssertExpectations(t)
		})
	}
}
