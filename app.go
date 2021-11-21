package main

import (
	"api/authentication"
	"api/cron"
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
	// todo move this logic to logging package
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
	r.GET("/api/clients/:id", C.GetCustomer)

	// Delete client by id
	r.DELETE("/api/clients/:id", C.DeleteCustomer)

	// Get all clients
	r.GET("/api/clients/", C.FindCustomers)

	// todo POST /api/clients/send {mailing_id: 1} to write a log message and delete all customers with given

	return r
}

func main() {
	if err := C.DAO.MigrateModels(); err != nil {
		panic(err)
	}
	if _, err := cron.Scheduler(); err != nil {
		panic(err)
	}
	// Listen and serve in 0.0.0.0:8080
	_ = SetupRouter().Run(":8080")
}
