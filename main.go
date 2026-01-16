package main

import (
	"github.com/MouslyCode/mrt-jakarta-api/modules/station"
	"github.com/gin-gonic/gin"
)

func main() {
	InitiateRouter()
}

func InitiateRouter() {
	router := gin.Default()
	api := router.Group("/v1/api")

	station.Initiate(api)

	router.Run(":8080")
}
