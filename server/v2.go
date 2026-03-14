package server

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/domain"
	ics "github.com/arran4/golang-ical"
	"github.com/go-chi/chi/v5"
)

func v2Router(r chi.Router) {
	r.Get("/", legacyV2ICSMigrationNotice)
	r.Get("/*", legacyV2ICSMigrationNotice)
}

func legacyV2ICSMigrationNotice(w http.ResponseWriter, r *http.Request) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)

	now := time.Now()
	startTime, endTime := domain.GetAcademicYearDates(now)
	// iCal all-day DTEND is exclusive, so add one day.
	endTime = endTime.AddDate(0, 0, 1)

	event := cal.AddEvent(fmt.Sprintf("migration-notice@etuedt-%d", now.Unix()))
	event.SetSummary("⚠️ API v2 is deprecated — please migrate to v3")
	event.SetDescription("The v2 API has been replaced by v3. See documentation for migration details.")
	event.SetDtStampTime(now)
	event.SetAllDayStartAt(startTime)
	event.SetAllDayEndAt(endTime)

	var buf bytes.Buffer
	if err := cal.SerializeTo(&buf); err != nil {
		http.Error(w, "failed to generate calendar", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	_, _ = buf.WriteTo(w)
}
