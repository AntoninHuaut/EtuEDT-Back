package ade

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

// minimalIcal is a valid, minimal iCalendar with one VEVENT.
const minimalIcal = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:1@test
DTSTART:20251001T090000Z
DTEND:20251001T100000Z
SUMMARY:Test Event
END:VEVENT
END:VCALENDAR
`

// resetHTTPClient replaces the package-level HTTPClient with the given client
// and restores the original after the test.
func resetHTTPClient(t *testing.T, client *http.Client) {
	t.Helper()
	original := HTTPClient
	HTTPClient = client
	t.Cleanup(func() { HTTPClient = original })
}

// zeroBackoff sets InitialBackoff to 0 for the duration of the test so
// retries complete immediately. With BackOffDelay, 0 * 2^n = 0 for all n.
func zeroBackoff(t *testing.T) {
	t.Helper()
	original := InitialBackoff
	InitialBackoff = 0
	t.Cleanup(func() { InitialBackoff = original })
}

// resetSfGroup replaces the singleflight group so each test starts clean.
func resetSfGroup(t *testing.T) {
	t.Helper()
	sfGroup = singleflight.Group{}
	t.Cleanup(func() { sfGroup = singleflight.Group{} })
}

// TestFetchTimetable_SuccessOnValidResponse checks that FetchTimetable
// returns a parsed calendar when the server responds with valid iCal data.
func TestFetchTimetable_SuccessOnValidResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(minimalIcal))
	}))
	defer srv.Close()

	resetHTTPClient(t, srv.Client())
	zeroBackoff(t)
	resetSfGroup(t)

	cal, err := FetchTimetable(1, srv.URL, 100, 42, 7)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cal == nil {
		t.Fatal("expected non-nil calendar")
	}
	events := cal.Events()
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

// TestFetchTimetable_ReturnsErrorOnNon2xxStatus checks that a non-2xx
// response causes FetchTimetable to return an error (after retries).
func TestFetchTimetable_ReturnsErrorOnNon2xxStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	resetHTTPClient(t, srv.Client())
	zeroBackoff(t)
	resetSfGroup(t)

	_, err := FetchTimetable(2, srv.URL, 100, 42, 7)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status code 500") {
		t.Errorf("expected error to mention 'unexpected status code 500', got: %v", err)
	}
}

// TestFetchTimetable_RetriesBeforeFailing checks that the client retries
// exactly maxAttempts times before returning an error.
func TestFetchTimetable_RetriesBeforeFailing(t *testing.T) {
	var callCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer srv.Close()

	resetHTTPClient(t, srv.Client())
	zeroBackoff(t)
	resetSfGroup(t)

	_, err := FetchTimetable(3, srv.URL, 200, 42, 7)
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if int(callCount.Load()) != int(MaxAttempts) {
		t.Errorf("expected %d attempts, got %d", MaxAttempts, callCount.Load())
	}
}

// TestFetchTimetable_SucceedsAfterTransientFailure checks that FetchTimetable
// succeeds when the server fails the first two attempts then returns valid data.
func TestFetchTimetable_SucceedsAfterTransientFailure(t *testing.T) {
	var callCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := callCount.Add(1)
		if n < 3 {
			http.Error(w, "transient", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(minimalIcal))
	}))
	defer srv.Close()

	resetHTTPClient(t, srv.Client())
	zeroBackoff(t)
	resetSfGroup(t)

	cal, err := FetchTimetable(4, srv.URL, 300, 42, 7)
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if cal == nil {
		t.Fatal("expected non-nil calendar")
	}
	if int(callCount.Load()) != 3 {
		t.Errorf("expected 3 attempts, got %d", callCount.Load())
	}
}

// TestFetchTimetable_ReturnsErrorOnInvalidIcal checks that a 200 response
// with unparseable iCal content causes an error from the iCal parser,
// not from the HTTP layer.
func TestFetchTimetable_ReturnsErrorOnInvalidIcal(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("this is not valid ical data"))
	}))
	defer srv.Close()

	resetHTTPClient(t, srv.Client())
	zeroBackoff(t)
	resetSfGroup(t)

	_, err := FetchTimetable(5, srv.URL, 400, 42, 7)
	if err == nil {
		t.Fatal("expected error for invalid iCal body, got nil")
	}
	// The error must come from the iCal parser, not from a non-2xx status or
	// network failure. A successful HTTP round-trip followed by parse failure
	// will not contain "status code" in the error message.
	if strings.Contains(err.Error(), "status code") {
		t.Errorf("error looks like an HTTP error, not an iCal parse error: %v", err)
	}
}

// TestFetchTimetable_ReturnsErrorOnBadBaseURL checks that an unparseable
// base URL causes an error immediately without hitting the network.
func TestFetchTimetable_ReturnsErrorOnBadBaseURL(t *testing.T) {
	resetSfGroup(t)

	_, err := FetchTimetable(6, "://bad url", 500, 42, 7)
	if err == nil {
		t.Fatal("expected error for invalid base URL, got nil")
	}
}

// TestFetchTimetable_SingleflightDedup checks that concurrent calls with
// the same key result in only one upstream request.
func TestFetchTimetable_SingleflightDedup(t *testing.T) {
	var callCount atomic.Int32
	// arrived counts goroutines that have entered sfGroup.Do and are blocking.
	// The server waits until all goroutines are inside the singleflight group
	// before responding, ensuring no goroutine slips through after the first
	// call returns (which would register as a second upstream hit).
	var arrived atomic.Int32
	const goroutines = 5

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		// Spin until all goroutines have entered sfGroup.Do.
		// Because singleflight deduplicates, only one goroutine actually
		// reaches this handler; the others are parked inside Do. We just
		// need to wait long enough for them to park — bounded by a timeout.
		deadline := time.Now().Add(2 * time.Second)
		for arrived.Load() < goroutines && time.Now().Before(deadline) {
			time.Sleep(time.Millisecond)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(minimalIcal))
	}))
	defer srv.Close()

	resetHTTPClient(t, srv.Client())
	zeroBackoff(t)
	resetSfGroup(t)

	var wg sync.WaitGroup
	errs := make([]error, goroutines)

	for i := range goroutines {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			arrived.Add(1)
			_, errs[idx] = FetchTimetable(7, srv.URL, 600, 42, 7)
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d got error: %v", i, err)
		}
	}
	if int(callCount.Load()) != 1 {
		t.Errorf("singleflight: expected 1 upstream call, got %d", callCount.Load())
	}
}
