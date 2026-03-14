package v3

import (
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/domain"

	"github.com/gofiber/fiber/v2"
)

func ListUniversities(c *fiber.Ctx) error {
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

func GetUniversity(c *fiber.Ctx) error {
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

func ListGroups(c *fiber.Ctx) error {
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

func ListTimetables(c *fiber.Ctx) error {
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

func GetTimetableMetadata(c *fiber.Ctx) error {
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

func GetTimetableEvents(c *fiber.Ctx) error {
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

func ListRooms(c *fiber.Ctx) error {
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

func GetRoomMetadata(c *fiber.Ctx) error {
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

func GetRoomEvents(c *fiber.Ctx) error {
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
