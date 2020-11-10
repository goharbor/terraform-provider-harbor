package client

func GetSchedule(schedule string) (typefmt string, cronfmt string) {
	var TypeStr string
	var CronStr string

	switch schedule {
	case "hourly", "Hourly":
		TypeStr = "Hourly"
		CronStr = "0 0 * * * *"
	case "daily", "Daily":
		TypeStr = "Daily"
		CronStr = "0 0 0 * * *"
	case "weekly", "Weekly":
		TypeStr = "Weekly"
		CronStr = "0 0 0 * * 0"

	}

	return TypeStr, CronStr
}
