package ade

import (
	"strings"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
)

// buildCalendar parses an iCal string and returns a *ics.Calendar.
// Panics on parse error (tests only).
func buildCalendar(ical string) *ics.Calendar {
	cal, err := ics.ParseCalendar(strings.NewReader(ical))
	if err != nil {
		panic("buildCalendar: " + err.Error())
	}
	return cal
}

// calendarWithEvent builds a minimal iCal string with one VEVENT.
func calendarWithEvent(uid, summary, dtstart, dtend, description, location string) string {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\r\nVERSION:2.0\r\n")
	sb.WriteString("BEGIN:VEVENT\r\n")
	sb.WriteString("UID:" + uid + "\r\n")
	if summary != "" {
		sb.WriteString("SUMMARY:" + summary + "\r\n")
	}
	if dtstart != "" {
		sb.WriteString("DTSTART:" + dtstart + "\r\n")
	}
	if dtend != "" {
		sb.WriteString("DTEND:" + dtend + "\r\n")
	}
	if description != "" {
		sb.WriteString("DESCRIPTION:" + description + "\r\n")
	}
	if location != "" {
		sb.WriteString("LOCATION:" + location + "\r\n")
	}
	sb.WriteString("END:VEVENT\r\n")
	sb.WriteString("END:VCALENDAR\r\n")
	return sb.String()
}

