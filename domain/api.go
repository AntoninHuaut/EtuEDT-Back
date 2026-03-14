package domain

import "time"

type UniversityResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	AdeUrl string `json:"adeUrl"`
}

type GroupResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TimetableResponse struct {
	AdeResources int       `json:"adeResources"`
	AdeProjectId int       `json:"adeProjectId"`
	Year         int       `json:"year"`
	Label        string    `json:"label"`
	AdeUrl       string    `json:"adeUrl"`
	LastUpdate   time.Time `json:"lastUpdate"`
}

type RoomResponse struct {
	AdeResources int       `json:"adeResources"`
	AdeProjectId int       `json:"adeProjectId"`
	Label        string    `json:"label"`
	AdeUrl       string    `json:"adeUrl"`
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
