package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// weekdays is a mapping of weekday strings to its respective values on the time package
	weekdays = map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}

	// hourMinuteFormat represents the HH:MM time format
	hourMinuteFormat = "15:04"

	// unitToDuration is a mapping of time units to their respective durations
	unitToDuration = map[string]time.Duration{
		"second":  time.Second,
		"seconds": time.Second,
		"minutes": time.Minute,
		"minute":  time.Minute,
		"hours":   time.Hour,
		"hour":    time.Hour,
		"days":    24 * time.Hour,
		"day":     24 * time.Hour,
		"months":  30 * 24 * time.Hour,
		"month":   30 * 24 * time.Hour,
		"years":   365 * 24 * time.Hour,
		"year":    365 * 24 * time.Hour,
	}
)

// getNextScheduleDate parses a time schedule string into the date of the next execution
func getNextScheduleDate(schedule string) (time.Time, error) {
	// Get the current time
	now := now()

	// Split the input schedule string into words
	words := strings.Fields(schedule)

	// Check if the first word is a number and the second word is a valid time unit
	if len(words) == 2 {
		num, err := strconv.Atoi(words[0])
		if err != nil {
			return time.Time{}, fmt.Errorf("Failed to parse schedule format '%s', invalid duration: %s, Error: %v", schedule, words[0], err)
		}

		unit := words[1]
		duration, found := unitToDuration[unit]
		if !found {
			return time.Time{}, fmt.Errorf("Failed to parse schedule format '%s', invalid time unit: %s", schedule, unit)
		}

		nextTime := now.Add(time.Duration(num) * duration)
		return nextTime, nil
	}

	if len(words) == 1 {
		// Check if the input is a weekday
		weekday, found := weekdays[words[0]]
		if found {
			daysUntilNextWeekday := int(weekday-now.Weekday()+7) % 7
			nextTime := now.Add(time.Duration(daysUntilNextWeekday) * 24 * time.Hour)
			nextTime = time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 1, 0, 0, time.UTC)
			// If the specified time has already passed for today, set it for the next week
			if nextTime.Before(now) {
				nextTime = nextTime.Add(time.Hour * 24 * 7)
			}

			return nextTime, nil
		}

		// Check if the input is a duration
		duration, found := unitToDuration[words[0]]
		if found {
			nextTime := now.Add(duration)
			return nextTime, nil
		}

		// Check if the input is a time in HH:MM format
		t, err := time.Parse(hourMinuteFormat, words[0])
		if err == nil {
			nextTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)

			// If the specified time has already passed for today, set it for the next day
			if now.After(nextTime) {
				nextTime = nextTime.Add(24 * time.Hour)
			}

			return nextTime, nil
		}
	}

	// Check if the input is a combination of a weekday and time
	if len(words) == 3 {
		if words[1] == "at" {
			weekday, found := weekdays[words[0]]
			if !found {
				return time.Time{}, fmt.Errorf("Failed to parse schedule format '%s', invalid weekday: %s", schedule, words[0])
			}

			t, err := time.Parse(hourMinuteFormat, words[2])
			if err != nil {
				return time.Time{}, fmt.Errorf("Failed to parse schedule format '%s', invalid duration: %s, Error: %v", schedule, words[0], err)
			}

			daysUntilNextWeekday := int(weekday-now.Weekday()+7) % 7
			nextTime := now.Add(time.Duration(daysUntilNextWeekday) * 24 * time.Hour)
			nextTime = time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
			// If the specified time has already passed for today, set it for the next week
			if nextTime.Before(now) {
				nextTime = nextTime.Add(time.Hour * 24 * 7)
			}
			return nextTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("Invalid schedule format: %s", schedule)
}
