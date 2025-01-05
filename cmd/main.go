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

func main() {
	app.Run()
}
