package cache

import (
	"EtuEDT-Go/domain"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/go-co-op/gocron"
)

func StartScheduler() error {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(domain.AppConfig.RefreshMinutes).Minutes().Do(func() {
		timeStart := time.Now()
		log.Println("[Info] [Cron] Refreshing timetables")
		refreshTimetables(domain.AppConfig.Universities)
		log.Println("[Info] [Cron] Refresh done in " + time.Since(timeStart).String())
	})
	if err != nil {
		return err
	}

	scheduler.StartAsync()

	return nil
}

func refreshTimetables(universities []domain.UniversityConfig) {
	for _, university := range universities {
		for _, timetable := range university.Timetables {
			calendar, err := fetchTimetable(university, timetable)
			if err == nil {
				SetTimetableByIds(university.NumUniv, timetable.AdeResources, calendar.Serialize(), calendarToJson(calendar))
			}
		}
	}
}

func fetchTimetable(university domain.UniversityConfig, timetable domain.TimetableConfig) (*ics.Calendar, error) {
	firstDate := time.Now().AddDate(0, -4, 0).Format("2006-01-02")
	lastDate := time.Now().AddDate(0, 4, 0).Format("2006-01-02")
	req, err := http.NewRequest(http.MethodGet, university.AdeUniv, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("resources", strconv.Itoa(timetable.AdeResources))
	q.Add("projectId", strconv.Itoa(timetable.AdeProjectId))
	q.Add("calType", "ical")
	q.Add("firstDate", firstDate)
	q.Add("lastDate", lastDate)
	req.URL.RawQuery = q.Encode()

	body, err := MakeRequest(fmt.Sprintf("%d", timetable.AdeResources), req)
	if err != nil {
		return nil, err
	}

	ical, err := ics.ParseCalendar(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return ical, nil
}
