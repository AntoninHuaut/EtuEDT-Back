package cache

import (
	"EtuEDT-Go/domain"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
)

type TimetableCache struct {
	AdeResources int                `json:"adeResources"`
	LastUpdate   time.Time          `json:"lastUpdate"`
	Ical         string             `json:"calendar"`
	Json         []domain.JsonEvent `json:"json"`
}

var cacheMap = make(map[string]TimetableCache)
var cacheMu sync.RWMutex

func GetTimetableByAdeResources(univID int, adeResources int) (TimetableCache, bool) {
	key := getKey(univID, adeResources)
	cacheMu.RLock()
	timetable, ok := cacheMap[key]
	cacheMu.RUnlock()
	return timetable, ok
}

func SetTimetableByAdeResources(univID int, adeResources int, ical string, json []domain.JsonEvent) TimetableCache {
	key := getKey(univID, adeResources)
	cacheMu.Lock()
	timetable := TimetableCache{
		AdeResources: adeResources,
		LastUpdate:   time.Now(),
		Ical:         ical,
		Json:         json,
	}
	cacheMap[key] = timetable
	cacheMu.Unlock()
	return timetable
}

func getKey(univID int, adeResources int) string {
	return strconv.Itoa(univID) + "-" + strconv.Itoa(adeResources)
}

func FetchTimetable(adeBaseUrl string, adeResources int, adeProjectId int) (*ics.Calendar, error) {
	firstDate, lastDate := domain.GetAcademicYearDates(time.Now())
	fullUrl := domain.BuildAdeUrl(adeBaseUrl, adeResources, adeProjectId, firstDate, lastDate)

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}

	body, err := MakeRequest(fmt.Sprintf("%d", adeResources), req)
	if err != nil {
		return nil, err
	}

	ical, err := ics.ParseCalendar(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return ical, nil
}
