package client

import (
	"regexp"
	"testing"
)

func TestGetSchedule(t *testing.T) {

	schedule, cron := GetSchedule("hourly")

	re := regexp.MustCompile(`([A-Z][^\s]*)`)
	matched := re.MatchString(schedule)
	if matched == false {
		t.Error("Didn't return a Titled string")
	}

	matchCronStr := regexCron.MatchString(cron)
	if matchCronStr == false {
		t.Error("Invalid cron string")
	}
}

func TestParseBoolOrDefault(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		defaultVal  bool
		expectedVal bool
		expectError bool
	}{
		{"Empty string, default false", "", false, false, false},
		{"Empty string, default true", "", true, true, false},
		{"Value 'true', default false", "true", false, true, false},
		{"Value 'false', default true", "false", true, false, false},
		{"Invalid value", "invalid", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBoolOrDefault(tt.value, tt.defaultVal)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseBoolOrDefault() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if result != tt.expectedVal {
				t.Errorf("ParseBoolOrDefault() = %v, want %v", result, tt.expectedVal)
			}
		})
	}
}
