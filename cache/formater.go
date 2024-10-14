package cache

import (
	"EtuEDT-Go/domain"
	ics "github.com/arran4/golang-ical"
	"regexp"
	"sort"
	"strings"
)

func jsonMergeSimilarEvents(events []domain.JsonEvent) []domain.JsonEvent {
	type MergeJsonEvent struct {
		Event       domain.JsonEvent
		OutputIndex int
	}

	outputs := make([]domain.JsonEvent, 0)
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

func calendarToJson(calendar *ics.Calendar) []domain.JsonEvent {
	var jsonEvents []domain.JsonEvent

	formatTitle := func(title string) string {
		title = regexp.MustCompile(`(_s\d+)$`).ReplaceAllString(title, "")
		return regexp.MustCompile(`(?m)^\w+.\d+ (:\s?)?`).ReplaceAllString(title, "")
	}

	removeExportedDescription := func(description string) string {
		return regexp.MustCompile(`\n\(Export(é|ed).*\n?`).ReplaceAllString(description, "")
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

	for _, event := range calendar.Events() {
		summary := event.GetProperty("SUMMARY")
		description := event.GetProperty("DESCRIPTION")
		location := event.GetProperty("LOCATION")
		if summary != nil && description != nil && location != nil {
			startAt, errStartAt := event.GetStartAt()
			endAt, errEndAt := event.GetEndAt()
			if errStartAt == nil && errEndAt == nil {
				formattedDescription := formatDescription(description.Value)
				teacher := getTeacher(formattedDescription)

				jsonEvents = append(jsonEvents, domain.JsonEvent{
					Title:       formatTitle(summary.Value),
					Teacher:     teacher,
					Description: formattedDescription,
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
