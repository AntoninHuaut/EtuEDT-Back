package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

type TimetableConfig struct {
	AdeResources int    `json:"adeResources" validate:"gt=0"`
	Year         int    `json:"year"         validate:"gt=0"`
	Label        string `json:"label"        validate:"required"`
}

type RoomConfig struct {
	AdeResources int    `json:"adeResources" validate:"gt=0"`
	Label        string `json:"label"        validate:"required"`
}

type GroupConfig struct {
	ID         int               `json:"id"         validate:"gt=0"`
	Name       string            `json:"name"       validate:"required"`
	Timetables []TimetableConfig `json:"timetables" validate:"required,min=1,dive"`
}

type AdeProjectIdCycleConfig struct {
	StartYear  int   `json:"startYear"  validate:"gt=0"`
	SplitMonth int   `json:"splitMonth" validate:"min=1,max=12"`
	Cycle      []int `json:"cycle"      validate:"required,min=1,dive,gt=0"`
}

func (c *AdeProjectIdCycleConfig) GetProjectId(now time.Time) int {
	year := now.Year()
	if now.Month() < time.Month(c.SplitMonth) {
		year--
	}
	offset := (year - c.StartYear) % len(c.Cycle)
	if offset < 0 {
		offset += len(c.Cycle)
	}
	return c.Cycle[offset]
}

type UniversityConfig struct {
	ID                int                       `json:"id"                validate:"gt=0"`
	Name              string                    `json:"name"              validate:"required"`
	AdeUrl            string                    `json:"adeUrl"            validate:"required,http_url"`
	AdeProjectId      int                       `json:"adeProjectId"`
	AdeProjectIdCycle *AdeProjectIdCycleConfig  `json:"adeProjectIdCycle,omitempty"`
	Rooms             []RoomConfig              `json:"rooms"             validate:"dive"`
	Groups            []GroupConfig             `json:"groups"            validate:"dive"`
}

func (u *UniversityConfig) GetEffectiveProjectId(now time.Time) int {
	if u.AdeProjectIdCycle != nil {
		return u.AdeProjectIdCycle.GetProjectId(now)
	}
	return u.AdeProjectId
}

func (u *UniversityConfig) GetSplitMonth() int {
	if u.AdeProjectIdCycle != nil {
		return u.AdeProjectIdCycle.SplitMonth
	}
	return 7
}

type Config struct {
	Universities []UniversityConfig `json:"univs" validate:"required,min=1,dive"`
}

var AppConfig Config

var validate = validator.New()

func LoadConfig() error {
	file, err := os.Open("./config.json")
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return validateConfig(&AppConfig)
}

func validateConfig(config *Config) error {
	if err := validate.Struct(config); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
			return fmt.Errorf("config validation error: %s", ve[0].Error())
		}
		return err
	}

	univIDs := make(map[int]bool)
	for _, univ := range config.Universities {
		if univ.AdeProjectIdCycle == nil && univ.AdeProjectId <= 0 {
			return fmt.Errorf("adeProjectId must be > 0 when adeProjectIdCycle is not set")
		}

		if univIDs[univ.ID] {
			return fmt.Errorf("duplicate university id: %d", univ.ID)
		}
		univIDs[univ.ID] = true

		adeResourcesSet := make(map[int]bool)
		for _, room := range univ.Rooms {
			if adeResourcesSet[room.AdeResources] {
				return fmt.Errorf("duplicate adeResources: %d", room.AdeResources)
			}
			adeResourcesSet[room.AdeResources] = true
		}

		groupIDs := make(map[int]bool)
		for _, group := range univ.Groups {
			if groupIDs[group.ID] {
				return fmt.Errorf("duplicate group id: %d in university %d", group.ID, univ.ID)
			}
			groupIDs[group.ID] = true
			for _, tt := range group.Timetables {
				if adeResourcesSet[tt.AdeResources] {
					return fmt.Errorf("duplicate adeResources: %d", tt.AdeResources)
				}
				adeResourcesSet[tt.AdeResources] = true
			}
		}
	}

	return nil
}
