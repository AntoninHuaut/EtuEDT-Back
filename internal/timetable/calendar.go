package timetable

import (
	"fmt"
	"strings"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
)

// CalendarURL preserves the public subscription-link contract for either source.
func CalendarURL(u *config.UniversityConfig, resource int, now, firstDate, lastDate time.Time) (string, error) {
	if u.UsesReparser() {
		return fmt.Sprintf("%s/v3/ics/%d/%d", strings.TrimRight(u.SourceOverride.BaseURL, "/"), u.ID, resource), nil
	}
	return ade.BuildURL(u.AdeUrl, resource, u.GetEffectiveProjectId(now), firstDate, lastDate)
}
