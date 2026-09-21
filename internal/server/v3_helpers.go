package server

import (
	"log/slog"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/sync/errgroup"
)

const (
	cacheTTL = 12 * time.Hour
)

var (
	errUniversityNotFound   = huma.NewError(http.StatusNotFound, "university not found")
	errGroupNotFound        = huma.NewError(http.StatusNotFound, "group not found")
	errCampusNotFound       = huma.NewError(http.StatusNotFound, "campus not found")
	errTimetableNotFound    = huma.NewError(http.StatusNotFound, "timetable not found")
	errRoomNotFound         = huma.NewError(http.StatusNotFound, "room not found")
	errTimetableUnavailable = huma.NewError(http.StatusServiceUnavailable, "could not fetch timetable and no cache available, try again later")
	errCampusEmpty          = huma.NewError(http.StatusNotFound, "campus is empty")
	errEndBeforeStart       = huma.NewError(http.StatusBadRequest, "end before start")
)

func findUniversity(univID int) (*config.UniversityConfig, error) {
	idx := slices.IndexFunc(config.AppConfig.Universities, func(u config.UniversityConfig) bool {
		return u.ID == univID
	})
	if idx < 0 {
		return nil, errUniversityNotFound
	}
	return &config.AppConfig.Universities[idx], nil
}

func findCampus(univ *config.UniversityConfig, campusID int) (*config.CampusConfig, error) {
	idx := slices.IndexFunc(univ.Campuses, func(g config.CampusConfig) bool {
		return g.ID == campusID
	})
	if idx < 0 {
		return nil, errCampusNotFound
	}
	return &univ.Campuses[idx], nil
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

func findCampusRooms(univ *config.UniversityConfig, campusID int) ([]config.RoomConfig, error) {
	var campusRooms []config.RoomConfig
	for _, r := range univ.Rooms {
		if r.CampusID != nil && *r.CampusID == campusID {
			campusRooms = append(campusRooms, r)
		}
	}
	if len(campusRooms) == 0 {
		return nil, errCampusEmpty
	}
	return campusRooms, nil
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
	campusID := -1
	if room.CampusID != nil {
		campusID = *room.CampusID
	}
	return roomResponse{
		AdeResources: room.AdeResources,
		AdeProjectId: projectId,
		Label:        room.Label,
		AdeUrl:       adeUrl,
		LastUpdate:   cached.LastUpdate,
		CampusID:     campusID,
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

func findFreeRoom(univ *config.UniversityConfig, input *freeRoomsInput) (*roomListOutput, error) {
	now := time.Now()
	start := now
	if !input.Start.IsZero() {
		start = input.Start
	}

	end := start.Add(1 * time.Hour)
	if !input.End.IsZero() {
		end = input.End
	}

	if !start.Before(end) {
		return nil, errEndBeforeStart
	}

	firstDate, lastDate := ade.GetAcademicYearDates(now, univ.GetSplitMonth())

	freeRooms := make([]roomResponse, 0)
	var mu sync.Mutex

	var eg errgroup.Group

	for _, room := range univ.Rooms {
		if input.CampusID != 0 {
			if room.CampusID == nil || *room.CampusID != input.CampusID {
				continue
			}
		}

		eg.Go(func() error {
			events, err := fetchEvents(univ, room.AdeResources)
			if err != nil {
				return nil
			}

			isFree := true
			for _, event := range events {
				if event.Start.Before(end) && event.End.After(start) {
					isFree = false
					break
				}
			}

			if isFree {
				resp := buildRoomResponse(univ, &room, now, firstDate, lastDate)

				mu.Lock()
				freeRooms = append(freeRooms, resp)
				mu.Unlock()
			}

			return nil
		})
	}

	grpErr := eg.Wait()
	return &roomListOutput{Body: freeRooms}, grpErr
}
