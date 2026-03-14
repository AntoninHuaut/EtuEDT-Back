package ade

import (
	"regexp"
	"sort"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Event struct {
	Title       string    `json:"title"`
	Teacher     string    `json:"teacher"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
}

var (
	reSuffix      = regexp.MustCompile(`(_s\d+)$`)
	rePrefix      = regexp.MustCompile(`(?m)^\w+.\d+ (:\s?)?`)
	reExportedMsg = regexp.MustCompile(`\n\(Export(é|ed).*\n?`)
)

func mergeSimilarEvents(events []Event) []Event {
	type mergeEvent struct {
		event       Event
		outputIndex int
	}

	outputs := make([]Event, 0)
	for _, event := range events {
		var existingEvents []mergeEvent
		for index, output := range outputs {
			if output.Title == event.Title &&
				output.Location == event.Location &&
				output.Teacher == event.Teacher &&
				(event.Start.Equal(output.End) || event.End.Equal(output.Start)) {
				existingEvents = append(existingEvents, mergeEvent{
					event:       output,
					outputIndex: index,
				})
			}
		}

		if len(existingEvents) == 0 {
			outputs = append(outputs, event)
		} else {
			for _, existing := range existingEvents {
				if event.Start.Before(existing.event.End) {
					event.End = existing.event.End
				} else {
					event.Start = existing.event.Start
				}
				outputs[existing.outputIndex] = event
			}
		}
	}

	return outputs
}

func CalendarToEvents(calendar *ics.Calendar) []Event {
	var events []Event

	formatTitle := func(title string) string {
		title = reSuffix.ReplaceAllString(title, "")
		return rePrefix.ReplaceAllString(title, "")
	}

	removeExportedDescription := func(description string) string {
		return reExportedMsg.ReplaceAllString(description, "")
	}

	formatDescription := func(description string) string {
		return strings.TrimSpace(removeExportedDescription(strings.ReplaceAll(description, "\\n", "\n")))
	}

	getTeacher := func(description string) string {
		if len(description) == 0 {
			return "?"
		}

		var teachers []string
		descSplit := strings.Split(description, "\n")
		firstTeacherIndex := 1
		for index, line := range descSplit {
			if strings.HasPrefix(line, "GRP") {
				firstTeacherIndex++
			}
			if index >= firstTeacherIndex {
				teachers = append(teachers, line)
			}
		}
		if len(teachers) == 0 {
			return "?"
		}
		return strings.Join(teachers, ",")
	}

	getLocation := func(location string) string {
		if len(location) == 0 {
			return "?"
		}
		return location
	}

	for _, event := range calendar.Events() {
		summary := event.GetProperty("SUMMARY")
		if summary == nil {
			continue
		}
		startAt, errStartAt := event.GetStartAt()
		endAt, errEndAt := event.GetEndAt()
		if errStartAt != nil || errEndAt != nil {
			continue
		}

		descriptionValue := ""
		if d := event.GetProperty("DESCRIPTION"); d != nil {
			descriptionValue = d.Value
		}
		locationValue := ""
		if l := event.GetProperty("LOCATION"); l != nil {
			locationValue = l.Value
		}

		formattedDescription := formatDescription(descriptionValue)
		events = append(events, Event{
			Title:       formatTitle(summary.Value),
			Teacher:     getTeacher(formattedDescription),
			Description: formattedDescription,
			Start:       startAt,
			End:         endAt,
			Location:    getLocation(locationValue),
		})
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].Start.Before(events[j].Start) {
			return true
		}
		if events[i].Start.After(events[j].Start) {
			return false
		}
		return events[i].Title < events[j].Title
	})

	return mergeSimilarEvents(events)
}
