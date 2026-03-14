package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type TimetableConfig struct {
	AdeResources int    `json:"adeResources"`
	Year         int    `json:"year"`
	Label        string `json:"label"`
}

type RoomConfig struct {
	AdeResources int    `json:"adeResources"`
	Label        string `json:"label"`
}

type GroupConfig struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Timetables []TimetableConfig `json:"timetables"`
}

type UniversityConfig struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	AdeUrl       string        `json:"adeUrl"`
	AdeProjectId int           `json:"adeProjectId"`
	Rooms        []RoomConfig  `json:"rooms"`
	Groups       []GroupConfig `json:"groups"`
}

type Config struct {
	Universities []UniversityConfig `json:"univs"`
}

var AppConfig Config

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
	if len(config.Universities) == 0 {
		return errors.New("at least one university must be configured")
	}

	univIDs := make(map[int]bool)
	adeResourcesSet := make(map[int]bool)

	for _, univ := range config.Universities {
		if univ.ID <= 0 {
			return fmt.Errorf("invalid university id: %d", univ.ID)
		}
		if univ.Name == "" {
			return fmt.Errorf("university %d name cannot be empty", univ.ID)
		}
		if univ.AdeUrl == "" {
			return fmt.Errorf("university %d adeUrl cannot be empty", univ.ID)
		}
		if univ.AdeProjectId <= 0 {
			return fmt.Errorf("university %d adeProjectId must be greater than 0", univ.ID)
		}

		if univIDs[univ.ID] {
			return fmt.Errorf("duplicate university id: %d", univ.ID)
		}
		univIDs[univ.ID] = true

		groupIDs := make(map[int]bool)

		for _, room := range univ.Rooms {
			if room.AdeResources <= 0 {
				return fmt.Errorf("invalid room adeResources in university %d", univ.ID)
			}
			if room.Label == "" {
				return fmt.Errorf("empty room label in university %d for adeResources %d", univ.ID, room.AdeResources)
			}
			if adeResourcesSet[room.AdeResources] {
				return fmt.Errorf("duplicate adeResources: %d", room.AdeResources)
			}
			adeResourcesSet[room.AdeResources] = true
		}

		for _, group := range univ.Groups {
			if group.ID <= 0 {
				return fmt.Errorf("invalid group id in university %d", univ.ID)
			}
			if group.Name == "" {
				return fmt.Errorf("empty group name in university %d", univ.ID)
			}
			if groupIDs[group.ID] {
				return fmt.Errorf("duplicate group id: %d in university %d", group.ID, univ.ID)
			}
			groupIDs[group.ID] = true
			if len(group.Timetables) == 0 {
				return fmt.Errorf("group %d in university %d has no timetables", group.ID, univ.ID)
			}

			for _, tt := range group.Timetables {
				if tt.AdeResources <= 0 {
					return fmt.Errorf("invalid timetable adeResources in group %d of university %d", group.ID, univ.ID)
				}
				if tt.Year <= 0 {
					return fmt.Errorf("invalid timetable year in group %d of university %d", group.ID, univ.ID)
				}
				if tt.Label == "" {
					return fmt.Errorf("empty timetable label in group %d of university %d", group.ID, univ.ID)
				}
				if adeResourcesSet[tt.AdeResources] {
					return fmt.Errorf("duplicate adeResources: %d", tt.AdeResources)
				}
				adeResourcesSet[tt.AdeResources] = true
			}
		}
	}

	return nil
}
