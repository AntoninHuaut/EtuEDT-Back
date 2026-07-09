package config

import (
	"strings"
	"testing"
	"time"
)

func validConfig() Config {
	return Config{
		Universities: []UniversityConfig{
			{
				ID:           1,
				Name:         "Test University",
				AdeUrl:       "https://ade.example.com",
				AdeProjectId: 42,
				Rooms: []RoomConfig{
					{AdeResources: 10, Label: "Room A"},
				},
				Groups: []GroupConfig{
					{
						ID:   1,
						Name: "Group 1",
						Timetables: []TimetableConfig{
							{AdeResources: 20, Year: 2025, Label: "Timetable A"},
						},
					},
				},
			},
		},
	}
}

func TestValidateConfig_Valid(t *testing.T) {
	cfg := validConfig()
	if err := validateConfig(&cfg); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateConfig_NoUniversities(t *testing.T) {
	cfg := Config{Universities: []UniversityConfig{}}
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for empty universities, got nil")
	}
}

func TestValidateConfig_UniversityIDZero(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].ID = 0
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for university ID=0, got nil")
	}
}

func TestValidateConfig_UniversityIDNegative(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].ID = -1
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for negative university ID, got nil")
	}
}

func TestValidateConfig_UniversityNameEmpty(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Name = ""
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for empty university name, got nil")
	}
}

func TestValidateConfig_UniversityAdeUrlEmpty(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].AdeUrl = ""
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for empty adeUrl, got nil")
	}
}

func TestValidateConfig_UniversityAdeUrlInvalid(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].AdeUrl = "not-a-url"
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for invalid adeUrl, got nil")
	}
}

func TestValidateConfig_UniversityAdeProjectIdZero(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].AdeProjectId = 0
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for adeProjectId=0, got nil")
	}
}

func TestValidateConfig_DuplicateUniversityID(t *testing.T) {
	cfg := validConfig()
	univ2 := cfg.Universities[0]
	univ2.Rooms = nil
	univ2.Groups[0].Timetables[0].AdeResources = 999 // avoid adeResources conflict
	cfg.Universities = append(cfg.Universities, univ2)
	err := validateConfig(&cfg)
	if err == nil {
		t.Fatal("expected error for duplicate university ID, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate university id") {
		t.Fatalf("expected 'duplicate university id' in error, got: %v", err)
	}
}

func TestValidateConfig_RoomAdeResourcesZero(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Rooms[0].AdeResources = 0
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for room AdeResources=0, got nil")
	}
}

func TestValidateConfig_RoomLabelEmpty(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Rooms[0].Label = ""
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for empty room label, got nil")
	}
}

func TestValidateConfig_GroupIDZero(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Groups[0].ID = 0
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for group ID=0, got nil")
	}
}

func TestValidateConfig_GroupNameEmpty(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Groups[0].Name = ""
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for empty group name, got nil")
	}
}

func TestValidateConfig_GroupNoTimetables(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Groups[0].Timetables = []TimetableConfig{}
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for group with no timetables, got nil")
	}
}

func TestValidateConfig_TimetableAdeResourcesZero(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Groups[0].Timetables[0].AdeResources = 0
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for timetable AdeResources=0, got nil")
	}
}

func TestValidateConfig_TimetableYearZero(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Groups[0].Timetables[0].Year = 0
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for timetable Year=0, got nil")
	}
}

func TestValidateConfig_TimetableLabelEmpty(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Groups[0].Timetables[0].Label = ""
	if err := validateConfig(&cfg); err == nil {
		t.Fatal("expected error for empty timetable label, got nil")
	}
}

