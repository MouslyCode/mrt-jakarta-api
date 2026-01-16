package station

type Station struct {
	Id   string `json:"nid"`
	Name string `json:"title"`
}

type StationResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Schedule struct {
	StationId          string `json:"nid"`
	StationName        string `json:"title"`
	ScheduleBundaranHI string `json:"jadwal_hi_biasa"`
	ScheduleLebakBulus string `json:"jadwal_lb_biasa"`
}

type ScheduleResponse struct {
	StationName string `json:"station"`
	Time        string `json:"time"`
}

func (s ScheduleResponse) Format(param1 string) {
	panic("unimplemented")
}

type EstimateRequest struct {
	Estimate []struct {
		StationId string `json:"stasiun_nid"`
		Fare      string `json:"tarif"`
		Time      string `json:"waktu"`
	} `json:"estimasi"`
}

type EstimateResponse struct {
	Fare string `json:"fare"`
	Time string `json:"time"`
}
