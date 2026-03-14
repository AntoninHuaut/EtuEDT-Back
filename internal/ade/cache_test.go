package ade

import (
	"sync"
	"testing"
	"time"
)

// resetCacheMap clears the global cache between tests.
func resetCacheMap() {
	cacheMu.Lock()
	cacheMap = make(map[string]TimetableCache)
	cacheMu.Unlock()
}

func TestGetTimetableByAdeResources_MissOnEmptyCache(t *testing.T) {
	resetCacheMap()

	_, ok := GetTimetableByAdeResources(1, 100)
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestSetAndGetTimetableByAdeResources_RoundTrip(t *testing.T) {
	resetCacheMap()

	events := []Event{{Title: "Maths", Location: "Room 1"}}
	SetTimetableByAdeResources(1, 100, events)

	entry, ok := GetTimetableByAdeResources(1, 100)
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if entry.AdeResources != 100 {
		t.Errorf("AdeResources: got %d, want 100", entry.AdeResources)
	}
	if len(entry.Events) != 1 || entry.Events[0].Title != "Maths" {
		t.Errorf("Events: got %v", entry.Events)
	}
}

func TestSetTimetableByAdeResources_SetsLastUpdate(t *testing.T) {
	resetCacheMap()

	before := time.Now()
	entry := SetTimetableByAdeResources(1, 200, nil)
	after := time.Now()

	if entry.LastUpdate == nil {
		t.Fatal("LastUpdate should not be nil")
	}
	if entry.LastUpdate.Before(before) || entry.LastUpdate.After(after) {
		t.Errorf("LastUpdate %v not within [%v, %v]", entry.LastUpdate, before, after)
	}
}

func TestGetTimetableByAdeResources_UnivIDIsolation(t *testing.T) {
	resetCacheMap()

	// Different univID → different cache entries
	events1 := []Event{{Title: "Physics"}}
	events2 := []Event{{Title: "Chemistry"}}
	SetTimetableByAdeResources(1, 100, events1)
	SetTimetableByAdeResources(2, 100, events2)

	entry1, ok1 := GetTimetableByAdeResources(1, 100)
	entry2, ok2 := GetTimetableByAdeResources(2, 100)

	if !ok1 || !ok2 {
		t.Fatal("both entries should be present")
	}
	if entry1.Events[0].Title == entry2.Events[0].Title {
		t.Errorf("entries should be isolated by univID but both have Title=%q", entry1.Events[0].Title)
	}
}

func TestGetTimetableByAdeResources_AdeResourcesIsolation(t *testing.T) {
	resetCacheMap()

	// Different adeResources → different cache entries
	events100 := []Event{{Title: "Math"}}
	events200 := []Event{{Title: "Science"}}
	SetTimetableByAdeResources(1, 100, events100)
	SetTimetableByAdeResources(1, 200, events200)

	e100, ok100 := GetTimetableByAdeResources(1, 100)
	e200, ok200 := GetTimetableByAdeResources(1, 200)

	if !ok100 || !ok200 {
		t.Fatal("both entries should be present")
	}
	if e100.Events[0].Title != "Math" {
		t.Errorf("adeResources=100: got %q, want %q", e100.Events[0].Title, "Math")
	}
	if e200.Events[0].Title != "Science" {
		t.Errorf("adeResources=200: got %q, want %q", e200.Events[0].Title, "Science")
	}
}

func TestSetTimetableByAdeResources_Overwrites(t *testing.T) {
	resetCacheMap()

	SetTimetableByAdeResources(1, 100, []Event{{Title: "Old"}})
	SetTimetableByAdeResources(1, 100, []Event{{Title: "New"}})

	entry, ok := GetTimetableByAdeResources(1, 100)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if entry.Events[0].Title != "New" {
		t.Errorf("expected overwrite: got %q, want %q", entry.Events[0].Title, "New")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	resetCacheMap()

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			SetTimetableByAdeResources(1, i, nil)
			GetTimetableByAdeResources(1, i)
		}(i)
	}
	wg.Wait()
	// No panic or race condition = pass
}
