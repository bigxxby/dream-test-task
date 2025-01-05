package app

import (
	"log"

	"github.com/bigxxby/dream-test-task/internal/config"
	"github.com/bigxxby/dream-test-task/internal/database/connection"
	"github.com/bigxxby/dream-test-task/internal/database/migration"
	"github.com/bigxxby/dream-test-task/internal/router"
)

func Run() {
	//make log flags to show file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config, err := config.GetCofig()
	if err != nil {
		log.Println(err)
		return
	}

	db, err := connection.GetDB(config)
	if err != nil {
		log.Println(err)
		return
	}

	// migrate db
	err = migration.Migrate(db)
	if err != nil {
		log.Println(err)
		return
	}

	// err = utils.CreateAdmin(db)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	router, err := router.NewRouter(db)
	if err != nil {
		log.Println(err)
		return
	}
	router.Run(":" + config.AppPort)
}