func TestValidateConfig_DuplicateAdeResourcesRoomAndTimetable(t *testing.T) {
	cfg := validConfig()
	// room has AdeResources=10, timetable has AdeResources=20
	// set timetable to same as room
	cfg.Universities[0].Groups[0].Timetables[0].AdeResources = 10
	err := validateConfig(&cfg)
	if err == nil {
		t.Fatal("expected error for duplicate adeResources across room and timetable, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate adeResources") {
		t.Fatalf("expected 'duplicate adeResources' in error, got: %v", err)
	}
}

func TestValidateConfig_DuplicateGroupID(t *testing.T) {
	cfg := validConfig()
	group2 := cfg.Universities[0].Groups[0]
	group2.Timetables = []TimetableConfig{
		{AdeResources: 999, Year: 2025, Label: "Timetable B"},
	}
	cfg.Universities[0].Groups = append(cfg.Universities[0].Groups, group2)
	err := validateConfig(&cfg)
	if err == nil {
		t.Fatal("expected error for duplicate group ID, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate group id") {
		t.Fatalf("expected 'duplicate group id' in error, got: %v", err)
	}
}

func TestValidateConfig_EmptyRoomsAllowed(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Rooms = nil
	if err := validateConfig(&cfg); err != nil {
		t.Fatalf("expected no error for nil rooms, got: %v", err)
	}
}

func TestValidateConfig_EmptyGroupsAllowed(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].Groups = nil
	if err := validateConfig(&cfg); err != nil {
		t.Fatalf("expected no error for nil groups, got: %v", err)
	}
}

func TestValidateConfig_CycleValid(t *testing.T) {
	cfg := validConfig()
	cfg.Universities[0].AdeProjectId = 0
	cfg.Universities[0].AdeProjectIdCycle = &AdeProjectIdCycleConfig{
		StartYear:  2025,
		SplitMonth: 7,
		Cycle:      []int{2, 1},
	}
	if err := validateConfig(&cfg); err != nil {
		t.Fatalf("expected no error for valid cycle, got: %v", err)
	}
}

func TestValidateConfig_CycleInvalid(t *testing.T) {
	tests := []struct {
		name  string
		cycle *AdeProjectIdCycleConfig
	}{
		{"empty cycle", &AdeProjectIdCycleConfig{StartYear: 2025, SplitMonth: 7, Cycle: []int{}}},
		{"startYear=0", &AdeProjectIdCycleConfig{StartYear: 0, SplitMonth: 7, Cycle: []int{2, 1}}},
		{"cycle value=0", &AdeProjectIdCycleConfig{StartYear: 2025, SplitMonth: 7, Cycle: []int{2, 0}}},
		{"splitMonth=0", &AdeProjectIdCycleConfig{StartYear: 2025, SplitMonth: 0, Cycle: []int{2, 1}}},
		{"splitMonth=13", &AdeProjectIdCycleConfig{StartYear: 2025, SplitMonth: 13, Cycle: []int{2, 1}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validConfig()
			cfg.Universities[0].AdeProjectId = 0
			cfg.Universities[0].AdeProjectIdCycle = tc.cycle
			if err := validateConfig(&cfg); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestGetProjectId_Alternating(t *testing.T) {
	cfg := &AdeProjectIdCycleConfig{StartYear: 2025, SplitMonth: 7, Cycle: []int{2, 1}}

	tests := []struct {
		date string
		want int
	}{
		{"2025-07-01", 2},
		{"2026-06-30", 2},
		{"2026-07-01", 1},
		{"2027-06-30", 1},
		{"2027-07-01", 2},
	}

	for _, tc := range tests {
		now, _ := time.Parse("2006-01-02", tc.date)
		got := cfg.GetProjectId(now)
		if got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.date, got, tc.want)
		}
	}
}

func TestGetEffectiveProjectId(t *testing.T) {
	now, _ := time.Parse("2006-01-02", "2025-07-01")

	t.Run("cycle", func(t *testing.T) {
		univ := &UniversityConfig{
			AdeProjectId: 99,
			AdeProjectIdCycle: &AdeProjectIdCycleConfig{
				StartYear: 2025, SplitMonth: 7, Cycle: []int{2, 1},
			},
		}
		if got := univ.GetEffectiveProjectId(now); got != 2 {
			t.Errorf("got %d, want 2", got)
		}
	})

	t.Run("fallback", func(t *testing.T) {
		univ := &UniversityConfig{AdeProjectId: 42}
		if got := univ.GetEffectiveProjectId(now); got != 42 {
			t.Errorf("got %d, want 42", got)
		}
	})
}
