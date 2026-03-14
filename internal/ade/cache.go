package ade

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/api"

	ics "github.com/arran4/golang-ical"
	"golang.org/x/sync/singleflight"
)

type TimetableCache struct {
	AdeResources int         `json:"adeResources"`
	LastUpdate   *time.Time  `json:"lastUpdate"`
	Ical         string      `json:"calendar"`
	Events       []api.Event `json:"events"`
}

var cacheMap = make(map[string]TimetableCache)
var cacheMu sync.RWMutex
var sfGroup singleflight.Group

func GetTimetableByAdeResources(univID int, adeResources int) (TimetableCache, bool) {
	key := cacheKey(univID, adeResources)
	cacheMu.RLock()
	timetable, ok := cacheMap[key]
	cacheMu.RUnlock()
	return timetable, ok
}

func SetTimetableByAdeResources(univID int, adeResources int, ical string, events []api.Event) TimetableCache {
	key := cacheKey(univID, adeResources)
	now := time.Now()
	cacheMu.Lock()
	timetable := TimetableCache{
		AdeResources: adeResources,
		LastUpdate:   &now,
		Ical:         ical,
		Events:       events,
	}
	cacheMap[key] = timetable
	cacheMu.Unlock()
	return timetable
}

func cacheKey(univID int, adeResources int) string {
	return strconv.Itoa(univID) + "-" + strconv.Itoa(adeResources)
}

func FetchTimetable(univID int, adeBaseUrl string, adeResources int, adeProjectId int) (*ics.Calendar, error) {
	key := cacheKey(univID, adeResources)
	result, err, _ := sfGroup.Do(key, func() (interface{}, error) {
		firstDate, lastDate := GetAcademicYearDates(time.Now())
		fullUrl, err := BuildURL(adeBaseUrl, adeResources, adeProjectId, firstDate, lastDate)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
		if err != nil {
			return nil, err
		}

		body, err := makeRequest(req)
		if err != nil {
			return nil, err
		}

		ical, err := ics.ParseCalendar(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}

		return ical, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*ics.Calendar), nil
}
