package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/cache"
	"github.com/AntoninHuaut/EtuEDT-Back/domain"
	"github.com/danielgtaylor/huma/v2"
)

var (
	errUniversityNotFound   = errors.New("university not found")
	errGroupNotFound        = errors.New("group not found")
	errTimetableNotFound    = errors.New("timetable not found")
	errRoomNotFound         = errors.New("room not found")
	errTimetableUnavailable = errors.New("could not fetch timetable and no cache available, try again later")
)

func humaError(err error) error {
	switch {
	case errors.Is(err, errUniversityNotFound),
		errors.Is(err, errGroupNotFound),
		errors.Is(err, errTimetableNotFound),
		errors.Is(err, errRoomNotFound):
		return huma.Error404NotFound(err.Error())
	case errors.Is(err, errTimetableUnavailable):
		return huma.Error503ServiceUnavailable(err.Error())
	default:
		return huma.Error500InternalServerError(err.Error())
	}
}

func findUniversity(univID int) (*domain.UniversityConfig, error) {
	idx := slices.IndexFunc(domain.AppConfig.Universities, func(u domain.UniversityConfig) bool {
		return u.ID == univID
	})
	if idx < 0 {
		return nil, errUniversityNotFound
	}
	return &domain.AppConfig.Universities[idx], nil
}

func findGroup(univ *domain.UniversityConfig, groupID int) (*domain.GroupConfig, error) {
	idx := slices.IndexFunc(univ.Groups, func(g domain.GroupConfig) bool {
		return g.ID == groupID
	})
	if idx < 0 {
		return nil, errGroupNotFound
	}
	return &univ.Groups[idx], nil
}

func findTimetable(group *domain.GroupConfig, adeResources int) (*domain.TimetableConfig, error) {
	idx := slices.IndexFunc(group.Timetables, func(tt domain.TimetableConfig) bool {
		return tt.AdeResources == adeResources
	})
	if idx < 0 {
		return nil, errTimetableNotFound
	}
	return &group.Timetables[idx], nil
}

func findRoom(univ *domain.UniversityConfig, adeResources int) (*domain.RoomConfig, error) {
	idx := slices.IndexFunc(univ.Rooms, func(r domain.RoomConfig) bool {
		return r.AdeResources == adeResources
	})
	if idx < 0 {
		return nil, errRoomNotFound
	}
	return &univ.Rooms[idx], nil
}

func buildTimetableResponse(univ *domain.UniversityConfig, tt *domain.TimetableConfig, firstDate time.Time, lastDate time.Time) domain.TimetableResponse {
	cached, _ := cache.GetTimetableByAdeResources(univ.ID, tt.AdeResources)
	adeUrl, _ := domain.BuildAdeUrl(univ.AdeUrl, tt.AdeResources, univ.AdeProjectId, firstDate, lastDate)
	return domain.TimetableResponse{
		AdeResources: tt.AdeResources,
		AdeProjectId: univ.AdeProjectId,
		Year:         tt.Year,
		Label:        tt.Label,
		AdeUrl:       adeUrl,
		LastUpdate:   cached.LastUpdate,
	}
}

func buildRoomResponse(univ *domain.UniversityConfig, room *domain.RoomConfig, firstDate time.Time, lastDate time.Time) domain.RoomResponse {
	cached, _ := cache.GetTimetableByAdeResources(univ.ID, room.AdeResources)
	adeUrl, _ := domain.BuildAdeUrl(univ.AdeUrl, room.AdeResources, univ.AdeProjectId, firstDate, lastDate)
	return domain.RoomResponse{
		AdeResources: room.AdeResources,
		AdeProjectId: univ.AdeProjectId,
		Label:        room.Label,
		AdeUrl:       adeUrl,
		LastUpdate:   cached.LastUpdate,
	}
}

func fetchEvents(univ *domain.UniversityConfig, adeResources int) ([]domain.JsonEvent, error) {
	calendar, err := cache.FetchTimetable(univ.ID, univ.AdeUrl, adeResources, univ.AdeProjectId)
	if err == nil {
		fresh := cache.SetTimetableByAdeResources(univ.ID, adeResources, calendar.Serialize(), cache.CalendarToJson(calendar))
		return fresh.Json, nil
	}

	stale, ok := cache.GetTimetableByAdeResources(univ.ID, adeResources)
	if !ok {
		return nil, errTimetableUnavailable
	}

	slog.Warn("upstream fetch failed, serving stale cache", "univId", univ.ID, "adeResources", adeResources, "err", err)
	return stale.Json, nil
}

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

type universityListOutput struct {
	Body []domain.UniversityResponse
}

type universityOutput struct {
	Body domain.UniversityResponse
}

type groupListOutput struct {
	Body []domain.GroupResponse
}

type timetableListOutput struct {
	Body []domain.TimetableResponse
}

type timetableOutput struct {
	Body domain.TimetableResponse
}

type eventListOutput struct {
	Body []domain.JsonEvent
}

type roomListOutput struct {
	Body []domain.RoomResponse
}

type roomOutput struct {
	Body domain.RoomResponse
}

func registerV3Handlers(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-universities",
		Method:      http.MethodGet,
		Path:        "/v3/univs",
		Summary:     "List all universities",
		Tags:        []string{"Universities"},
	}, func(ctx context.Context, _ *struct{}) (*universityListOutput, error) {
		resp := make([]domain.UniversityResponse, 0, len(domain.AppConfig.Universities))
		for _, u := range domain.AppConfig.Universities {
			resp = append(resp, domain.UniversityResponse{ID: u.ID, Name: u.Name, AdeUrl: u.AdeUrl})
		}
		return &universityListOutput{Body: resp}, nil
	})

	huma.Register(api, huma.Operation{
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
		return &universityOutput{Body: domain.UniversityResponse{ID: univ.ID, Name: univ.Name, AdeUrl: univ.AdeUrl}}, nil
	})

	huma.Register(api, huma.Operation{
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
		resp := make([]domain.GroupResponse, 0, len(univ.Groups))
		for _, g := range univ.Groups {
			resp = append(resp, domain.GroupResponse{ID: g.ID, Name: g.Name})
		}
		return &groupListOutput{Body: resp}, nil
	})

	huma.Register(api, huma.Operation{
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
		firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
		resp := make([]domain.TimetableResponse, 0, len(group.Timetables))
		for i := range group.Timetables {
			resp = append(resp, buildTimetableResponse(univ, &group.Timetables[i], firstDate, lastDate))
		}
		return &timetableListOutput{Body: resp}, nil
	})

	huma.Register(api, huma.Operation{
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
		firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
		return &timetableOutput{Body: buildTimetableResponse(univ, tt, firstDate, lastDate)}, nil
	})

	huma.Register(api, huma.Operation{
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

	huma.Register(api, huma.Operation{
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
		firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
		resp := make([]domain.RoomResponse, 0, len(univ.Rooms))
		for i := range univ.Rooms {
			resp = append(resp, buildRoomResponse(univ, &univ.Rooms[i], firstDate, lastDate))
		}
		return &roomListOutput{Body: resp}, nil
	})

	huma.Register(api, huma.Operation{
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
		firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
		return &roomOutput{Body: buildRoomResponse(univ, room, firstDate, lastDate)}, nil
	})

	huma.Register(api, huma.Operation{
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
