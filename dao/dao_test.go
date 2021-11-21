package dao

import (
	"api/customer"
	"api/postgresql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"gorm.io/gorm"

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
			err := dao.Create(&test.input)
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

func TestCustomerDAO_First(t *testing.T) {
	tests := map[string]struct {
		id        int64
		c         customer.Customer
		withError *regexp.Regexp
		db        *postgresql.DataBaseMock
	}{
		"ok": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("First", int64(1)).Return(customer.Customer{
					Model:     gorm.Model{ID: 1},
					Email:     "oroparece@platano.es",
					Title:     "a client",
					MailingID: 1,
				}, nil)
				return &m
			}(),
			c: customer.Customer{
				Email:     "oroparece@platano.es",
				Title:     "a client",
				MailingID: 1,
				Model:     gorm.Model{ID: 1},
			},
		},
		"not ok, db error": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("First", int64(1)).Return(nil, fmt.Errorf("an error"))
				return &m
			}(),
			withError: regexp.MustCompile("database error: first: an error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao := CustomerDAO{Db: test.db}
			c, err := dao.First(1)
			if test.withError != nil {
				require.Error(t, err)
				assert.Regexp(t, test.withError, err.Error())
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.c.ID, c.ID)
			assert.Equal(t, test.c.Title, c.Title)
			assert.Equal(t, test.c.MailingID, c.MailingID)
			assert.Equal(t, test.c.Content, c.Content)

			test.db.AssertExpectations(t)
		})
	}
}

func TestCustomerDAO_Find(t *testing.T) {
	tests := map[string]struct {
		db        *postgresql.DataBaseMock
		withError *regexp.Regexp
		expected  []customer.Customer
	}{
		"ok": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("Find").Return([]customer.Customer{{
					Email:     "oroparece@platano.es",
					Title:     "a client",
					MailingID: 1,
					Model:     gorm.Model{ID: 1},
				}}, nil).Once()
				return &m
			}(),
			expected: []customer.Customer{{
				Email:     "oroparece@platano.es",
				Title:     "a client",
				MailingID: 1,
				Model:     gorm.Model{ID: 1},
			}},
		},
		"db returns an error": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("Find").Return([]customer.Customer{}, fmt.Errorf("an error")).Once()
				return &m
			}(),
			withError: regexp.MustCompile("database error: find: an error"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao := CustomerDAO{Db: test.db}
			cs, err := dao.Find()
			if test.withError != nil {
				require.Error(t, err)
				assert.Regexp(t, test.withError, err.Error())
				return
			}
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(cs, test.expected))

			test.db.AssertExpectations(t)
		})
	}
}

func TestDeleteOld(t *testing.T) {
	tests := map[string]struct {
		withError    *regexp.Regexp
		db           *postgresql.DataBaseMock
		rowsAffected int64
	}{
		"ok": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("DeleteOld", 1).Return(int64(1), nil)
				return &m
			}(),
			rowsAffected: 1,
		},
		"nok, orm error": {
			db: func() *postgresql.DataBaseMock {
				m := postgresql.DataBaseMock{}
				m.On("DeleteOld", 1).Return(int64(0), fmt.Errorf("an error"))
				return &m
			}(),
			withError: regexp.MustCompile("database error: delete old: an error"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao := CustomerDAO{Db: test.db}
			rowsAffected, err := dao.DeleteOld(1)
			if test.withError != nil {
				require.Error(t, err)
				assert.Regexp(t, test.withError, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.rowsAffected, rowsAffected)
		})
	}
}
