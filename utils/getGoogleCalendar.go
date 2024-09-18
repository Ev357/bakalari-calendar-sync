package utils

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

func getGoogleCalendarEvents(srv *calendar.Service) (*calendar.Events, error) {
	now := time.Now()
	startOfWeek := getStartOfWeek(now).Format(time.RFC3339)
	endOfTwoWeeks := getEndOfTwoWeeks(now).Format(time.RFC3339)

	return srv.Events.List("primary").
		SingleEvents(true).TimeMin(startOfWeek).TimeMax(endOfTwoWeeks).PrivateExtendedProperty("forBakalariCalendarSync=true").Do()

}

func getStartOfWeek(t time.Time) time.Time {
	return t.AddDate(0, 0, -getIntWeek(t)).Truncate(24 * time.Hour)
}

func getEndOfTwoWeeks(t time.Time) time.Time {
	daysUntilEndOfWeek := 13 - getIntWeek(t)
	endOfTwoWeeks := t.AddDate(0, 0, daysUntilEndOfWeek).Truncate(24 * time.Hour)

	return endOfTwoWeeks.Add(24*time.Hour - time.Nanosecond)
}

func getIntWeek(t time.Time) int {
	weekday := int(t.Weekday()) - 1

	if weekday == -1 {
		weekday = 6
	}

	return weekday
}
