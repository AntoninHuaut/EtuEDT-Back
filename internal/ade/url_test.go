package ade

import (
	"net/url"
	"testing"
	"time"
)

func TestGetAcademicYearDates_BeforeAugust(t *testing.T) {
	// July is before August → rolls back one year
	now := time.Date(2025, time.July, 15, 12, 0, 0, 0, time.UTC)
	first, last := GetAcademicYearDates(now)

	wantFirst := time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC)
	wantLast := time.Date(2025, time.July, 31, 0, 0, 0, 0, time.UTC)

	if !first.Equal(wantFirst) {
		t.Errorf("first date: got %v, want %v", first, wantFirst)
	}
	if !last.Equal(wantLast) {
		t.Errorf("last date: got %v, want %v", last, wantLast)
	}
}

func TestGetAcademicYearDates_InAugust(t *testing.T) {
	// August is the start month → stays in current year
	now := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	first, last := GetAcademicYearDates(now)

	wantFirst := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	wantLast := time.Date(2026, time.July, 31, 0, 0, 0, 0, time.UTC)

	if !first.Equal(wantFirst) {
		t.Errorf("first date: got %v, want %v", first, wantFirst)
	}
	if !last.Equal(wantLast) {
		t.Errorf("last date: got %v, want %v", last, wantLast)
	}
}

func TestGetAcademicYearDates_AfterAugust(t *testing.T) {
	// November is after August → stays in current year
	now := time.Date(2025, time.November, 10, 9, 0, 0, 0, time.UTC)
	first, last := GetAcademicYearDates(now)

	wantFirst := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	wantLast := time.Date(2026, time.July, 31, 0, 0, 0, 0, time.UTC)

	if !first.Equal(wantFirst) {
		t.Errorf("first date: got %v, want %v", first, wantFirst)
	}
	if !last.Equal(wantLast) {
		t.Errorf("last date: got %v, want %v", last, wantLast)
	}
}

func TestGetAcademicYearDates_January(t *testing.T) {
	// January is before August → rolls back one year
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	first, last := GetAcademicYearDates(now)

	wantFirst := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	wantLast := time.Date(2026, time.July, 31, 0, 0, 0, 0, time.UTC)

	if !first.Equal(wantFirst) {
		t.Errorf("first date: got %v, want %v", first, wantFirst)
	}
	if !last.Equal(wantLast) {
		t.Errorf("last date: got %v, want %v", last, wantLast)
	}
}

func TestBuildURL_ContainsAllQueryParams(t *testing.T) {
	baseUrl := "https://ade.example.com/jsp/custom/modules/plannings/online.jsp"
	firstDate := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	lastDate := time.Date(2026, time.July, 31, 0, 0, 0, 0, time.UTC)

	result, err := BuildURL(baseUrl, 1177, 42, firstDate, lastDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("result is not a valid URL: %v", err)
	}
	q := parsed.Query()

	if got := q.Get("resources"); got != "1177" {
		t.Errorf("resources: got %q, want %q", got, "1177")
	}
	if got := q.Get("projectId"); got != "42" {
		t.Errorf("projectId: got %q, want %q", got, "42")
	}
	if got := q.Get("calType"); got != "ical" {
		t.Errorf("calType: got %q, want %q", got, "ical")
	}
	if got := q.Get("firstDate"); got != "2025-08-01" {
		t.Errorf("firstDate: got %q, want %q", got, "2025-08-01")
	}
	if got := q.Get("lastDate"); got != "2026-07-31" {
		t.Errorf("lastDate: got %q, want %q", got, "2026-07-31")
	}
}

func TestBuildURL_PreservesBaseURLPath(t *testing.T) {
	baseUrl := "https://ade.example.com/jsp/custom/modules/plannings/online.jsp"
	firstDate := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	lastDate := time.Date(2026, time.July, 31, 0, 0, 0, 0, time.UTC)

	result, err := BuildURL(baseUrl, 1177, 42, firstDate, lastDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("result is not a valid URL: %v", err)
	}

	if parsed.Scheme != "https" {
		t.Errorf("scheme: got %q, want %q", parsed.Scheme, "https")
	}
	if parsed.Host != "ade.example.com" {
		t.Errorf("host: got %q, want %q", parsed.Host, "ade.example.com")
	}
	if parsed.Path != "/jsp/custom/modules/plannings/online.jsp" {
		t.Errorf("path: got %q", parsed.Path)
	}
}

func TestBuildURL_PreservesExistingQueryParams(t *testing.T) {
	baseUrl := "https://ade.example.com/?existing=value"
	firstDate := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	lastDate := time.Date(2026, time.July, 31, 0, 0, 0, 0, time.UTC)

	result, err := BuildURL(baseUrl, 100, 10, firstDate, lastDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("result is not a valid URL: %v", err)
	}
	q := parsed.Query()

	if got := q.Get("existing"); got != "value" {
		t.Errorf("existing param: got %q, want %q", got, "value")
	}
	if got := q.Get("resources"); got != "100" {
		t.Errorf("resources: got %q, want %q", got, "100")
	}
}
