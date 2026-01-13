package station

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MouslyCode/mrt-schedules/common/client"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
	CheckScheduleByStations(id string) (response []ScheduleResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStations() (response []StationResponse, err error) {
	// layer Service
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"
	// url := "https://www.jakartamrt.co.id/"

	// Hit URL
	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var stations []Station
	err = json.Unmarshal(byteResponse, &stations)
	if err != nil {
		return
	}

	// Response
	for _, item := range stations {
		response = append(response, StationResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}
	return
}

func (s *service) CheckScheduleByStations(id string) (response []ScheduleResponse, err error) {
	// Layer Service
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	// Hit URL
	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var schedules []Schedule
	err = json.Unmarshal(byteResponse, &schedules)

	// Response
	var scheduleSelected Schedule
	for _, item := range schedules {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		err = errors.New("Station not Found")
		return
	}

	response, err := ConvertDataToScheduleResponse(scheduleSelected)
	if err != nil {
		return
	}

	return

}

func ConvertDataToScheduleResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		lebakBulusTripName = "Stasiun Lebak Bulus Grab"
		bundaranHITripName = "Stasiun Bundaran HI Grab"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParsed, err := ConvertDataToScheduleResponse(scheduleLebakBulus)
	if err != nil {
		return 
	}

	scheduleBunderanHIParsed, err := ConvertDataToScheduleResponse(scheduleBundaranHI)
	if err != nil {
		return 
	}

	// convert to Response
	// for _, item := range scheduleLebakBulusParsed {
	// 	if item.Format("15:04") > time.Now().Format("15:04") {
	// 		response = append(response, ScheduleResponse{
	// 			StationName: lebakBulusTripName,
	// 			Time: item.Format("15:04"),
	// 		})
	// 	}
	// }

	for _, item := range scheduleBunderanHIParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: bundaranHITripName,
				Time: item.Format("15:04"),
			})
		}
	}
	
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules = strings.Split(schedule, "," )
	)

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item)
		if trimmedTime == "" {
			continue
		}
		parsedTime, err = time.Parse("15:04", trimmedTime)
		if err != nil {
			err = errors.New("Invalid Time Format" + trimmedTime)
			return
		}
		response = append(response, parsedTime)
	}
	return
}
