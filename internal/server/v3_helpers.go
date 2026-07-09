package server

import (
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
	"github.com/danielgtaylor/huma/v2"
)

const (
	cacheTTL = 12 * time.Hour
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

func findUniversity(univID int) (*config.UniversityConfig, error) {
	idx := slices.IndexFunc(config.AppConfig.Universities, func(u config.UniversityConfig) bool {
		return u.ID == univID
	})
	if idx < 0 {
		return nil, errUniversityNotFound
	}
	return &config.AppConfig.Universities[idx], nil
}

func findGroup(univ *config.UniversityConfig, groupID int) (*config.GroupConfig, error) {
	idx := slices.IndexFunc(univ.Groups, func(g config.GroupConfig) bool {
		return g.ID == groupID
	})
	if idx < 0 {
		return nil, errGroupNotFound
	}
	return &univ.Groups[idx], nil
}

func findTimetable(group *config.GroupConfig, adeResources int) (*config.TimetableConfig, error) {
	idx := slices.IndexFunc(group.Timetables, func(tt config.TimetableConfig) bool {
		return tt.AdeResources == adeResources
	})
	if idx < 0 {
		return nil, errTimetableNotFound
	}
	return &group.Timetables[idx], nil
}

func findRoom(univ *config.UniversityConfig, adeResources int) (*config.RoomConfig, error) {
	idx := slices.IndexFunc(univ.Rooms, func(r config.RoomConfig) bool {
		return r.AdeResources == adeResources
	})
	if idx < 0 {
		return nil, errRoomNotFound
	}
	return &univ.Rooms[idx], nil
}

func buildTimetableResponse(univ *config.UniversityConfig, tt *config.TimetableConfig, now time.Time, firstDate time.Time, lastDate time.Time) timetableResponse {
	cached, _ := ade.GetTimetableByAdeResources(univ.ID, tt.AdeResources)
	projectId := univ.GetEffectiveProjectId(now)
	adeUrl, _ := ade.BuildURL(univ.AdeUrl, tt.AdeResources, projectId, firstDate, lastDate)
	return timetableResponse{
		AdeResources: tt.AdeResources,
		AdeProjectId: projectId,
		Year:         tt.Year,
		Label:        tt.Label,
		AdeUrl:       adeUrl,
		LastUpdate:   cached.LastUpdate,
	}
}

func buildRoomResponse(univ *config.UniversityConfig, room *config.RoomConfig, now time.Time, firstDate time.Time, lastDate time.Time) roomResponse {
	cached, _ := ade.GetTimetableByAdeResources(univ.ID, room.AdeResources)
	projectId := univ.GetEffectiveProjectId(now)
	adeUrl, _ := ade.BuildURL(univ.AdeUrl, room.AdeResources, projectId, firstDate, lastDate)
	return roomResponse{
		AdeResources: room.AdeResources,
		AdeProjectId: projectId,
		Label:        room.Label,
		AdeUrl:       adeUrl,
		LastUpdate:   cached.LastUpdate,
	}
}

func fetchEvents(univ *config.UniversityConfig, adeResources int) ([]ade.Event, error) {
	cached, ok := ade.GetTimetableByAdeResources(univ.ID, adeResources)
	if ok && cached.LastUpdate != nil && time.Since(*cached.LastUpdate) < cacheTTL {
		return cached.Events, nil
	}

	now := time.Now()
	projectID := univ.GetEffectiveProjectId(now)
	calendar, err := ade.FetchTimetable(univ.ID, univ.AdeUrl, adeResources, projectID, univ.GetSplitMonth())
	if err == nil {
		fresh := ade.SetTimetableByAdeResources(univ.ID, adeResources, ade.CalendarToEvents(calendar))
		return fresh.Events, nil
	}

	if ok {
		slog.Warn("upstream fetch failed, serving stale cache", "univId", univ.ID, "adeResources", adeResources, "err", err)
		return cached.Events, nil
	}

	return nil, errTimetableUnavailable
}
