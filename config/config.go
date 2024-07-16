package config

import (
	"template/database"
	"template/router"
)

func Initialization() {
	database.Initialization()
	router.Routes("3000")
}
