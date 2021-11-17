package handler

import (
	"api/customer"
	"api/dao"
	"api/logging"
	"api/tracing"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

type (
	customerHandler interface {
		CreateCustomer(ctx *gin.Context) error
	}

	CustomerHandler struct {
		DAO *dao.CustomerDAO
	}
)

func GetCustomerHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		logging.WarnLogger.Printf("%s: %s", c.Request.Header.Get(tracing.XRequestID), err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cust, err := dao.DAO.First(id)
	if err == nil {
		c.IndentedJSON(http.StatusOK, cust)
	} else {
		if errors.Is(err, logger.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		logging.ErrorLogger.Printf("error querying the DB: %s\n", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *CustomerHandler) CreateCustomer(ctx *gin.Context) {
	var newCustomer customer.Customer
	if err := ctx.BindJSON(&newCustomer); err != nil {
		return
	}

	if err := newCustomer.Validate(); err != nil {
		logging.WarnLogger.Printf("%s: %s", ctx.Request.Header.Get(tracing.XRequestID), err.Error())
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := c.DAO.Create(newCustomer)
	if err != nil {
		if errors.Is(err, dao.ErrPgIndex) {
			logging.WarnLogger.Printf("%s: %s", ctx.Request.Header.Get(tracing.XRequestID), err.Error())
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		logging.ErrorLogger.Printf("%s: %s", ctx.Request.Header.Get(tracing.XRequestID), err.Error())
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	js, err := json.Marshal(newCustomer)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logging.InfoLogger.Printf("%s : %s", ctx.Request.Header.Get(tracing.XRequestID), string(js))

	ctx.Status(http.StatusCreated)
}

func DeleteCustomerHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		logging.WarnLogger.Printf("%s: %s", c.Request.Header.Get(tracing.XRequestID), err.Error())
		// todo how to use AborWithError?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = dao.DAO.Delete(&customer.Customer{}, id)
	if err == nil {
		c.Status(http.StatusNoContent)
	} else {
		logging.ErrorLogger.Printf("error deleting from the DB: %s\n", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func FindCustomersHandler(c *gin.Context) {
	var customers []customer.Customer
	// todo gorm with pagination FindInBatches
	err := dao.DAO.Find(&customers)
	if err == nil {
		c.IndentedJSON(http.StatusOK, customers)
	} else {
		logging.WarnLogger.Printf("error querying the DB: %s\n", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
