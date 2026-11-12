package timetable

import (
	"context"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"time"
)

type ADE struct{}

func (ADE) Fetch(ctx context.Context, req Request) ([]ade.Event, error) {
	u := req.University
	cal, err := ade.FetchTimetableContext(ctx, u.ID, u.AdeUrl, req.Resource, u.GetEffectiveProjectId(time.Now()), u.GetSplitMonth())
	if err != nil {
		return nil, ErrUnavailable
	}
	return ade.CalendarToEvents(cal), nil
}
