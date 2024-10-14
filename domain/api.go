package domain

import "time"

type UniversityResponse struct {
	NumUniv  int    `json:"numUniv"`
	NameUniv string `json:"nameUniv"`
	AdeUniv  string `json:"adeUniv"`
}

type TimetableResponse struct {
	NumUniv      int       `json:"numUniv"`
	NameUniv     string    `json:"nameUniv"`
	NameTT       string    `json:"nameTT"`
	DescTT       string    `json:"descTT"`
	NumYearTT    int       `json:"numYearTT"`
	AdeResources int       `json:"adeResources"`
	AdeProjectId int       `json:"adeProjectId"`
	LastUpdate   time.Time `json:"lastUpdate"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type JsonEvent struct {
	Title       string    `json:"title"`
	Teacher     string    `json:"teacher"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
}
