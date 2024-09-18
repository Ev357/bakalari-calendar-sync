package main

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

func getGoogleCalendarEvents(srv *calendar.Service) (*calendar.Events, error) {
	now := time.Now()
	startOfWeek := getStartOfWeek(now).Format(time.RFC3339)
	endOfWeek := getEndOfWeek(now).Format(time.RFC3339)

	return srv.Events.List("primary").
		SingleEvents(true).TimeMin(startOfWeek).TimeMax(endOfWeek).PrivateExtendedProperty("forBakalariCalendarSync=true").Do()

}

func getStartOfWeek(t time.Time) time.Time {
	daysBeforeStartOfWeek := t.Weekday() - time.Sunday

	return t.AddDate(0, 0, -int(daysBeforeStartOfWeek)).Truncate(24 * time.Hour)
}

func getEndOfWeek(t time.Time) time.Time {
	daysUntilEndOfWeek := time.Saturday - t.Weekday()

	endOfWeek := t.AddDate(0, 0, int(daysUntilEndOfWeek)).Truncate(24 * time.Hour)

	return endOfWeek.Add(24*time.Hour - time.Nanosecond)
}