// TestCalendarToEvents_BasicEvent verifies a simple event is parsed correctly.
func TestCalendarToEvents_BasicEvent(t *testing.T) {
	ical := calendarWithEvent(
		"uid-1",
		"INFO3-Algo : Cours",
		"20250901T090000Z",
		"20250901T110000Z",
		"Group A\\nDr. Smith",
		"Amphi A",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e := events[0]
	if e.Location != "Amphi A" {
		t.Errorf("Location: got %q, want %q", e.Location, "Amphi A")
	}
	if e.Start.IsZero() || e.End.IsZero() {
		t.Errorf("Start/End should be set: start=%v end=%v", e.Start, e.End)
	}
}

// TestCalendarToEvents_SkipsEventWithoutSummary verifies events without SUMMARY are dropped.
func TestCalendarToEvents_SkipsEventWithoutSummary(t *testing.T) {
	ical := calendarWithEvent(
		"uid-1",
		"", // no summary
		"20250901T090000Z",
		"20250901T110000Z",
		"",
		"",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 0 {
		t.Errorf("expected 0 events (no summary), got %d", len(events))
	}
}

// TestCalendarToEvents_SkipsEventWithBadDates verifies events with invalid dates are dropped.
func TestCalendarToEvents_SkipsEventWithBadDates(t *testing.T) {
	ical := calendarWithEvent(
		"uid-1",
		"Math",
		"NOT_A_DATE",
		"ALSO_NOT_A_DATE",
		"",
		"",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 0 {
		t.Errorf("expected 0 events (bad dates), got %d", len(events))
	}
}

// TestCalendarToEvents_TitleStripsAdeSuffix verifies _s<digits> suffix is stripped.
func TestCalendarToEvents_TitleStripsAdeSuffix(t *testing.T) {
	ical := calendarWithEvent(
		"uid-1",
		"INFO3.123 : Algo_s1",
		"20250901T090000Z",
		"20250901T110000Z",
		"", "",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if strings.HasSuffix(events[0].Title, "_s1") {
		t.Errorf("Title should not end with _s1, got %q", events[0].Title)
	}
}

// TestCalendarToEvents_TitleStripsAdePrefix verifies WORD.DIGITS prefix is stripped.
func TestCalendarToEvents_TitleStripsAdePrefix(t *testing.T) {
	ical := calendarWithEvent(
		"uid-1",
		"INFO3.456 : Algorithmique",
		"20250901T090000Z",
		"20250901T110000Z",
		"", "",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Title != "Algorithmique" {
		t.Errorf("Title: got %q, want %q", events[0].Title, "Algorithmique")
	}
}

// TestCalendarToEvents_LocationFallback verifies "?" is used when location is empty.
func TestCalendarToEvents_LocationFallback(t *testing.T) {
	ical := calendarWithEvent(
		"uid-1",
		"Math",
		"20250901T090000Z",
		"20250901T110000Z",
		"", "",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Location != "?" {
		t.Errorf("Location: got %q, want %q", events[0].Location, "?")
	}
}

// TestCalendarToEvents_TeacherFallback verifies "?" is returned when description is empty.
func TestCalendarToEvents_TeacherFallback(t *testing.T) {
	ical := calendarWithEvent(
		"uid-1",
		"Math",
		"20250901T090000Z",
		"20250901T110000Z",
		"", "",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Teacher != "?" {
		t.Errorf("Teacher: got %q, want %q", events[0].Teacher, "?")
	}
}

// TestCalendarToEvents_TeacherExtracted verifies teacher is extracted from description (second line onward).
func TestCalendarToEvents_TeacherExtracted(t *testing.T) {
	// Description: first line is ignored, second is teacher
	// iCal uses \n as literal backslash-n
	ical := calendarWithEvent(
		"uid-1",
		"Math",
		"20250901T090000Z",
		"20250901T110000Z",
		"Group A\\nDr. Smith",
		"",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Teacher != "Dr. Smith" {
		t.Errorf("Teacher: got %q, want %q", events[0].Teacher, "Dr. Smith")
	}
}

// TestCalendarToEvents_TeacherSkipsGRPLines verifies GRP lines are skipped when extracting teacher.
func TestCalendarToEvents_TeacherSkipsGRPLines(t *testing.T) {
	// First line: group name (ignored), then GRP line (skipped), then teacher
	ical := calendarWithEvent(
		"uid-1",
		"Math",
		"20250901T090000Z",
		"20250901T110000Z",
		"Group A\\nGRP B\\nDr. Jones",
		"",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Teacher != "Dr. Jones" {
		t.Errorf("Teacher: got %q, want %q", events[0].Teacher, "Dr. Jones")
	}
}

// TestCalendarToEvents_DescriptionStripsExportedFooter verifies export footer is removed.
func TestCalendarToEvents_DescriptionStripsExportedFooter(t *testing.T) {
	// The description contains an "(Exporté..." footer that should be stripped
	ical := calendarWithEvent(
		"uid-1",
		"Math",
		"20250901T090000Z",
		"20250901T110000Z",
		"Group A\\nDr. Smith\\n(Exporté le 01/01/2025)",
		"",
	)

	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if strings.Contains(events[0].Description, "Exporté") {
		t.Errorf("Description should not contain export footer, got %q", events[0].Description)
	}
}

// TestCalendarToEvents_SortedByStartThenTitle verifies events are sorted by start time, then title.
func TestCalendarToEvents_SortedByStartThenTitle(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\r\nVERSION:2.0\r\n")

	// Event B at 10:00
	sb.WriteString("BEGIN:VEVENT\r\nUID:b\r\nSUMMARY:Biology\r\nDTSTART:20250901T100000Z\r\nDTEND:20250901T110000Z\r\nEND:VEVENT\r\n")
	// Event A at 09:00
	sb.WriteString("BEGIN:VEVENT\r\nUID:a\r\nSUMMARY:Algorithmique\r\nDTSTART:20250901T090000Z\r\nDTEND:20250901T100000Z\r\nEND:VEVENT\r\n")
	// Event C at 09:00 (same time as A, sorted by title after A)
	sb.WriteString("BEGIN:VEVENT\r\nUID:c\r\nSUMMARY:Chemistry\r\nDTSTART:20250901T090000Z\r\nDTEND:20250901T100000Z\r\nEND:VEVENT\r\n")
	sb.WriteString("END:VCALENDAR\r\n")

	events := CalendarToEvents(buildCalendar(sb.String()))
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	// Algorithmique and Chemistry are at 09:00, Biology at 10:00
	if events[0].Title != "Algorithmique" {
		t.Errorf("events[0].Title: got %q, want %q", events[0].Title, "Algorithmique")
	}
	if events[1].Title != "Chemistry" {
		t.Errorf("events[1].Title: got %q, want %q", events[1].Title, "Chemistry")
	}
	if events[2].Title != "Biology" {
		t.Errorf("events[2].Title: got %q, want %q", events[2].Title, "Biology")
	}
}

// TestCalendarToEvents_EmptyCalendar verifies an empty calendar returns no events.
func TestCalendarToEvents_EmptyCalendar(t *testing.T) {
	ical := "BEGIN:VCALENDAR\r\nVERSION:2.0\r\nEND:VCALENDAR\r\n"
	events := CalendarToEvents(buildCalendar(ical))
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

// --- mergeSimilarEvents tests ---

func makeEvent(title, location, teacher string, start, end time.Time) Event {
	return Event{Title: title, Location: location, Teacher: teacher, Start: start, End: end}
}

// TestMergeSimilarEvents_NoMerge verifies distinct events are not merged.
func TestMergeSimilarEvents_NoMerge(t *testing.T) {
	t0 := time.Date(2025, 9, 1, 8, 0, 0, 0, time.UTC)
	t1 := time.Date(2025, 9, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 9, 1, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2025, 9, 1, 14, 0, 0, 0, time.UTC)

	events := []Event{
		makeEvent("Math", "A", "Smith", t0, t1),
		makeEvent("Physics", "B", "Jones", t2, t3),
	}

	result := mergeSimilarEvents(events)
	if len(result) != 2 {
		t.Errorf("expected 2 events after no-merge, got %d", len(result))
	}
}

// TestMergeSimilarEvents_MergesConsecutive verifies two back-to-back identical events are merged.
func TestMergeSimilarEvents_MergesConsecutive(t *testing.T) {
	t0 := time.Date(2025, 9, 1, 8, 0, 0, 0, time.UTC)
	t1 := time.Date(2025, 9, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 9, 1, 12, 0, 0, 0, time.UTC)

	events := []Event{
		makeEvent("Math", "A", "Smith", t0, t1),
		makeEvent("Math", "A", "Smith", t1, t2), // starts exactly when first ends
	}

	result := mergeSimilarEvents(events)
	if len(result) != 1 {
		t.Fatalf("expected 1 merged event, got %d", len(result))
	}
	if !result[0].Start.Equal(t0) {
		t.Errorf("merged Start: got %v, want %v", result[0].Start, t0)
	}
	if !result[0].End.Equal(t2) {
		t.Errorf("merged End: got %v, want %v", result[0].End, t2)
	}
}

// TestMergeSimilarEvents_DifferentTitleNotMerged verifies events with different titles are not merged.
func TestMergeSimilarEvents_DifferentTitleNotMerged(t *testing.T) {
	t0 := time.Date(2025, 9, 1, 8, 0, 0, 0, time.UTC)
	t1 := time.Date(2025, 9, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 9, 1, 12, 0, 0, 0, time.UTC)

	events := []Event{
		makeEvent("Math", "A", "Smith", t0, t1),
		makeEvent("Physics", "A", "Smith", t1, t2),
	}

	result := mergeSimilarEvents(events)
	if len(result) != 2 {
		t.Errorf("expected 2 events (different titles), got %d", len(result))
	}
}

// TestMergeSimilarEvents_DifferentLocationNotMerged verifies events with different locations are not merged.
func TestMergeSimilarEvents_DifferentLocationNotMerged(t *testing.T) {
	t0 := time.Date(2025, 9, 1, 8, 0, 0, 0, time.UTC)
	t1 := time.Date(2025, 9, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 9, 1, 12, 0, 0, 0, time.UTC)

	events := []Event{
		makeEvent("Math", "Room A", "Smith", t0, t1),
		makeEvent("Math", "Room B", "Smith", t1, t2),
	}

	result := mergeSimilarEvents(events)
	if len(result) != 2 {
		t.Errorf("expected 2 events (different locations), got %d", len(result))
	}
}

// TestMergeSimilarEvents_NonAdjacentNotMerged verifies non-adjacent events are not merged.
func TestMergeSimilarEvents_NonAdjacentNotMerged(t *testing.T) {
	t0 := time.Date(2025, 9, 1, 8, 0, 0, 0, time.UTC)
	t1 := time.Date(2025, 9, 1, 10, 0, 0, 0, time.UTC)
	// Gap between t1 and t2
	t2 := time.Date(2025, 9, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2025, 9, 1, 13, 0, 0, 0, time.UTC)

	events := []Event{
		makeEvent("Math", "A", "Smith", t0, t1),
		makeEvent("Math", "A", "Smith", t2, t3),
	}

	result := mergeSimilarEvents(events)
	if len(result) != 2 {
		t.Errorf("expected 2 events (gap between them), got %d", len(result))
	}
}
