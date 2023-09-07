package api

import (
	"EtuEDT-Go/cache"
	"EtuEDT-Go/config"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"slices"
	"strconv"
	"time"
)

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

func v2Router(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		var universitiesResponse []UniversityResponse
		for _, university := range config.AppConfig.Universities {
			universitiesResponse = append(universitiesResponse, createUniversityResponse(&university))
		}
		return c.JSON(universitiesResponse)
	})

	router.Get("/:numUniv", func(c *fiber.Ctx) error {
		university, statusCode, err := getUniversityFromParam(c)
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		var timetablesResponse []TimetableResponse
		for _, timetable := range university.Timetables {
			timetablesResponse = append(timetablesResponse, createTimetableResponse(university, &timetable))
		}

		return c.JSON(timetablesResponse)
	})

	router.Get("/:numUniv/:adeResources/:format?", func(c *fiber.Ctx) error {
		university, statusCode, err := getUniversityFromParam(c)
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		timetable, statusCode, err := getTimetableFromParam(c, university)
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}

		timetableCache, ok := cache.GetTimetableByIds(university.NumUniv, timetable.AdeResources)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no cache found for this timetable"})
		}

		format := c.Params("format")
		if len(format) == 0 {
			return c.JSON(createTimetableResponse(university, timetable))
		} else if format == "json" {
			return c.JSON(timetableCache.Json)
		} else if format == "ics" {
			c.Set("Content-Type", "text/calendar")
			return c.SendString(timetableCache.Ical)
		} else {
			return c.JSON(fiber.Map{
				"error": "invalid format",
			})
		}
	})
}

func getUniversityFromParam(c *fiber.Ctx) (*config.UniversityConfig, int, error) {
	numUniv, err := strconv.Atoi(c.Params("numUniv"))
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.New("invalid parameter")
	}
	universityIndex := slices.IndexFunc(config.AppConfig.Universities, func(c config.UniversityConfig) bool { return c.NumUniv == numUniv })
	if universityIndex < 0 {
		return nil, fiber.StatusNotFound, errors.New("university not found")
	}
	university := config.AppConfig.Universities[universityIndex]
	return &university, 0, nil
}

func getTimetableFromParam(c *fiber.Ctx, university *config.UniversityConfig) (*config.TimetableConfig, int, error) {
	adeResources, err := strconv.Atoi(c.Params("adeResources"))
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.New("invalid parameter")
	}
	timetableIndex := slices.IndexFunc(university.Timetables, func(c config.TimetableConfig) bool { return c.AdeResources == adeResources })
	if timetableIndex < 0 {
		return nil, fiber.StatusNotFound, errors.New("timetable not found")
	}
	timetable := university.Timetables[timetableIndex]
	return &timetable, 0, nil
}

func createUniversityResponse(university *config.UniversityConfig) UniversityResponse {
	return UniversityResponse{
		NumUniv:  university.NumUniv,
		NameUniv: university.NameUniv,
		AdeUniv:  university.AdeUniv,
	}
}

func createTimetableResponse(university *config.UniversityConfig, timetable *config.TimetableConfig) TimetableResponse {
	timetableCache, ok := cache.GetTimetableByIds(university.NumUniv, timetable.AdeResources)
	lastUpdate := time.Time{}
	if ok {
		lastUpdate = timetableCache.LastUpdate
	}
	return TimetableResponse{
		NumUniv:      university.NumUniv,
		NameUniv:     university.NameUniv,
		NameTT:       fmt.Sprintf("%d%s%s", timetable.NumYearTT, "A ", timetable.DescTT),
		DescTT:       timetable.DescTT,
		NumYearTT:    timetable.NumYearTT,
		AdeResources: timetable.AdeResources,
		AdeProjectId: timetable.AdeProjectId,
		LastUpdate:   lastUpdate,
	}
}
