package main

import (
	"log"

	"github.com/bigxxby/dream-test-task/internal/app"
	"github.com/bigxxby/dream-test-task/internal/config"
)

// set config first
func init() {
	err := config.SetConfig()
	if err != nil {
		log.Println(err)
		return
	}
}

// @title dream-shortener API
// @version 1.0
// @description This is a dream-shortener API
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Provide your Bearer token in the format: Bearer <token>
func main() {
	app.App()
}
