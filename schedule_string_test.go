package scheduler

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNextScheduleDate(t *testing.T) {
	// TIME INTERVAL TESTS
	t.Run("Should schedule a time interval in minutes correctly", func(t *testing.T) {
		result, err := getNextScheduleDate("5 minutes")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(5*time.Minute).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("1 minute")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(time.Minute).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("minute")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(time.Minute).Round(time.Second), result.Round(time.Second))
	})
	t.Run("Should schedule a time interval in hours correctly", func(t *testing.T) {
		result, err := getNextScheduleDate("3 hours")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(3*time.Hour).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("1 hour")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(time.Hour).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("hour")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(time.Hour).Round(time.Second), result.Round(time.Second))
	})
	t.Run("Should schedule a time interval in days correctly", func(t *testing.T) {
		result, err := getNextScheduleDate("2 days")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(2*24*time.Hour).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("1 day")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(24*time.Hour).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("day")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(24*time.Hour).Round(time.Second), result.Round(time.Second))
	})
	t.Run("Should schedule a time interval in months correctly", func(t *testing.T) {
		result, err := getNextScheduleDate("4 months")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(4*30*24*time.Hour).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("1 month")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(30*24*time.Hour).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("month")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(30*24*time.Hour).Round(time.Second), result.Round(time.Second))
	})
	t.Run("Should schedule a time interval in years correctly", func(t *testing.T) {
		result, err := getNextScheduleDate("5 years")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(5*365*24*time.Hour).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("1 year")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(365*24*time.Hour).Round(time.Second), result.Round(time.Second))

		result, err = getNextScheduleDate("year")
		assert.NoError(t, err)
		assert.Equal(t, now().Add(365*24*time.Hour).Round(time.Second), result.Round(time.Second))
	})
	t.Run("Should fail if the duration unit is invalid", func(t *testing.T) {
		_, err := getNextScheduleDate("5 bananas")
		assert.Equal(t, "Failed to parse schedule format '5 bananas', invalid time unit: bananas", err.Error())
	})
	t.Run("Should fail if the time unit is invalid", func(t *testing.T) {
		_, err := getNextScheduleDate("two minutes")
		assert.Contains(t, err.Error(), "Failed to parse schedule format 'two minutes', invalid duration: two")
	})

	// TIME FORMATS TEST
	t.Run("Should set the time for today, if the specified hour did not pass already", func(t *testing.T) {
		now := now()
		scheduledTime := now.Add(2 * time.Hour)

		scheduleString := fmt.Sprintf("%02d:%02d", scheduledTime.Hour(), scheduledTime.Minute())
		result, err := getNextScheduleDate(scheduleString)

		assert.NoError(t, err)

		assert.Equal(t, scheduledTime.Day(), result.Day())
		assert.Equal(t, scheduledTime.Hour(), result.Hour())
		assert.Equal(t, scheduledTime.Minute(), result.Minute())
	})
	t.Run("Should set the time for tomorrow, if the specified hour already passed", func(t *testing.T) {
		now := now()
		scheduledTime := now.Add(-2 * time.Hour)

		scheduleString := fmt.Sprintf("%02d:%02d", scheduledTime.Hour(), scheduledTime.Minute())
		result, err := getNextScheduleDate(scheduleString)

		assert.NoError(t, err)

		assert.Equal(t, scheduledTime.Day()+1, result.Day())
		assert.Equal(t, scheduledTime.Hour(), result.Hour())
		assert.Equal(t, scheduledTime.Minute(), result.Minute())
	})

	// WEEKDAY TESTS
	t.Run("Should schedule for the next weekday correctly when the specified weekday has not passed", func(t *testing.T) {
		now := now()
		scheduledWeekday := now.Add(time.Hour * 48).Weekday()
		scheduledString := strings.ToLower(scheduledWeekday.String())

		result, err := getNextScheduleDate(scheduledString)

		assert.NoError(t, err)
		assert.Equal(t, 2, result.Day()-now.Day())
	})
	t.Run("Should schedule for the next weekday correctly when the specified weekday already passed", func(t *testing.T) {
		now := now()
		scheduledWeekday := now.Add(-time.Hour * 48).Weekday()
		scheduledString := strings.ToLower(scheduledWeekday.String())

		result, err := getNextScheduleDate(scheduledString)

		assert.NoError(t, err)
		assert.Equal(t, 5, result.Day()-now.Day())
	})
	t.Run("Should schedule for the next week if the specified weekday is today", func(t *testing.T) {
		now := now()
		scheduledString := strings.ToLower(now.Weekday().String())

		result, err := getNextScheduleDate(scheduledString)

		assert.NoError(t, err)
		assert.Equal(t, 7, result.Day()-now.Day())
	})

	// WEEKDAY AND TIME TESTS
	t.Run("Should schedule for the next weekday with the correct time when the specified weekday has not passed", func(t *testing.T) {
		now := now()
		scheduledTime := now.Add(time.Hour * 48)

		weekdayString := strings.ToLower(scheduledTime.Weekday().String())
		timeString := fmt.Sprintf("%02d:%02d", scheduledTime.Hour(), scheduledTime.Minute())

		scheduledString := fmt.Sprintf("%s at %s", weekdayString, timeString)

		result, err := getNextScheduleDate(scheduledString)

		assert.NoError(t, err)

		assert.Equal(t, 2, result.Day()-now.Day())
		assert.Equal(t, scheduledTime.Hour(), result.Hour())
		assert.Equal(t, scheduledTime.Minute(), result.Minute())
	})
	t.Run("Should schedule for the next weekday with the correct time when the specified weekday already passed", func(t *testing.T) {
		now := now()
		scheduledTime := now.Add(-time.Hour * 48)

		weekdayString := strings.ToLower(scheduledTime.Weekday().String())
		timeString := fmt.Sprintf("%02d:%02d", scheduledTime.Hour(), scheduledTime.Minute())

		scheduledString := fmt.Sprintf("%s at %s", weekdayString, timeString)

		result, err := getNextScheduleDate(scheduledString)

		assert.NoError(t, err)

		assert.Equal(t, 5, result.Day()-now.Day())
		assert.Equal(t, scheduledTime.Hour(), result.Hour())
		assert.Equal(t, scheduledTime.Minute(), result.Minute())
	})
	t.Run("Should schedule for today if the weekday is today, but the time has not passed already", func(t *testing.T) {
		now := now()
		scheduledTime := now.Add(time.Hour * 2)

		weekdayString := strings.ToLower(scheduledTime.Weekday().String())
		timeString := fmt.Sprintf("%02d:%02d", scheduledTime.Hour(), scheduledTime.Minute())

		scheduledString := fmt.Sprintf("%s at %s", weekdayString, timeString)

		result, err := getNextScheduleDate(scheduledString)

		assert.NoError(t, err)

		assert.Equal(t, now.Day(), result.Day())
		assert.Equal(t, scheduledTime.Hour(), result.Hour())
		assert.Equal(t, scheduledTime.Minute(), result.Minute())
	})
	t.Run("Should schedule for next week if the weekday is today and the time has already passed", func(t *testing.T) {
		now := now()
		scheduledTime := now.Add(-time.Hour * 2)

		weekdayString := strings.ToLower(scheduledTime.Weekday().String())
		timeString := fmt.Sprintf("%02d:%02d", scheduledTime.Hour(), scheduledTime.Minute())

		scheduledString := fmt.Sprintf("%s at %s", weekdayString, timeString)

		result, err := getNextScheduleDate(scheduledString)

		assert.NoError(t, err)

		assert.Equal(t, 7, result.Day()-now.Day())
		assert.Equal(t, scheduledTime.Hour(), result.Hour())
		assert.Equal(t, scheduledTime.Minute(), result.Minute())
	})
}
