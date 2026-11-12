package timetable

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
)

type sourceFunc func(context.Context, Request) ([]ade.Event, error)

func (f sourceFunc) Fetch(ctx context.Context, req Request) ([]ade.Event, error) { return f(ctx, req) }

func TestServiceCacheAndStaleFallback(t *testing.T) {
	u := testUniversity("https://cache.example.test")
	var calls int
	s := New(sourceFunc(func(context.Context, Request) ([]ade.Event, error) {
		calls++
		if calls > 1 {
			return nil, errors.New("upstream failure")
		}
		return []ade.Event{{Title: "Cached"}}, nil
	}))
	for i := 0; i < 2; i++ {
		events, err := s.Fetch(context.Background(), u, 20)
		if err != nil || len(events) != 1 {
			t.Fatalf("events %v, error %v", events, err)
		}
	}
	if calls != 1 {
		t.Fatalf("fresh cache fetched %d times", calls)
	}
	cached, _ := Cached(u, 20)
	*cached.LastUpdate = time.Now().Add(-time.Hour)
	events, err := s.Fetch(context.Background(), u, 20)
	if err != nil || len(events) != 1 || calls != 2 {
		t.Fatalf("stale fallback: %v, %v, calls %d", events, err, calls)
	}
	if time.Since(*cached.LastUpdate) < time.Hour {
		t.Fatal("failure advanced cache timestamp")
	}
	u.SourceOverride.BaseURL = "https://other.example.test"
	if _, err := s.Fetch(context.Background(), u, 20); err != ErrUnavailable {
		t.Fatalf("cache crossed provider URLs: %v", err)
	}
	u.SourceOverride = nil
	if _, ok := Cached(u, 20); ok {
		t.Fatal("cache crossed providers")
	}
	Store(u, 20, events)
	u.AdeProjectId++
	if _, ok := Cached(u, 20); ok {
		t.Fatal("cache crossed ADE projects")
	}
}

func TestServiceSharedRefreshSurvivesCancelledSubscriber(t *testing.T) {
	u := testUniversity("https://shared.example.test")
	started, release := make(chan struct{}), make(chan struct{})
	var calls atomic.Int32
	s := New(sourceFunc(func(ctx context.Context, req Request) ([]ade.Event, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		select {
		case <-release:
			return []ade.Event{}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	first := make(chan error, 1)
	go func() { _, err := s.Fetch(ctx, u, 20); first <- err }()
	<-started
	cancel()
	if err := <-first; !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
	second := make(chan error, 1)
	go func() { _, err := s.Fetch(context.Background(), u, 20); second <- err }()
	close(release)
	if err := <-second; err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatalf("duplicate requests: %d", calls.Load())
	}
}
