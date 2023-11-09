package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/ardanlabs/conf/v3"
	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/api"
)

type appConfig struct {
	Port int `conf:"default:8080,env:APP_PORT"`
}

func main() {

	// CONFIG
	var cfg appConfig
	help, err := conf.Parse("APP", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
		log.Fatalf("Error parsing configuration: %v\n", err)
		return
	}

	// Router
	r := gin.Default()

	v1api := r.Group("/api/v1")

	v1api.GET("/user", api.HandleGetUsers)
	v1api.GET("/user/:id", api.HandleGetUserById)

	addr := fmt.Sprintf(":%d", cfg.Port)
	r.Run(addr)

}
