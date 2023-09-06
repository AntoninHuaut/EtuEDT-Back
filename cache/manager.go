package cache

import (
	"strconv"
	"time"
)

type TimetableCache struct {
	NumUniv      int         `json:"numUniv"`
	AdeResources int         `json:"adeResources"`
	LastUpdate   time.Time   `json:"lastUpdate"`
	Ical         string      `json:"calendar"`
	Json         []JsonEvent `json:"json"`
}

var cache = make(map[string]TimetableCache)

func GetTimetableByIds(numUniv int, adeResources int) (TimetableCache, bool) {
	key := getKey(numUniv, adeResources)
	timetable, ok := cache[key]
	return timetable, ok
}

func SetTimetableByIds(numUniv int, adeResources int, ical string, json []JsonEvent) TimetableCache {
	key := getKey(numUniv, adeResources)
	timetable, ok := cache[key]
	if ok {
		timetable.LastUpdate = time.Now()
		timetable.Ical = ical
		timetable.Json = json
	} else {
		timetable = TimetableCache{
			NumUniv:      numUniv,
			AdeResources: adeResources,
			LastUpdate:   time.Now(),
			Ical:         ical,
			Json:         json,
		}
	}
	cache[key] = timetable
	return timetable
}

func getKey(numUniv int, adeResources int) string {
	return strconv.Itoa(numUniv) + "-" + strconv.Itoa(adeResources)
}
