package client

import (
	"regexp"
	"strconv"
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
	case "":
		TypeStr = "None"
		CronStr = ""
	default:
		TypeStr = "Custom"
		CronStr = schedule
	}

	return TypeStr, CronStr
}

var regexCron = regexp.MustCompile(`(?m)((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)

func ParseBoolOrDefault(value string, defaultValue bool) (bool, error) {
	if value == "" {
		return defaultValue, nil
	}
	return strconv.ParseBool(value)
}
