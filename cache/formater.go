package cache

import (
	ics "github.com/arran4/golang-ical"
	"regexp"
	"sort"
	"strings"
)

func jsonMergeSimilarEvents(events []JsonEvent) []JsonEvent {
	type MergeJsonEvent struct {
		Event       JsonEvent
		OutputIndex int
	}

	outputs := make([]JsonEvent, 0)
	for _, event := range events {
		var existingEvents []MergeJsonEvent
		for index, output := range outputs {
			if output.Title == event.Title &&
				output.Location == event.Location &&
				output.Teacher == event.Teacher &&
				(event.Start.Equal(output.End) || event.End.Equal(output.Start)) {
				existingEvents = append(existingEvents, MergeJsonEvent{
					Event:       output,
					OutputIndex: index,
				})
			}
		}

		if len(existingEvents) == 0 {
			outputs = append(outputs, event)
		} else {
			for _, existing := range existingEvents {
				if event.Start.Before(existing.Event.End) {
					event.End = existing.Event.End
				} else {
					event.Start = existing.Event.Start
				}
				outputs[existing.OutputIndex] = event
			}
		}
	}

	return outputs
}

func calendarToJson(calendar *ics.Calendar) []JsonEvent {
	var jsonEvents []JsonEvent

	formatTitle := func(title string) string {
		title = regexp.MustCompile(`(_s\d+)$`).ReplaceAllString(title, "")
		return regexp.MustCompile(`(?m)^\w+.\d+ (: )?`).ReplaceAllString(title, "")
	}

	removeExportedDescription := func(description string) string {
		return regexp.MustCompile(`\(Exported.*\n?`).ReplaceAllString(description, "")
	}

	formatDescription := func(description string) string {
		return strings.TrimSpace(strings.ReplaceAll(removeExportedDescription(description), "\\n", "\n"))
	}

	getTeacher := func(description string) string {
		if len(description) == 0 {
			return "?"
		}

		descSplit := strings.Split(removeExportedDescription(description), "\\n")
		teacherIndex := 1
		counter := 0
		for _, line := range descSplit {
			if len(line) == 0 {
				continue
			}
			if strings.HasPrefix(line, "GRP") {
				teacherIndex++
			}
			if counter == teacherIndex {
				return line
			}
			counter++
		}
		return "?"
	}

	for _, event := range calendar.Events() {
		summary := event.GetProperty("SUMMARY")
		description := event.GetProperty("DESCRIPTION")
		location := event.GetProperty("LOCATION")
		if summary != nil && description != nil && location != nil {
			startAt, errStartAt := event.GetStartAt()
			endAt, errEndAt := event.GetEndAt()
			if errStartAt == nil && errEndAt == nil {
				teacher := getTeacher(description.Value)
				jsonEvents = append(jsonEvents, JsonEvent{
					Title:       formatTitle(summary.Value),
					Enseignant:  teacher, // TODO to be removed
					Teacher:     teacher,
					Description: formatDescription(description.Value),
					Start:       startAt,
					End:         endAt,
					Location:    location.Value,
				})
			}
		}
	}

	sort.Slice(jsonEvents, func(i, j int) bool {
		if jsonEvents[i].Start.Before(jsonEvents[j].Start) {
			return true
		}
		if jsonEvents[i].Start.After(jsonEvents[j].Start) {
			return false
		}
		return jsonEvents[i].Title < jsonEvents[j].Title
	})

	return jsonMergeSimilarEvents(jsonEvents)
}
