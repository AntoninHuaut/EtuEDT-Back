package config

import (
	"fmt"
	"net/url"
	"time"
)

// SourceOverrideConfig selects an alternative provider for a university.
// It is deployment configuration, never part of the public API.
type SourceOverrideConfig struct {
	Provider     string `json:"provider"`
	BaseURL      string `json:"baseUrl"`
	CacheMinutes int    `json:"cacheMinutes,omitempty"`
}

func loadSourceSettings(cfg *Config) error {
	for i := range cfg.Universities {
		u := &cfg.Universities[i]
		override := u.SourceOverride
		if override == nil {
			continue
		}
		if override.Provider != "reparser" {
			return fmt.Errorf("university %d: unsupported sourceOverride.provider", u.ID)
		}
		parsed, err := url.Parse(override.BaseURL)
		if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" || (parsed.Path != "" && parsed.Path != "/") {
			return fmt.Errorf("university %d: sourceOverride.baseUrl must be an HTTPS origin", u.ID)
		}
		if override.CacheMinutes == 0 {
			override.CacheMinutes = 30
		}
		if override.CacheMinutes < 5 || override.CacheMinutes > 720 {
			return fmt.Errorf("university %d: sourceOverride.cacheMinutes must be between 5 and 720", u.ID)
		}
	}
	return nil
}

func (u *UniversityConfig) UsesReparser() bool {
	return u.SourceOverride != nil && u.SourceOverride.Provider == "reparser"
}

func (u *UniversityConfig) CacheTTL() time.Duration {
	if !u.UsesReparser() {
		return 12 * time.Hour
	}
	if u.SourceOverride.CacheMinutes == 0 {
		return 30 * time.Minute
	}
	return time.Duration(u.SourceOverride.CacheMinutes) * time.Minute
}
