package server

import (
	"context"
	"net/http"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
	"github.com/danielgtaylor/huma/v2"
)

func registerV3Handlers(humaAPI huma.API) {
	huma.Register(humaAPI, huma.Operation{
		OperationID: "list-universities",
		Method:      http.MethodGet,
		Path:        "/v3/univs",
		Summary:     "List all universities",
		Tags:        []string{"Universities"},
	}, func(ctx context.Context, _ *struct{}) (*universityListOutput, error) {
		resp := make([]universityResponse, 0, len(config.AppConfig.Universities))
		for _, u := range config.AppConfig.Universities {
			resp = append(resp, universityResponse{ID: u.ID, Name: u.Name, AdeUrl: u.AdeUrl})
		}
		return &universityListOutput{Body: resp}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-university",
		Method:      http.MethodGet,
		Path:        "/v3/univs/{univId}",
		Summary:     "Get a university by ID",
		Tags:        []string{"Universities"},
	}, func(ctx context.Context, input *univInput) (*universityOutput, error) {
		univ, err := findUniversity(input.UnivID)
		if err != nil {
			return nil, humaError(err)
		}
		return &universityOutput{Body: universityResponse{ID: univ.ID, Name: univ.Name, AdeUrl: univ.AdeUrl}}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "list-groups",
		Method:      http.MethodGet,
		Path:        "/v3/univs/{univId}/groups",
		Summary:     "List groups for a university",
		Tags:        []string{"Groups"},
	}, func(ctx context.Context, input *univInput) (*groupListOutput, error) {
		univ, err := findUniversity(input.UnivID)
		if err != nil {
			return nil, humaError(err)
		}
		resp := make([]groupResponse, 0, len(univ.Groups))
		for _, g := range univ.Groups {
			resp = append(resp, groupResponse{ID: g.ID, Name: g.Name})
		}
		return &groupListOutput{Body: resp}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "list-timetables",
		Method:      http.MethodGet,
		Path:        "/v3/univs/{univId}/groups/{groupId}",
		Summary:     "List timetables for a group",
		Tags:        []string{"Timetables"},
	}, func(ctx context.Context, input *univGroupInput) (*timetableListOutput, error) {
		univ, err := findUniversity(input.UnivID)
		if err != nil {
			return nil, humaError(err)
		}
		group, err := findGroup(univ, input.GroupID)
		if err != nil {
			return nil, humaError(err)
		}
		firstDate, lastDate := ade.GetAcademicYearDates(time.Now())
		resp := make([]timetableResponse, 0, len(group.Timetables))
		for i := range group.Timetables {
			resp = append(resp, buildTimetableResponse(univ, &group.Timetables[i], firstDate, lastDate))
		}
		return &timetableListOutput{Body: resp}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-timetable-metadata",
		Method:      http.MethodGet,
		Path:        "/v3/univs/{univId}/groups/{groupId}/{adeResources}",
		Summary:     "Get timetable metadata",
		Tags:        []string{"Timetables"},
	}, func(ctx context.Context, input *univGroupAdeInput) (*timetableOutput, error) {
		univ, err := findUniversity(input.UnivID)
		if err != nil {
			return nil, humaError(err)
		}
		group, err := findGroup(univ, input.GroupID)
		if err != nil {
			return nil, humaError(err)
		}
		tt, err := findTimetable(group, input.AdeResources)
		if err != nil {
			return nil, humaError(err)
		}
		firstDate, lastDate := ade.GetAcademicYearDates(time.Now())
		return &timetableOutput{Body: buildTimetableResponse(univ, tt, firstDate, lastDate)}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-timetable-events",
		Method:      http.MethodGet,
		Path:        "/v3/univs/{univId}/groups/{groupId}/{adeResources}/events",
		Summary:     "Get timetable events (JSON)",
		Tags:        []string{"Timetables"},
	}, func(ctx context.Context, input *univGroupAdeInput) (*eventListOutput, error) {
		univ, err := findUniversity(input.UnivID)
		if err != nil {
			return nil, humaError(err)
		}
		group, err := findGroup(univ, input.GroupID)
		if err != nil {
			return nil, humaError(err)
		}
		if _, err := findTimetable(group, input.AdeResources); err != nil {
			return nil, humaError(err)
		}
		events, err := fetchEvents(univ, input.AdeResources)
		if err != nil {
			return nil, humaError(err)
		}
		return &eventListOutput{Body: events}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "list-rooms",
		Method:      http.MethodGet,
		Path:        "/v3/univs/{univId}/rooms",
		Summary:     "List rooms for a university",
		Tags:        []string{"Rooms"},
	}, func(ctx context.Context, input *univInput) (*roomListOutput, error) {
		univ, err := findUniversity(input.UnivID)
		if err != nil {
			return nil, humaError(err)
		}
		firstDate, lastDate := ade.GetAcademicYearDates(time.Now())
		resp := make([]roomResponse, 0, len(univ.Rooms))
		for i := range univ.Rooms {
			resp = append(resp, buildRoomResponse(univ, &univ.Rooms[i], firstDate, lastDate))
		}
		return &roomListOutput{Body: resp}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-room-metadata",
		Method:      http.MethodGet,
		Path:        "/v3/univs/{univId}/rooms/{adeResources}",
		Summary:     "Get room metadata",
		Tags:        []string{"Rooms"},
	}, func(ctx context.Context, input *univAdeInput) (*roomOutput, error) {
		univ, err := findUniversity(input.UnivID)
		if err != nil {
			return nil, humaError(err)
		}
		room, err := findRoom(univ, input.AdeResources)
		if err != nil {
			return nil, humaError(err)
		}
		firstDate, lastDate := ade.GetAcademicYearDates(time.Now())
		return &roomOutput{Body: buildRoomResponse(univ, room, firstDate, lastDate)}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-room-events",
		Method:      http.MethodGet,
		Path:        "/v3/univs/{univId}/rooms/{adeResources}/events",
		Summary:     "Get room events (JSON)",
		Tags:        []string{"Rooms"},
	}, func(ctx context.Context, input *univAdeInput) (*eventListOutput, error) {
		univ, err := findUniversity(input.UnivID)
		if err != nil {
			return nil, humaError(err)
		}
		if _, err := findRoom(univ, input.AdeResources); err != nil {
			return nil, humaError(err)
		}
		events, err := fetchEvents(univ, input.AdeResources)
		if err != nil {
			return nil, humaError(err)
		}
		return &eventListOutput{Body: events}, nil
	})
}
