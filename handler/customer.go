package handler

import (
	"api/customer"
	"api/db"
	"api/logging"
	"api/tracing"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func GetCustomerHandler(c *gin.Context) {
	id := c.Params.ByName("id")

	var cust customer.Customer
	tx := db.DB.First(&cust, id)
	if tx.Error == nil {
		c.IndentedJSON(http.StatusOK, cust)
	} else {
		if errors.Is(tx.Error, logger.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		logging.WarnLogger.Printf("error querying the DB: %s\n", tx.Error.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func CreateCustomerHandler(ctx *gin.Context) {
	var c customer.Customer
	if err := ctx.BindJSON(&c); err != nil {
		return
	}

	if err := c.Validate(); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx := db.DB.Create(&c)
	if tx.Error != nil {
		logging.ErrorLogger.Printf("error writing into DB: %s\n", tx.Error.Error())
		return
	}

	cs, _ := json.Marshal(c)
	logging.InfoLogger.Printf("%s : %s", ctx.Request.Header.Get(tracing.XRequestID), string(cs))

	ctx.Status(http.StatusCreated)
}

func DeleteCustomerHandler(c *gin.Context) {
	id := c.Params.ByName("id")

	var cust customer.Customer
	tx := db.DB.Delete(&cust, id)
	if tx.Error == nil {
		c.Status(http.StatusNoContent)
	} else {
		logging.WarnLogger.Printf("error deleting from the DB: %s\n", tx.Error.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func FindCustomersHandler(c *gin.Context) {
	var customers []customer.Customer
	tx := db.DB.Find(&customers)
	if tx.Error == nil {
		c.IndentedJSON(http.StatusOK, customers)
	} else {
		logging.WarnLogger.Printf("error querying the DB: %s\n", tx.Error.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
