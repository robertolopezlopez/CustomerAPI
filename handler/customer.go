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
		// CreateCustomer handles POST /api/clients
		CreateCustomer(*gin.Context)
		// GetCustomer handles GET /api/clients/:id
		GetCustomer(*gin.Context)
		// DeleteCustomer handles DELETE /api/clients/:id
		DeleteCustomer(*gin.Context)
		// FindCustomers handles GET /api/clients
		FindCustomers(*gin.Context)
	}

	CustomerHandler struct {
	}
)

func (c *CustomerHandler) GetCustomer(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Params.ByName("id"), 10, 64)
	if err != nil {
		logging.WarnLogger.Printf("%s: %s", ctx.Request.Header.Get(tracing.XRequestID), err.Error())
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cust, err := dao.DAO.First(id)
	if errors.Is(err, logger.ErrRecordNotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		logging.ErrorLogger.Printf("error querying the DB: %s\n", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.IndentedJSON(http.StatusOK, cust)
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

	err := dao.DAO.Create(&newCustomer)
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

func (c *CustomerHandler) DeleteCustomer(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Params.ByName("id"), 10, 64)
	if err != nil {
		logging.WarnLogger.Printf("%s: %s", ctx.Request.Header.Get(tracing.XRequestID), err.Error())
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = dao.DAO.Delete(&customer.Customer{}, id)
	if err != nil {
		logging.ErrorLogger.Printf("error deleting from the DB: %s\n", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.Status(http.StatusNoContent)
}

func (c *CustomerHandler) FindCustomers(ctx *gin.Context) {
	// todo pagination
	customers, err := dao.DAO.Find()
	if err != nil {
		logging.WarnLogger.Printf("error querying the DB: %s\n", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.IndentedJSON(http.StatusOK, customers)
}
