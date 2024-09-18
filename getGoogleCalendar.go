package main

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

func getGoogleCalendarEvents(config *Config, srv *calendar.Service) (*calendar.Events, error) {
	now := time.Now()
	startOfWeek := getStartOfWeek(now).Format(time.RFC3339)
	endOfWeek := getEndOfWeek(now).Format(time.RFC3339)

	return srv.Events.List("primary").
		SingleEvents(true).TimeMin(startOfWeek).TimeMax(endOfWeek).PrivateExtendedProperty("forBakalariCalendarSync=true").Do()

}

func getStartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	diff := weekday - 1

	return t.AddDate(0, 0, -diff).Truncate(24 * time.Hour)
}

func getEndOfWeek(t time.Time) time.Time {
	daysUntilEndOfWeek := time.Saturday - t.Weekday()

	endOfWeek := t.AddDate(0, 0, int(daysUntilEndOfWeek)+1)

	endOfDay := time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, endOfWeek.Location())

	return endOfDay
}
