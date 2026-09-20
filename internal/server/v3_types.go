package server

import (
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
)

// Path parameter input structs used by Huma handlers.

type univInput struct {
	UnivID int `path:"univId" doc:"University ID"`
}

type univGroupInput struct {
	UnivID  int `path:"univId" doc:"University ID"`
	GroupID int `path:"groupId" doc:"Group ID"`
}

type univGroupAdeInput struct {
	UnivID       int `path:"univId" doc:"University ID"`
	GroupID      int `path:"groupId" doc:"Group ID"`
	AdeResources int `path:"adeResources" doc:"ADE resources identifier"`
}

type univAdeInput struct {
	UnivID       int `path:"univId" doc:"University ID"`
	AdeResources int `path:"adeResources" doc:"ADE resources identifier"`
}

type freeRoomsInput struct {
	UnivID   int       `path:"univId" doc:"University ID"`
	Start    time.Time `query:"start" required:"false" doc:"Search start hour (default: now)"`
	End      time.Time `query:"end" required:"false" doc:"Search end hour (default: start + 1h)"`
	CampusID int       `query:"campusId" required:"false" default:"0" doc:"Campus ID (optional)"`
}
type campusInput struct {
	UnivID   int `path:"univId" doc:"University ID"`
	CampusID int `path:"campusId" doc:"Campus ID"`
}
type campusesInput struct {
	UnivID int `path:"univId" doc:"University ID"`
}

// JSON body response shapes

type universityResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	AdeUrl string `json:"adeUrl"`
}

type groupResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type timetableResponse struct {
	AdeResources int        `json:"adeResources"`
	AdeProjectId int        `json:"adeProjectId"`
	Year         int        `json:"year"`
	Label        string     `json:"label"`
	AdeUrl       string     `json:"adeUrl"`
	LastUpdate   *time.Time `json:"lastUpdate"`
}

type roomResponse struct {
	AdeResources int        `json:"adeResources"`
	AdeProjectId int        `json:"adeProjectId"`
	Label        string     `json:"label"`
	AdeUrl       string     `json:"adeUrl"`
	LastUpdate   *time.Time `json:"lastUpdate"`
	CampusID     int        `json:"campusId"`
}

type campusResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Huma response envelope types — each wraps a Body field for the JSON response.

type universityListOutput struct {
	Body []universityResponse
}

type universityOutput struct {
	Body universityResponse
}

type groupListOutput struct {
	Body []groupResponse
}

type timetableListOutput struct {
	Body []timetableResponse
}

type timetableOutput struct {
	Body timetableResponse
}

type eventListOutput struct {
	Body []ade.Event
}

type roomListOutput struct {
	Body []roomResponse
}

type roomOutput struct {
	Body roomResponse
}
type campusListOutput struct {
	Body []campusResponse
}
type campusOutput struct {
	Body campusResponse
}
