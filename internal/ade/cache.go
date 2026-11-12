package ade

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
	"golang.org/x/sync/singleflight"
)

type TimetableCache struct {
	AdeResources int        `json:"adeResources"`
	LastUpdate   *time.Time `json:"lastUpdate"`
	Events       []Event    `json:"events"`
}

var cacheMap = make(map[string]TimetableCache)
var cacheMu sync.RWMutex
var sfGroup singleflight.Group

func GetTimetableByAdeResources(univID int, adeResources int, scope ...string) (TimetableCache, bool) {
	key := scopedCacheKey(univID, adeResources, scope)
	cacheMu.RLock()
	timetable, ok := cacheMap[key]
	cacheMu.RUnlock()
	return timetable, ok
}

func SetTimetableByAdeResources(univID int, adeResources int, events []Event, scope ...string) TimetableCache {
	key := scopedCacheKey(univID, adeResources, scope)
	cacheMu.Lock()
	timetable := TimetableCache{
		AdeResources: adeResources,
		LastUpdate:   new(time.Now()),
		Events:       events,
	}
	cacheMap[key] = timetable
	cacheMu.Unlock()
	return timetable
}

func cacheKey(univID int, adeResources int) string {
	return strconv.Itoa(univID) + "-" + strconv.Itoa(adeResources)
}

func scopedCacheKey(univID, resource int, scope []string) string {
	key := cacheKey(univID, resource)
	if len(scope) > 0 {
		key += ":" + scope[0]
	}
	return key
}

func FetchTimetable(univID int, adeBaseUrl string, adeResources int, adeProjectId int, splitMonth int) (*ics.Calendar, error) {
	return FetchTimetableContext(context.Background(), univID, adeBaseUrl, adeResources, adeProjectId, splitMonth)
}

func FetchTimetableContext(ctx context.Context, univID int, adeBaseUrl string, adeResources int, adeProjectId int, splitMonth int) (*ics.Calendar, error) {
	firstDate, lastDate := GetAcademicYearDates(time.Now(), splitMonth)
	fullUrl, err := BuildURL(adeBaseUrl, adeResources, adeProjectId, firstDate, lastDate)
	if err != nil {
		return nil, err
	}
	key := cacheKey(univID, adeResources) + ":" + fullUrl
	result, err, _ := sfGroup.Do(key, func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
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
