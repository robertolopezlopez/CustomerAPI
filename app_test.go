package main

import (
	"api/authentication"
	"api/customer"
	"api/dao"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"gorm.io/gorm"

	"gorm.io/gorm/logger"

	"github.com/stretchr/testify/assert"
)

var (
	oKheaders = map[string]string{authentication.AuthTokenHeader: authentication.AuthTokenValue}
)

func TestPingRoute(t *testing.T) {

	tests := map[string]struct {
		headers      map[string]string
		expectedCode int
		expectedBody string
	}{
		"200 ok": {
			headers:      oKheaders,
			expectedCode: http.StatusOK,
			expectedBody: "pong",
		},
		"401 nok": {
			headers:      map[string]string{authentication.AuthTokenHeader: ""},
			expectedCode: http.StatusUnauthorized,
		},
		"401 empty": {
			headers:      map[string]string{},
			expectedCode: http.StatusUnauthorized,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			router := SetupRouter()
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
			for key, value := range test.headers {
				req.Header.Add(key, value)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

func TestGetCustomer(t *testing.T) {
	aCustomer := customer.Customer{
		Model:     gorm.Model{ID: uint(1)},
		Email:     "oroparece@platano.es",
		Title:     "ninja",
		Content:   "content",
		MailingID: 1,
	}

	tests := map[string]struct {
		m            *dao.CustomerDaoMock
		expectedCode int
		expected     *customer.Customer
		id           string
	}{
		"OK 200": {
			id:           "1",
			expectedCode: http.StatusOK,
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("First", int64(1)).Return(&aCustomer, nil)
				return &m
			}(),
			expected: &aCustomer,
		},
		"NOK 500": {
			id:           "1",
			expectedCode: http.StatusInternalServerError,
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("First", int64(1)).Return(nil, errors.New("an error"))
				return &m
			}(),
		},
		"NOK 404": {
			id:           "1",
			expectedCode: http.StatusNotFound,
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("First", int64(1)).Return(nil, fmt.Errorf("%w: not found", logger.ErrRecordNotFound))
				return &m
			}(),
		},
		"NOK 400": {
			id:           "---",
			expectedCode: http.StatusBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao.DAO = test.m
			router := SetupRouter()
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/clients/%s", test.id), nil)
			for key, value := range oKheaders {
				req.Header.Add(key, value)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)

			response := w.Body.String()
			if test.expected == nil {
				assert.Empty(t, response)
				return
			}

			customerResponse := customer.Customer{}
			bytesResponse := []byte(response)
			err := json.Unmarshal(bytesResponse, &customerResponse)
			if err != nil {
				panic(err.Error())
			}

			assert.True(t, reflect.DeepEqual(customerResponse, *test.expected))

			test.m.AssertExpectations(t)
		})
	}
}

func TestCreateCustomer(t *testing.T) {
	tests := map[string]struct {
		m            *dao.CustomerDaoMock
		expectedCode int
		body         customer.Customer
		saved        customer.Customer
	}{
		"201 created": {
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("Create", mock.Anything).Return(nil)
				return &m
			}(),
			expectedCode: http.StatusCreated,
			saved: customer.Customer{
				Model:     gorm.Model{ID: uint(1)},
				Email:     "hello@example.com",
				Title:     "dev",
				Content:   "no content",
				MailingID: 1,
			},
			body: customer.Customer{
				Model:     gorm.Model{ID: uint(1)},
				Email:     "hello@example.com",
				Title:     "dev",
				Content:   "no content",
				MailingID: 1,
			},
		},
		"400 bad request validation": {
			body: customer.Customer{
				Model:     gorm.Model{ID: uint(1)},
				Email:     "---",
				Title:     "dev",
				Content:   "no content",
				MailingID: 1,
			},
			m:            func() *dao.CustomerDaoMock { return &dao.CustomerDaoMock{} }(),
			expectedCode: http.StatusBadRequest,
		},
		"500 internal error pg": {
			expectedCode: http.StatusInternalServerError,
			body: customer.Customer{
				Model:     gorm.Model{ID: uint(1)},
				Email:     "hello@example.com",
				Title:     "dev",
				Content:   "no content",
				MailingID: 1,
			},
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("Create", mock.Anything).Return(fmt.Errorf("an error"))
				return &m
			}(),
		},
		"400 bad request pg": {
			expectedCode: http.StatusBadRequest,
			body: customer.Customer{
				Model:     gorm.Model{ID: uint(1)},
				Email:     "hello@example.com",
				Title:     "dev",
				Content:   "no content",
				MailingID: 1,
			},
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("Create", mock.Anything).Return(fmt.Errorf("%w an error", dao.ErrPgIndex))
				return &m
			}(),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao.DAO = test.m
			router := SetupRouter()
			w := httptest.NewRecorder()

			bodyBytes, err := json.Marshal(test.body)
			if err != nil {
				panic(err.Error())
			}
			req, _ := http.NewRequest(http.MethodPost, "/api/clients", bytes.NewBuffer(bodyBytes))
			for key, value := range oKheaders {
				req.Header.Add(key, value)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)

			test.m.AssertExpectations(t)
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	tests := map[string]struct {
		m            *dao.CustomerDaoMock
		expectedCode int
		id           string
	}{
		"400 bad request": {
			id:           "--",
			m:            func() *dao.CustomerDaoMock { return &dao.CustomerDaoMock{} }(),
			expectedCode: http.StatusBadRequest,
		},
		"500 server error": {
			id: "1",
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("Delete", mock.Anything, int64(1)).Return(fmt.Errorf("an error"))
				return &m
			}(),
			expectedCode: http.StatusInternalServerError,
		},
		"204 no content": {
			id: "1",
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("Delete", mock.Anything, int64(1)).Return(nil)
				return &m
			}(),
			expectedCode: http.StatusNoContent,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao.DAO = test.m
			router := SetupRouter()
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/clients/%s", test.id), nil)
			for key, value := range oKheaders {
				req.Header.Add(key, value)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)

			test.m.AssertExpectations(t)
		})
	}
}

func TestFind(t *testing.T) {
	aCustomer := customer.Customer{
		Model:     gorm.Model{ID: uint(1)},
		Email:     "oroparece@platano.es",
		Title:     "ninja",
		Content:   "content",
		MailingID: 1,
	}

	tests := map[string]struct {
		m            *dao.CustomerDaoMock
		expectedCode int
		expected     []customer.Customer
	}{
		"500 nok": {
			expectedCode: http.StatusInternalServerError,
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("Find").Return([]customer.Customer{}, errors.New("an error"))
				return &m
			}(),
			expected: []customer.Customer{},
		},
		"200 ok": {
			expectedCode: http.StatusOK,
			m: func() *dao.CustomerDaoMock {
				m := dao.CustomerDaoMock{}
				m.On("Find").Return([]customer.Customer{aCustomer, aCustomer}, nil)
				return &m
			}(),
			expected: []customer.Customer{aCustomer, aCustomer},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao.DAO = test.m
			router := SetupRouter()
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/api/clients", nil)
			for key, value := range oKheaders {
				req.Header.Add(key, value)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)

			response := w.Body.String()

			var listResponse []customer.Customer
			bytesResponse := []byte(response)
			err := json.Unmarshal(bytesResponse, &listResponse)
			if err != nil {
				panic(err.Error())
			}

			assert.True(t, reflect.DeepEqual(listResponse, test.expected))

			test.m.AssertExpectations(t)
		})
	}
}
