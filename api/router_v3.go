package api

import (
	"EtuEDT-Go/cache"
	"EtuEDT-Go/domain"
	"slices"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func v3Router(router fiber.Router) {
	univ := router.Group("/univ")

	// GET /v3/univ — list universities
	univ.Get("/", listUniversities)

	// GET /v3/univ/:univId — university detail
	univ.Get("/:univId", getUniversity)

	// GET /v3/univ/:univId/groups — list groups
	univ.Get("/:univId/groups", listGroups)

	// GET /v3/univ/:univId/groups/:groupId — list timetables in group
	univ.Get("/:univId/groups/:groupId", listTimetables)

	// GET /v3/univ/:univId/groups/:groupId/:adeResources — timetable metadata
	univ.Get("/:univId/groups/:groupId/:adeResources", getTimetableMetadata)

	// GET /v3/univ/:univId/groups/:groupId/:adeResources/events — timetable events
	univ.Get("/:univId/groups/:groupId/:adeResources/events", getTimetableEvents)

	// GET /v3/univ/:univId/rooms — list rooms
	univ.Get("/:univId/rooms", listRooms)

	// GET /v3/univ/:univId/rooms/:adeResources — room metadata
	univ.Get("/:univId/rooms/:adeResources", getRoomMetadata)

	// GET /v3/univ/:univId/rooms/:adeResources/events — room events
	univ.Get("/:univId/rooms/:adeResources/events", getRoomEvents)
}

// --- Handlers ---

func listUniversities(c *fiber.Ctx) error {
	resp := make([]domain.UniversityResponse, 0, len(domain.AppConfig.Universities))
	for _, u := range domain.AppConfig.Universities {
		resp = append(resp, domain.UniversityResponse{
			ID:     u.ID,
			Name:   u.Name,
			AdeUrl: u.AdeUrl,
		})
	}
	return c.JSON(resp)
}

func getUniversity(c *fiber.Ctx) error {
	univ, ok := findUniversity(c)
	if !ok {
		return nil
	}
	return c.JSON(domain.UniversityResponse{
		ID:     univ.ID,
		Name:   univ.Name,
		AdeUrl: univ.AdeUrl,
	})
}

func listGroups(c *fiber.Ctx) error {
	univ, ok := findUniversity(c)
	if !ok {
		return nil
	}
	resp := make([]domain.GroupResponse, 0, len(univ.Groups))
	for _, g := range univ.Groups {
		resp = append(resp, domain.GroupResponse{
			ID:   g.ID,
			Name: g.Name,
		})
	}
	return c.JSON(resp)
}

func listTimetables(c *fiber.Ctx) error {
	univ, ok := findUniversity(c)
	if !ok {
		return nil
	}
	group, ok := findGroup(c, univ)
	if !ok {
		return nil
	}
	firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
	resp := make([]domain.TimetableResponse, 0, len(group.Timetables))
	for i := range group.Timetables {
		resp = append(resp, buildTimetableResponse(univ, &group.Timetables[i], firstDate, lastDate))
	}
	return c.JSON(resp)
}

func getTimetableMetadata(c *fiber.Ctx) error {
	univ, ok := findUniversity(c)
	if !ok {
		return nil
	}
	group, ok := findGroup(c, univ)
	if !ok {
		return nil
	}
	tt, ok := findTimetable(c, group)
	if !ok {
		return nil
	}
	firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
	return c.JSON(buildTimetableResponse(univ, tt, firstDate, lastDate))
}

func getTimetableEvents(c *fiber.Ctx) error {
	univ, ok := findUniversity(c)
	if !ok {
		return nil
	}
	group, ok := findGroup(c, univ)
	if !ok {
		return nil
	}
	tt, ok := findTimetable(c, group)
	if !ok {
		return nil
	}
	return serveEvents(c, univ, tt.AdeResources)
}

func listRooms(c *fiber.Ctx) error {
	univ, ok := findUniversity(c)
	if !ok {
		return nil
	}
	firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
	resp := make([]domain.RoomResponse, 0, len(univ.Rooms))
	for i := range univ.Rooms {
		resp = append(resp, buildRoomResponse(univ, &univ.Rooms[i], firstDate, lastDate))
	}
	return c.JSON(resp)
}

func getRoomMetadata(c *fiber.Ctx) error {
	univ, ok := findUniversity(c)
	if !ok {
		return nil
	}
	room, ok := findRoom(c, univ)
	if !ok {
		return nil
	}
	firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
	return c.JSON(buildRoomResponse(univ, room, firstDate, lastDate))
}

func getRoomEvents(c *fiber.Ctx) error {
	univ, ok := findUniversity(c)
	if !ok {
		return nil
	}
	room, ok := findRoom(c, univ)
	if !ok {
		return nil
	}
	return serveEvents(c, univ, room.AdeResources)
}

// --- Content Negotiation ---

// serveEvents returns cached events in the format requested by the Accept header.
// Defaults to JSON. Use Accept: text/calendar for iCal format.
// Always tries to refresh from ADE on each request.
// If refresh fails, returns cached data when available.
func serveEvents(c *fiber.Ctx, univ *domain.UniversityConfig, adeResources int) error {
	timetableCache, ok := cache.GetTimetableByAdeResources(univ.ID, adeResources)

	calendar, err := cache.FetchTimetable(univ.AdeUrl, adeResources, univ.AdeProjectId)
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
	case "text/calendar":
		c.Set("Content-Type", "text/calendar")
		return c.SendString(timetableCache.Ical)
	default:
		return c.JSON(timetableCache.Json)
	}
}

// --- Lookup Helpers ---

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
