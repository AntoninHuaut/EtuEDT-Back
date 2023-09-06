package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type TimetableConfig struct {
	NumYearTT    int    `json:"numYearTT"`
	DescTT       string `json:"descTT"`
	AdeResources int    `json:"adeResources"`
	AdeProjectId int    `json:"adeProjectId"`
}

type UniversityConfig struct {
	NumUniv    int               `json:"numUniv"`
	NameUniv   string            `json:"nameUniv"`
	AdeUniv    string            `json:"adeUniv"`
	Timetables []TimetableConfig `json:"timetable"`
}

type Config struct {
	RefreshMinutes int                `json:"refreshMinutes"`
	Universities   []UniversityConfig `json:"univ"`
}

var AppConfig Config

func LoadConfig() error {
	file, err := os.Open("./config/config.json")
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		return err
	}

	if AppConfig.RefreshMinutes < 1 {
		return errors.New("refreshMinutes must be greater than 0")
	}

	return nil
}
