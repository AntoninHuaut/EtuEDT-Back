package api

import (
	"EtuEDT-Go/cache"
	"EtuEDT-Go/domain"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"slices"
	"strconv"
	"time"
)

func v2Router(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		var universitiesResponse []domain.UniversityResponse
		for _, university := range domain.AppConfig.Universities {
			universitiesResponse = append(universitiesResponse, createUniversityResponse(&university))
		}
		return c.JSON(universitiesResponse)
	})

	router.Get("/:numUniv", func(c *fiber.Ctx) error {
		university, statusCode, err := getUniversityFromParam(c)
		if err != nil {
			return c.Status(statusCode).JSON(domain.ErrorResponse{Error: err.Error()})
		}
		var timetablesResponse []domain.TimetableResponse
		for _, timetable := range university.Timetables {
			timetablesResponse = append(timetablesResponse, createTimetableResponse(university, &timetable))
		}

		return c.JSON(timetablesResponse)
	})

	router.Get("/:numUniv/:adeResources/:format?", func(c *fiber.Ctx) error {
		university, statusCode, err := getUniversityFromParam(c)
		if err != nil {
			return c.Status(statusCode).JSON(domain.ErrorResponse{Error: err.Error()})
		}
		timetable, statusCode, err := getTimetableFromParam(c, university)
		if err != nil {
			return c.Status(statusCode).JSON(domain.ErrorResponse{Error: err.Error()})
		}

		timetableCache, ok := cache.GetTimetableByIds(university.NumUniv, timetable.AdeResources)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "no cache found for this timetable"})
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
			return c.JSON(domain.ErrorResponse{
				Error: "invalid format",
			})
		}
	})
}

func getUniversityFromParam(c *fiber.Ctx) (*domain.UniversityConfig, int, error) {
	numUniv, err := strconv.Atoi(c.Params("numUniv"))
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.New("invalid parameter")
	}
	universityIndex := slices.IndexFunc(domain.AppConfig.Universities, func(c domain.UniversityConfig) bool { return c.NumUniv == numUniv })
	if universityIndex < 0 {
		return nil, fiber.StatusNotFound, errors.New("university not found")
	}
	university := domain.AppConfig.Universities[universityIndex]
	return &university, 0, nil
}

func getTimetableFromParam(c *fiber.Ctx, university *domain.UniversityConfig) (*domain.TimetableConfig, int, error) {
	adeResources, err := strconv.Atoi(c.Params("adeResources"))
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.New("invalid parameter")
	}
	timetableIndex := slices.IndexFunc(university.Timetables, func(c domain.TimetableConfig) bool { return c.AdeResources == adeResources })
	if timetableIndex < 0 {
		return nil, fiber.StatusNotFound, errors.New("timetable not found")
	}
	timetable := university.Timetables[timetableIndex]
	return &timetable, 0, nil
}

func createUniversityResponse(university *domain.UniversityConfig) domain.UniversityResponse {
	return domain.UniversityResponse{
		NumUniv:  university.NumUniv,
		NameUniv: university.NameUniv,
		AdeUniv:  university.AdeUniv,
	}
}

func createTimetableResponse(university *domain.UniversityConfig, timetable *domain.TimetableConfig) domain.TimetableResponse {
	timetableCache, ok := cache.GetTimetableByIds(university.NumUniv, timetable.AdeResources)
	lastUpdate := time.Time{}
	if ok {
		lastUpdate = timetableCache.LastUpdate
	}
	return domain.TimetableResponse{
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
