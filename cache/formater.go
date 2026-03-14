package cache

import (
	"regexp"
	"sort"
	"strings"

	"github.com/AntoninHuaut/EtuEDT-Back/domain"
	ics "github.com/arran4/golang-ical"
)

var (
	reSuffix      = regexp.MustCompile(`(_s\d+)$`)
	rePrefix      = regexp.MustCompile(`(?m)^\w+.\d+ (:\s?)?`)
	reExportedMsg = regexp.MustCompile(`\n\(Export(é|ed).*\n?`)
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

func CalendarToJson(calendar *ics.Calendar) []domain.JsonEvent {
	var jsonEvents []domain.JsonEvent

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
		jsonEvents = append(jsonEvents, domain.JsonEvent{
			Title:       formatTitle(summary.Value),
			Teacher:     getTeacher(formattedDescription),
			Description: formattedDescription,
			Start:       startAt,
			End:         endAt,
			Location:    getLocation(locationValue),
		})
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
