package station

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MouslyCode/mrt-jakarta-api/common/client"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
	CheckScheduleByStations(id string) (response []ScheduleResponse, err error)
	CheckEstimateByStations(id string) (response []StationEstimateResponse, err error)
	CheckFacilityByStations(id string) (response []StationFacilityResponse, err error)
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

	response, err = ConvertDataToResponse(scheduleSelected)
	if err != nil {
		return
	}

	return

}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		lebakBulusTripName = "Stasiun Lebak Bulus Grab"
		bundaranHITripName = "Stasiun Bundaran HI Grab"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBunderanHIParsed, err := ConvertScheduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return
	}

	// convert to Response
	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04:05") > time.Now().Format("15:04:05") {
			response = append(response, ScheduleResponse{
				StationName: lebakBulusTripName,
				Time:        item.Format("15:04:05"),
			})
		}
	}

	for _, item := range scheduleBunderanHIParsed {
		if item.Format("15:04:05") > time.Now().Format("15:04:05") {
			response = append(response, ScheduleResponse{
				StationName: bundaranHITripName,
				Time:        item.Format("15:04:05"),
			})
		}
	}
	return
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules  = strings.Split(schedule, ",")
	)

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item)
		if trimmedTime == "" {
			continue
		}
		parsedTime, err = time.Parse("15:04:05", trimmedTime)
		if err != nil {
			err = errors.New("Invalid Time Format" + trimmedTime)
			return
		}
		response = append(response, parsedTime)
	}
	return
}

func (s *service) CheckEstimateByStations(id string) (response []StationEstimateResponse, err error) {
	// Layer Service
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

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

	stationNameById := make(map[string]string, len(stations))
	for _, s := range stations {
		stationNameById[s.Id] = s.Name
	}

	for _, station := range stations {

		if station.Id == id {
			continue
		}

		var estimates []EstimateResponse
		for _, est := range station.StationEstimate {
			stationName := stationNameById[est.StationId]
			estimates = append(estimates, EstimateResponse{
				StationName: stationName,
				Fare:        est.Fare,
				Time:        est.Time,
			})
		}

		response = append(response, StationEstimateResponse{
			StationName: station.Name,
			Estimates:   estimates,
		})

		return
	}

	err = errors.New("Station Not Found")

	return
}

func (s *service) CheckFacilityByStations(id string) (response []StationFacilityResponse, err error) {
	// Layer Service
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

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

	for _, station := range stations {

		if station.Id == id {
			continue
		}

		var facilities []FacilityResponse
		for _, facility := range station.StationFacility {
			facilities = append(facilities, FacilityResponse{
				Title: facility.Title,
				Type:  facility.Type,
				Img:   facility.Img,
			})
		}

		response = append(response, StationFacilityResponse{
			StationName: station.Name,
			Facilities:  facilities,
		})

		return
	}

	err = errors.New("Station Not found")

	return
}
