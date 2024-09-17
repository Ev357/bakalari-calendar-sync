package main

import (
	"context"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func main() {
	godotenv.Overload()

	config, err := getConfig()

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	srv, err := calendar.NewService(ctx, option.WithCredentialsJSON(config.serviceAccount))
	if err != nil {
		panic(err)
	}

	bakalariCalendar, err := getBakalariCalendar(config)

	if err != nil {
		panic(err)
	}

	googleCalendar, err := getGoogleCalendarEvents(config, srv)

	if err != nil {
		panic(err)
	}

	event := &calendar.Event{
		Summary:     "Google I/O 2015",
		Description: "A chance to hear more about Google's developer products.",
		Start: &calendar.EventDateTime{
			DateTime: "2024-09-17T11:35:00-00:00",
		},
		End: &calendar.EventDateTime{
			DateTime: "2024-09-17T12:20:00-00:00",
		},
	}

	_, err = srv.Events.Insert(config.account, event).Do()

	if err != nil {
		panic(err)
	}

	for _, day := range bakalariCalendar {
		for _, event := range day {
			googleEvent, err := findGoogleEvent(googleCalendar, event)

			if err != nil {
				panic(err)
			}

			if googleEvent != nil {
				fmt.Println(googleEvent.Summary)
			}
		}
	}
}

func findGoogleEvent(googleCalendar []calendar.Event, class Class) (*calendar.Event, error) {
	for _, event := range googleCalendar {
		parsedTime, err := time.Parse(time.RFC3339, event.Start.DateTime)

		if err != nil {
			return nil, err
		}

		h1, m1, s1 := parsedTime.Clock()
		h2, m2, s2 := class.from.Clock()
		y1, M1, d1 := parsedTime.Date()
		y2, M2, d2 := class.date.Date()

		if h1 == h2 && m1 == m2 && s1 == s2 && y1 == y2 && M1 == M2 && d1 == d2 {
			return &event, nil
		}
	}

	return nil, nil
}
