package domain

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func GetAcademicYearDates(now time.Time) (time.Time, time.Time) {
	year := now.Year()
	if now.Month() < time.August {
		year--
	}
	first := time.Date(year, time.August, 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(year+1, time.July, 31, 0, 0, 0, 0, time.UTC)
	return first, last
}

func BuildAdeUrl(baseUrl string, adeResources int, adeProjectId int, firstDate time.Time, lastDate time.Time) (string, error) {
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
