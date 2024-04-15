package main

import (
	"time"

	"github.com/Mehul-Kumar-27/HotelReservation/gateway/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(
		cors.New(
			cors.Config{
				AllowAllOrigins: true,
				AllowMethods:     []string{"OPTIONS", "PUT", "PATCH", "GET", "DELETE"},
				AllowHeaders:     []string{"Content-Type", "Content-Length", "Origin"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)

	group := router.Group("/api/v1/")
	api.Admin.AddToGroup(group)

	router.Run(":8080")

}
