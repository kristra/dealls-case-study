package utils

import "time"

func CountWeekdays(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}

	workingDays := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() >= time.Monday && d.Weekday() <= time.Friday {
			workingDays++
		}
	}
	return workingDays
}
