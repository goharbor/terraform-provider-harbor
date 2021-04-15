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
