package main

import (
	"context"
	"fmt"
	"strings"
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

	for _, day := range bakalariCalendar {
		for _, event := range day {
			googleEvent, err := findGoogleEvent(googleCalendar, event)

			if err != nil {
				panic(err)
			}

			if googleEvent != nil {
				switch event.status {
				case "normal":
					// Check if the event is still correct
				default:
					err := srv.Events.Delete(config.account, googleEvent.Id).Do()
					if err != nil {
						panic(err)
					}
				}
			} else {
				_, err := srv.Events.Insert(config.account, getClassEvent(event)).Do()
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func getClassEvent(class Class) *calendar.Event {
	summary := class.name
	location := class.room
	description := ""
	if class.theme != "" {
		description += fmt.Sprintf("Theme: %s\n", class.theme)
	}
	description += fmt.Sprintf("Teacher: %s\n", class.teacher)
	if len(class.homeworks) > 0 {
		description += fmt.Sprintf("Homeworks: %s\n", strings.Join(class.homeworks, ", "))
	}

	start := time.Date(class.date.Year(), class.date.Month(), class.date.Day(), class.from.Hour(), class.from.Minute(), 0, 0, class.from.Location()).Format(time.RFC3339)
	end := time.Date(class.date.Year(), class.date.Month(), class.date.Day(), class.to.Hour(), class.to.Minute(), 0, 0, class.to.Location()).Format(time.RFC3339)

	return &calendar.Event{
		Summary:     summary,
		Location:    location,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start,
		},
		End: &calendar.EventDateTime{
			DateTime: end,
		},
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
