package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
	"github.com/danielgtaylor/huma/v2"
)

// setupAppConfig sets config.AppConfig to a known test fixture.
func setupAppConfig() {
	config.AppConfig = config.Config{
		Universities: []config.UniversityConfig{
			{
				ID:           1,
				Name:         "Test University",
				AdeUrl:       "https://ade.example.com",
				AdeProjectId: 42,
				Rooms: []config.RoomConfig{
					{AdeResources: 10, Label: "Room A"},
				},
				Groups: []config.GroupConfig{
					{
						ID:   1,
						Name: "Group 1",
						Timetables: []config.TimetableConfig{
							{AdeResources: 20, Year: 2025, Label: "Timetable A"},
						},
					},
				},
			},
		},
	}
}

// --- humaError tests ---

func TestHumaError_UniversityNotFound_Returns404(t *testing.T) {
	he := humaError(errUniversityNotFound)
	var se huma.StatusError
	if !errors.As(he, &se) {
		t.Fatalf("expected huma.StatusError, got %T", he)
	}
	if se.GetStatus() != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", se.GetStatus(), http.StatusNotFound)
	}
}

func TestHumaError_GroupNotFound_Returns404(t *testing.T) {
	he := humaError(errGroupNotFound)
	var se huma.StatusError
	if !errors.As(he, &se) {
		t.Fatalf("expected huma.StatusError, got %T", he)
	}
	if se.GetStatus() != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", se.GetStatus(), http.StatusNotFound)
	}
}

func TestHumaError_TimetableNotFound_Returns404(t *testing.T) {
	he := humaError(errTimetableNotFound)
	var se huma.StatusError
	if !errors.As(he, &se) {
		t.Fatalf("expected huma.StatusError, got %T", he)
	}
	if se.GetStatus() != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", se.GetStatus(), http.StatusNotFound)
	}
}

func TestHumaError_RoomNotFound_Returns404(t *testing.T) {
	he := humaError(errRoomNotFound)
	var se huma.StatusError
	if !errors.As(he, &se) {
		t.Fatalf("expected huma.StatusError, got %T", he)
	}
	if se.GetStatus() != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", se.GetStatus(), http.StatusNotFound)
	}
}

func TestHumaError_TimetableUnavailable_Returns503(t *testing.T) {
	he := humaError(errTimetableUnavailable)
	var se huma.StatusError
	if !errors.As(he, &se) {
		t.Fatalf("expected huma.StatusError, got %T", he)
	}
	if se.GetStatus() != http.StatusServiceUnavailable {
		t.Errorf("status: got %d, want %d", se.GetStatus(), http.StatusServiceUnavailable)
	}
}

func TestHumaError_GenericError_Returns500(t *testing.T) {
	he := humaError(errors.New("something unexpected"))
	var se huma.StatusError
	if !errors.As(he, &se) {
		t.Fatalf("expected huma.StatusError, got %T", he)
	}
	if se.GetStatus() != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", se.GetStatus(), http.StatusInternalServerError)
	}
}

// --- findUniversity tests ---

func TestFindUniversity_Found(t *testing.T) {
	setupAppConfig()

	u, err := findUniversity(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID != 1 || u.Name != "Test University" {
		t.Errorf("unexpected university: %+v", u)
	}
}

func TestFindUniversity_NotFound(t *testing.T) {
	setupAppConfig()

	_, err := findUniversity(999)
	if !errors.Is(err, errUniversityNotFound) {
		t.Errorf("expected errUniversityNotFound, got %v", err)
	}
}

// --- findGroup tests ---

func TestFindGroup_Found(t *testing.T) {
	setupAppConfig()

	univ := &config.AppConfig.Universities[0]
	g, err := findGroup(univ, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.ID != 1 || g.Name != "Group 1" {
		t.Errorf("unexpected group: %+v", g)
	}
}

func TestFindGroup_NotFound(t *testing.T) {
	setupAppConfig()

	univ := &config.AppConfig.Universities[0]
	_, err := findGroup(univ, 999)
	if !errors.Is(err, errGroupNotFound) {
		t.Errorf("expected errGroupNotFound, got %v", err)
	}
}

// --- findTimetable tests ---

func TestFindTimetable_Found(t *testing.T) {
	setupAppConfig()

	group := &config.AppConfig.Universities[0].Groups[0]
	tt, err := findTimetable(group, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tt.AdeResources != 20 || tt.Label != "Timetable A" {
		t.Errorf("unexpected timetable: %+v", tt)
	}
}

func TestFindTimetable_NotFound(t *testing.T) {
	setupAppConfig()

	group := &config.AppConfig.Universities[0].Groups[0]
	_, err := findTimetable(group, 999)
	if !errors.Is(err, errTimetableNotFound) {
		t.Errorf("expected errTimetableNotFound, got %v", err)
	}
}

// --- findRoom tests ---

func TestFindRoom_Found(t *testing.T) {
	setupAppConfig()

	univ := &config.AppConfig.Universities[0]
	r, err := findRoom(univ, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.AdeResources != 10 || r.Label != "Room A" {
		t.Errorf("unexpected room: %+v", r)
	}
}

func TestFindRoom_NotFound(t *testing.T) {
	setupAppConfig()

	univ := &config.AppConfig.Universities[0]
	_, err := findRoom(univ, 999)
	if !errors.Is(err, errRoomNotFound) {
		t.Errorf("expected errRoomNotFound, got %v", err)
	}
}

// --- fetchEvents tests (using pre-warmed cache to avoid real HTTP) ---

// stubUpstreamServer starts an httptest.Server that always returns the given
// status code, wires ade.HTTPClient to use it, and zeros ade.InitialBackoff
// so that retries complete instantly. Both are restored after the test.
func stubUpstreamServer(t *testing.T, statusCode int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "stub", statusCode)
	}))
	t.Cleanup(srv.Close)

	origClient := ade.HTTPClient
	ade.HTTPClient = srv.Client()
	t.Cleanup(func() { ade.HTTPClient = origClient })

	origBackoff := ade.InitialBackoff
	ade.InitialBackoff = 0
	t.Cleanup(func() { ade.InitialBackoff = origBackoff })

	return srv
}

