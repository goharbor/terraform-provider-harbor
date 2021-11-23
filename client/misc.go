package client

import (
	"regexp"
	"strings"
)

func GetSchedule(schedule string) (typefmt string, cronfmt string) {
	var TypeStr string
	var CronStr string

	switch strings.ToLower(schedule) {
	case "hourly", "0 0 * * * *":
		TypeStr = "Hourly"
		CronStr = "0 0 * * * *"
	case "daily", "0 0 0 * * *":
		TypeStr = "Daily"
		CronStr = "0 0 0 * * *"
	case "weekly", "0 0 0 * * 0":
		TypeStr = "Weekly"
		CronStr = "0 0 0 * * 0"
	default:
		TypeStr = "Custom"
		CronStr = schedule
	}

	return TypeStr, CronStr
}

var regexCron = regexp.MustCompile(`(?m)((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
