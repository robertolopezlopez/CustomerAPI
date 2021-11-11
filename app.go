package main

import (
	"api/authentication"
	"api/customer"
	"api/db"
	"api/logging"
	"api/tracing"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	f, err := os.OpenFile(logging.GinLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.Use(authentication.HeaderAuthMiddleware())
	r.Use(tracing.XRequestIDMiddleware())

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// create a client entry
	r.POST("/api/clients", createCustomerHandler)

	// Get client by id
	r.GET("/api/clients/:id", getCustomerHandler)

	// Delete client by id
	r.DELETE("/api/clients/:id", deleteCustomerHandler)

	// Get all clients
	r.GET("/api/clients/", findCustomersHandler)

	return r
}

func findCustomersHandler(c *gin.Context) {
	var customers []customer.Customer
	tx := db.DB.Find(&customers)
	if tx.Error == nil {
		c.IndentedJSON(http.StatusOK, customers)
	} else {
		logging.WarnLogger.Printf("error querying the DB: %s\n", tx.Error.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func deleteCustomerHandler(c *gin.Context) {
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

func getCustomerHandler(c *gin.Context) {
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

func createCustomerHandler(ctx *gin.Context) {
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

func main() {
	r := SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	_ = r.Run(":8080")
}

// TODO CRON https://github.com/robfig/cron
// TODO MORE TESTS
