package main

import (
	"api/authentication"
	"api/dao"
	"api/handler"
	"api/logging"
	"api/tracing"
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

var (
	C = &handler.CustomerHandler{DAO: dao.DAO}
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(authentication.HeaderAuthMiddleware())
	r.Use(tracing.XRequestIDMiddleware())

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// create a client entry
	r.POST("/api/clients", C.CreateCustomer)

	// Get client by id
	r.GET("/api/clients/:id", handler.GetCustomerHandler)

	// Delete client by id
	r.DELETE("/api/clients/:id", handler.DeleteCustomerHandler)

	// Get all clients
	r.GET("/api/clients/", handler.FindCustomersHandler)

	return r
}

func main() {
	if err := C.DAO.MigrateModels(); err != nil {
		panic(err)
	}
	// Listen and serve in 0.0.0.0:8080
	_ = SetupRouter().Run(":8080")
}

// TODO CRON https://github.com/robfig/cron or https://github.com/go-co-op/gocron
// TODO MORE TESTS
