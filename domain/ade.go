package domain

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
)

// GetAcademicYearDates returns the first and last date strings for the academic year
// containing the given time. Academic year runs August 1 to July 31.
//   - If month >= August: firstDate = currentYear-08-01, lastDate = nextYear-07-31
//   - If month < August:  firstDate = previousYear-08-01, lastDate = currentYear-07-31
func GetAcademicYearDates(now time.Time) (firstDate string, lastDate string) {
	year := now.Year()
	if now.Month() >= time.August {
		firstDate = fmt.Sprintf("%d-08-01", year)
		lastDate = fmt.Sprintf("%d-07-31", year+1)
	} else {
		firstDate = fmt.Sprintf("%d-08-01", year-1)
		lastDate = fmt.Sprintf("%d-07-31", year)
	}
	return
}

// BuildAdeUrl constructs the full ADE iCal URL with query parameters.
func BuildAdeUrl(baseUrl string, adeResources int, adeProjectId int, firstDate string, lastDate string) string {
	u, err := url.Parse(baseUrl)
	if err != nil {
		log.Printf("BuildAdeUrl: failed to parse base URL %q: %v", baseUrl, err)
		return baseUrl
	}
	q := u.Query()
	q.Set("resources", strconv.Itoa(adeResources))
	q.Set("projectId", strconv.Itoa(adeProjectId))
	q.Set("calType", "ical")
	q.Set("firstDate", firstDate)
	q.Set("lastDate", lastDate)
	u.RawQuery = q.Encode()
	return u.String()
}
