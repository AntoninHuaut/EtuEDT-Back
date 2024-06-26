package cache

import (
	"EtuEDT-Go/config"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/go-co-op/gocron"
)

const (
	httpTimeout = 30 * time.Second
)

type JsonEvent struct {
	Title       string    `json:"title"`
	Enseignant  string    `json:"enseignant"` // TODO temporary for compatibility with the old frontend
	Teacher     string    `json:"teacher"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
}

func StartScheduler() error {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(config.AppConfig.RefreshMinutes).Minutes().Do(func() {
		log.Println("Refreshing timetables")
		refreshTimetables(config.AppConfig.Universities)
	})
	if err != nil {
		return err
	}

	scheduler.StartAsync()

	return nil
}

func refreshTimetables(universities []config.UniversityConfig) {
	client := &http.Client{Timeout: httpTimeout}

	wg := sync.WaitGroup{}
	for _, university := range universities {
		for _, timetable := range university.Timetables {
			wg.Add(1)
			go func(university config.UniversityConfig, timetable config.TimetableConfig) {
				calendar, err := fetchTimetable(client, university, timetable)
				if err != nil {
					log.Printf("refreshTimetables: %v", err)
				} else {
					SetTimetableByIds(university.NumUniv, timetable.AdeResources, calendar.Serialize(), calendarToJson(calendar))
				}
				wg.Done()
			}(university, timetable)
		}
	}
	wg.Wait()
}

func fetchTimetable(client *http.Client, university config.UniversityConfig, timetable config.TimetableConfig) (*ics.Calendar, error) {
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

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var deferErr error
	defer func(Body io.ReadCloser) {
		deferErr = Body.Close()
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	ical, err := ics.ParseCalendar(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return ical, deferErr
}