func TestFetchEvents_ReturnsCachedEventsWhenFresh(t *testing.T) {
	setupAppConfig()

	univ := &config.AppConfig.Universities[0]
	const adeRes = 21 // unique key, isolated from other fetchEvents tests
	events := []ade.Event{{Title: "Cached Maths"}}
	ade.SetTimetableByAdeResources(univ.ID, adeRes, events)

	result, err := fetchEvents(univ, adeRes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Title != "Cached Maths" {
		t.Errorf("expected cached event, got %v", result)
	}
}

func TestFetchEvents_ReturnsCacheNotExpiredAfterJustSet(t *testing.T) {
	setupAppConfig()

	univ := &config.AppConfig.Universities[0]
	const adeRes = 22 // unique key, isolated from other fetchEvents tests
	events := []ade.Event{{Title: "Fresh Event"}}
	ade.SetTimetableByAdeResources(univ.ID, adeRes, events)

	// Call twice; second call must also return from cache (no upstream hit)
	result1, err1 := fetchEvents(univ, adeRes)
	result2, err2 := fetchEvents(univ, adeRes)

	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected errors: %v, %v", err1, err2)
	}
	if len(result1) != 1 || len(result2) != 1 {
		t.Errorf("both calls should return the cached event")
	}
	if result1[0].Title != result2[0].Title {
		t.Errorf("second call returned different data: %q vs %q", result1[0].Title, result2[0].Title)
	}
}

func TestFetchEvents_ReturnsTimetableUnavailableWhenNoCache(t *testing.T) {
	setupAppConfig()

	// Point upstream at a local server returning 500 with zero backoff
	// so retries are instant and no real network calls are made.
	srv := stubUpstreamServer(t, http.StatusInternalServerError)

	univ := &config.AppConfig.Universities[0]
	univ.AdeUrl = srv.URL
	const nonCachedResource = 99999

	_, err := fetchEvents(univ, nonCachedResource)
	if !errors.Is(err, errTimetableUnavailable) {
		t.Errorf("expected errTimetableUnavailable when no cache and upstream fails, got %v", err)
	}
}

// --- buildTimetableResponse tests ---

func TestBuildTimetableResponse_FieldsSet(t *testing.T) {
	setupAppConfig()

	univ := &config.AppConfig.Universities[0]
	tt := &univ.Groups[0].Timetables[0]
	firstDate := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	lastDate := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)

	resp := buildTimetableResponse(univ, tt, firstDate, lastDate)

	if resp.AdeResources != tt.AdeResources {
		t.Errorf("AdeResources: got %d, want %d", resp.AdeResources, tt.AdeResources)
	}
	if resp.AdeProjectId != univ.AdeProjectId {
		t.Errorf("AdeProjectId: got %d, want %d", resp.AdeProjectId, univ.AdeProjectId)
	}
	if resp.Year != tt.Year {
		t.Errorf("Year: got %d, want %d", resp.Year, tt.Year)
	}
	if resp.Label != tt.Label {
		t.Errorf("Label: got %q, want %q", resp.Label, tt.Label)
	}
	if resp.AdeUrl == "" {
		t.Errorf("AdeUrl should not be empty")
	}
}

// --- buildRoomResponse tests ---

func TestBuildRoomResponse_FieldsSet(t *testing.T) {
	setupAppConfig()

	univ := &config.AppConfig.Universities[0]
	room := &univ.Rooms[0]
	firstDate := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	lastDate := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)

	resp := buildRoomResponse(univ, room, firstDate, lastDate)

	if resp.AdeResources != room.AdeResources {
		t.Errorf("AdeResources: got %d, want %d", resp.AdeResources, room.AdeResources)
	}
	if resp.AdeProjectId != univ.AdeProjectId {
		t.Errorf("AdeProjectId: got %d, want %d", resp.AdeProjectId, univ.AdeProjectId)
	}
	if resp.Label != room.Label {
		t.Errorf("Label: got %q, want %q", resp.Label, room.Label)
	}
	if resp.AdeUrl == "" {
		t.Errorf("AdeUrl should not be empty")
	}
}
