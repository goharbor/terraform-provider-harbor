package client

import (
	"regexp"
	"strings"
)

func GetSchedule(schedule string) (typefmt string, cronfmt string) {
	var TypeStr string
	var CronStr string

	switch strings.ToLower(schedule) {
	case "hourly":
		TypeStr = "Hourly"
		CronStr = "0 0 * * * *"
	case "daily":
		TypeStr = "Daily"
		CronStr = "0 0 0 * * *"
	case "weekly":
		TypeStr = "Weekly"
		CronStr = "0 0 0 * * 0"

	}
	if regexCron.MatchString(schedule) {
		TypeStr = "Custom"
		CronStr = schedule
	}

	return TypeStr, CronStr
}

var regexCron = regexp.MustCompile(`(?m)((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
