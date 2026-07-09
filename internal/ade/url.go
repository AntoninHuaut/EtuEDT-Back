package ade

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func GetAcademicYearDates(now time.Time, splitMonth int) (time.Time, time.Time) {
	year := now.Year()
	if now.Month() < time.Month(splitMonth) {
		year--
	}
	first := time.Date(year, time.Month(splitMonth), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(1, 0, -1)
	return first, last
}

func BuildURL(baseUrl string, adeResources int, adeProjectId int, firstDate time.Time, lastDate time.Time) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("failed to parse ADE base URL %q: %w", baseUrl, err)
	}
	const dateFmt = "2006-01-02"
	q := u.Query()
	q.Set("resources", strconv.Itoa(adeResources))
	q.Set("projectId", strconv.Itoa(adeProjectId))
	q.Set("calType", "ical")
	q.Set("firstDate", firstDate.Format(dateFmt))
	q.Set("lastDate", lastDate.Format(dateFmt))
	u.RawQuery = q.Encode()
	return u.String(), nil
}
