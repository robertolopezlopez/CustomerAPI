package main

import (
	"api/authentication"
	"api/customer"
	"api/db"
	"api/logging"
	"api/tracing"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

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

	r.POST("/api/clients", createCustomerHandler)

	// Get user value
	r.GET("/api/clients/:id", getCustomerHandler)

	return r
}

func getCustomerHandler(c *gin.Context) {
	id := c.Params.ByName("id")

	var cust customer.Customer
	value := db.DB.First(&cust, id)
	if value.Error == nil {
		c.JSON(http.StatusOK, cust)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
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

	db.DB.Create(&c)

	cs, _ := json.Marshal(c)
	logging.InfoLogger.Printf("%s : %s", ctx.Request.Header.Get(tracing.XRequestID), string(cs))

	ctx.Status(http.StatusCreated)

}

func main() {
	r := SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	_ = r.Run(":8080")
}
