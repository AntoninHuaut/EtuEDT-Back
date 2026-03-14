package api

import (
	"EtuEDT-Go/domain"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func v2Router(router fiber.Router) {
	router.Get("/", legacyV2ICSMigrationNotice)
	router.Get("/*", legacyV2ICSMigrationNotice)
}

func legacyV2ICSMigrationNotice(c *fiber.Ctx) error {
	firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
	startAt, err := time.Parse("2006-01-02", firstDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "failed to build legacy calendar start date"})
	}

	endAtInclusive, err := time.Parse("2006-01-02", lastDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "failed to build legacy calendar end date"})
	}

	// iCal DTEND for all-day events is exclusive, so we add one day.
	endAtExclusive := endAtInclusive.Add(24 * time.Hour)
	ical := buildLegacyV2MigrationICS(startAt, endAtExclusive)

	c.Set("Content-Type", "text/calendar; charset=utf-8")
	return c.SendString(ical)
}

func buildLegacyV2MigrationICS(startAt time.Time, endAt time.Time) string {
	dtStamp := time.Now().UTC().Format("20060102T150405Z")
	dtStart := startAt.Format("20060102")
	dtEnd := endAt.Format("20060102")

	summary := icsEscape("Action requise : mettez a jour votre URL ICS")
	description := icsEscape("Cette URL v2 est obsolete. Mettez a jour la synchronisation de votre agenda (Google/Microsoft/Autres) vers l'URL ADE directe.")
	location := icsEscape("ADE")

	return strings.Join([]string{
		"BEGIN:VCALENDAR",
		"PRODID:-//EtuEDT//Legacy V2 Migration Notice//FR",
		"VERSION:2.0",
		"CALSCALE:GREGORIAN",
		"METHOD:PUBLISH",
		"BEGIN:VEVENT",
		fmt.Sprintf("UID:legacy-v2-migration-%s@etuedt", dtStart),
		fmt.Sprintf("DTSTAMP:%s", dtStamp),
		fmt.Sprintf("DTSTART;VALUE=DATE:%s", dtStart),
		fmt.Sprintf("DTEND;VALUE=DATE:%s", dtEnd),
		fmt.Sprintf("SUMMARY:%s", summary),
		fmt.Sprintf("DESCRIPTION:%s", description),
		fmt.Sprintf("LOCATION:%s", location),
		"STATUS:CONFIRMED",
		"TRANSP:TRANSPARENT",
		"END:VEVENT",
		"END:VCALENDAR",
		"",
	}, "\r\n")
}

func icsEscape(input string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		";", "\\;",
		",", "\\,",
		"\n", "\\n",
	)
	return replacer.Replace(input)
}
