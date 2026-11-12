package timetable

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
	"golang.org/x/sync/singleflight"
)

var ErrUnavailable = errors.New("timetable unavailable")

type Request struct {
	University *config.UniversityConfig
	Resource   int
}
type Source interface {
	Fetch(context.Context, Request) ([]ade.Event, error)
}
type Service struct {
	reparser Source
	group    singleflight.Group
}

func New(source Source) *Service { return &Service{reparser: source} }

func Scope(u *config.UniversityConfig, resource int) string {
	kind := "timetable"
	for _, room := range u.Rooms {
		if room.AdeResources == resource {
			kind = "room"
			break
		}
	}
	now := time.Now()
	first, _ := ade.GetAcademicYearDates(now, u.GetSplitMonth())
	identity := fmt.Sprintf("ade:%s:%d:%s:%s", u.AdeUrl, u.GetEffectiveProjectId(now), first.Format("2006-01-02"), kind)
	if u.UsesReparser() {
		identity = "reparser-v1:" + u.SourceOverride.BaseURL + ":" + first.Format("2006-01-02") + ":" + kind
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(identity)))
}

func Cached(u *config.UniversityConfig, resource int) (ade.TimetableCache, bool) {
	return ade.GetTimetableByAdeResources(u.ID, resource, Scope(u, resource))
}
func Store(u *config.UniversityConfig, resource int, events []ade.Event) ade.TimetableCache {
	return ade.SetTimetableByAdeResources(u.ID, resource, events, Scope(u, resource))
}

func (s *Service) Fetch(ctx context.Context, u *config.UniversityConfig, resource int) ([]ade.Event, error) {
	if resource <= 0 {
		return nil, ErrUnavailable
	}
	if u.UsesReparser() {
		if _, err := reparserPath(Request{u, resource}); err != nil {
			return nil, err
		}
	}
	ttl := u.CacheTTL()
	if cached, ok := Cached(u, resource); ok && cached.LastUpdate != nil && time.Since(*cached.LastUpdate) < ttl {
		return cached.Events, nil
	}
	key := fmt.Sprintf("%d:%d:%s", u.ID, resource, Scope(u, resource))
	result := s.group.DoChan(key, func() (any, error) {
		// A cancelled subscriber must not cancel a shared refresh for others.
		fetchCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
		defer cancel()
		if cached, ok := Cached(u, resource); ok && cached.LastUpdate != nil && time.Since(*cached.LastUpdate) < ttl {
			return cached.Events, nil
		}
		var source Source = ADE{}
		if u.UsesReparser() {
			source = s.reparser
		}
		events, err := source.Fetch(fetchCtx, Request{u, resource})
		if err != nil {
			if cached, ok := Cached(u, resource); ok {
				return cached.Events, nil
			}
			return nil, ErrUnavailable
		}
		if events == nil {
			events = []ade.Event{}
		}
		return Store(u, resource, events).Events, nil
	})
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case out := <-result:
		if out.Err != nil {
			return nil, ErrUnavailable
		}
		return out.Val.([]ade.Event), nil
	}
}
