package station

import (
	"net/http"

	"github.com/MouslyCode/mrt-jakarta-api/common/response"
	"github.com/gin-gonic/gin"
)

func Initiate(router *gin.RouterGroup) {

	stationService := NewService()

	station := router.Group("/stations")
	station.GET("", func(c *gin.Context) {
		// Code Services
		GetAllStations(c, stationService)
	})

	station.GET("/:id", func(c *gin.Context) {
		CheckScheduleByStations(c, stationService)
	})

	station.GET("/:id/estimate", func(c *gin.Context) {
		CheckEstimateByStations(c, stationService)
	})
}

func GetAllStations(c *gin.Context, service Service) {
	datas, err := service.GetAllStations()

	if err != nil {
		// handle error
		c.JSON(http.StatusBadRequest,
			response.APIResponse{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			},
		)
		return
	} else {
		c.JSON(http.StatusOK, response.APIResponse{
			Success: true,
			Message: "Success Get All Stations",
			Data:    datas,
		})
	}

	// response

}

func CheckScheduleByStations(c *gin.Context, service Service) {
	id := c.Param("id")

	datas, err := service.CheckScheduleByStations(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	} else {
		c.JSON(http.StatusOK, response.APIResponse{
			Success: true,
			Message: "Success Get All Schedule",
			Data:    datas,
		})
	}

}

func CheckEstimateByStations(c *gin.Context, service Service) {
	id := c.Param("id")

	datas, err := service.CheckEstimateByStations(id)
	if err != nil {
		c.JSON(http.StatusBadGateway, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	} else {
		c.JSON(http.StatusOK, response.APIResponse{
			Success: true,
			Message: "Succes Get Estimate From ",
			Data:    datas,
		})
	}

}
