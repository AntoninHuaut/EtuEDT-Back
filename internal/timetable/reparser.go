package timetable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
)

const maxResponseBytes = 8 << 20

type Reparser struct {
	client *http.Client
	slots  chan struct{}
}

func NewReparser(client *http.Client) *Reparser {
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	copy := *client
	copy.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &Reparser{client: &copy, slots: make(chan struct{}, 3)}
}

type reparserEvent struct {
	Title       string    `json:"title"`
	Teacher     string    `json:"teacher"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

func (r *Reparser) Fetch(ctx context.Context, req Request) ([]ade.Event, error) {
	path, err := reparserPath(req)
	if err != nil {
		return nil, err
	}
	u := req.University
	select {
	case r.slots <- struct{}{}:
	case <-ctx.Done():
		return nil, ErrUnavailable
	}
	defer func() { <-r.slots }()
	for attempt := 0; attempt < 2; attempt++ {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(u.SourceOverride.BaseURL, "/")+path, nil)
		if err != nil {
			return nil, ErrUnavailable
		}
		response, err := r.client.Do(request)
		if err == nil {
			body, readErr := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				if readErr != nil || len(body) > maxResponseBytes {
					return nil, ErrUnavailable
				}
				return decodeEvents(body)
			}
			if response.StatusCode < 500 && response.StatusCode != http.StatusTooManyRequests {
				return nil, ErrUnavailable
			}
		}
		if attempt == 0 {
			timer := time.NewTimer(time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ErrUnavailable
			case <-timer.C:
			}
		}
	}
	return nil, ErrUnavailable
}

func decodeEvents(body []byte) ([]ade.Event, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 || trimmed[0] != '[' {
		return nil, ErrUnavailable
	}
	var rows []reparserEvent
	if err := json.Unmarshal(trimmed, &rows); err != nil {
		return nil, ErrUnavailable
	}
	events := make([]ade.Event, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Title) == "" || row.Start.IsZero() || row.End.IsZero() || !row.End.After(row.Start) {
			return nil, ErrUnavailable
		}
		teacher, room := row.Teacher, row.Location
		if teacher == "" {
			teacher = "?"
		}
		if room == "" {
			room = "?"
		}
		events = append(events, ade.Event{Title: ade.FormatTitle(row.Title), Teacher: teacher, Location: room,
			Description: row.Description, Start: row.Start.UTC(), End: row.End.UTC()})
	}
	return ade.FormatEvents(events), nil
}

// Resolve only resources present in the local catalogue.
func reparserPath(req Request) (string, error) {
	u := req.University
	if req.Resource <= 0 {
		return "", ErrUnavailable
	}
	path := ""
	for _, room := range u.Rooms {
		if room.AdeResources == req.Resource {
			path = fmt.Sprintf("/v3/univs/%d/rooms/%d/events", u.ID, req.Resource)
		}
	}
	for _, group := range u.Groups {
		for _, tt := range group.Timetables {
			if tt.AdeResources == req.Resource {
				path = fmt.Sprintf("/v3/univs/%d/groups/%d/%d/events", u.ID, group.ID, req.Resource)
			}
		}
	}
	if path == "" {
		return "", ErrUnavailable
	}
	return path, nil
}
