package v3

import (
	"slices"
	"strconv"

	"github.com/AntoninHuaut/EtuEDT-Back/cache"
	"github.com/AntoninHuaut/EtuEDT-Back/domain"

	"github.com/gofiber/fiber/v2"
)

func findUniversity(c *fiber.Ctx) (*domain.UniversityConfig, bool) {
	univId, err := strconv.Atoi(c.Params("univId"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "invalid univId parameter"})
		return nil, false
	}
	idx := slices.IndexFunc(domain.AppConfig.Universities, func(u domain.UniversityConfig) bool {
		return u.ID == univId
	})
	if idx < 0 {
		_ = c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse{Error: "university not found"})
		return nil, false
	}
	return &domain.AppConfig.Universities[idx], true
}

func findGroup(c *fiber.Ctx, univ *domain.UniversityConfig) (*domain.GroupConfig, bool) {
	groupId, err := strconv.Atoi(c.Params("groupId"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "invalid groupId parameter"})
		return nil, false
	}
	idx := slices.IndexFunc(univ.Groups, func(g domain.GroupConfig) bool {
		return g.ID == groupId
	})
	if idx < 0 {
		_ = c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse{Error: "group not found"})
		return nil, false
	}
	return &univ.Groups[idx], true
}

func findTimetable(c *fiber.Ctx, group *domain.GroupConfig) (*domain.TimetableConfig, bool) {
	adeResources, err := strconv.Atoi(c.Params("adeResources"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "invalid adeResources parameter"})
		return nil, false
	}
	idx := slices.IndexFunc(group.Timetables, func(tt domain.TimetableConfig) bool {
		return tt.AdeResources == adeResources
	})
	if idx < 0 {
		_ = c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse{Error: "timetable not found"})
		return nil, false
	}
	return &group.Timetables[idx], true
}

func findRoom(c *fiber.Ctx, univ *domain.UniversityConfig) (*domain.RoomConfig, bool) {
	adeResources, err := strconv.Atoi(c.Params("adeResources"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "invalid adeResources parameter"})
		return nil, false
	}
	idx := slices.IndexFunc(univ.Rooms, func(r domain.RoomConfig) bool {
		return r.AdeResources == adeResources
	})
	if idx < 0 {
		_ = c.Status(fiber.StatusNotFound).JSON(domain.ErrorResponse{Error: "room not found"})
		return nil, false
	}
	return &univ.Rooms[idx], true
}

func buildTimetableResponse(univ *domain.UniversityConfig, tt *domain.TimetableConfig, firstDate string, lastDate string) domain.TimetableResponse {
	cached, _ := cache.GetTimetableByAdeResources(univ.ID, tt.AdeResources)
	return domain.TimetableResponse{
		AdeResources: tt.AdeResources,
		AdeProjectId: univ.AdeProjectId,
		Year:         tt.Year,
		Label:        tt.Label,
		AdeUrl:       domain.BuildAdeUrl(univ.AdeUrl, tt.AdeResources, univ.AdeProjectId, firstDate, lastDate),
		LastUpdate:   cached.LastUpdate,
	}
}

func buildRoomResponse(univ *domain.UniversityConfig, room *domain.RoomConfig, firstDate string, lastDate string) domain.RoomResponse {
	cached, _ := cache.GetTimetableByAdeResources(univ.ID, room.AdeResources)
	return domain.RoomResponse{
		AdeResources: room.AdeResources,
		AdeProjectId: univ.AdeProjectId,
		Label:        room.Label,
		AdeUrl:       domain.BuildAdeUrl(univ.AdeUrl, room.AdeResources, univ.AdeProjectId, firstDate, lastDate),
		LastUpdate:   cached.LastUpdate,
	}
}

// serveEvents returns cached events in the format requested by the Accept header.
// Defaults to text/calendar (ICS). Use Accept: application/json for JSON format.
// Always tries to refresh from ADE on each request.
// If refresh fails, returns cached data when available.
func serveEvents(c *fiber.Ctx, univ *domain.UniversityConfig, adeResources int) error {
	timetableCache, ok := cache.GetTimetableByAdeResources(univ.ID, adeResources)

	calendar, err := cache.FetchTimetable(univ.ID, univ.AdeUrl, adeResources, univ.AdeProjectId)
	if err == nil {
		timetableCache = cache.SetTimetableByAdeResources(univ.ID, adeResources, calendar.Serialize(), cache.CalendarToJson(calendar))
		ok = true
	} else if !ok {
		return c.Status(fiber.StatusServiceUnavailable).JSON(domain.ErrorResponse{
			Error: "could not fetch timetable and no cache available, try again later",
		})
	}

	accept := c.Accepts("application/json", "text/calendar")
	switch accept {
	case "application/json":
		return c.JSON(timetableCache.Json)
	default:
		c.Set("Content-Type", "text/calendar; charset=utf-8")
		return c.SendString(timetableCache.Ical)
	}
}
