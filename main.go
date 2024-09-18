package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry       string `json:"expiry"`
}

func main() {
	godotenv.Overload()

	config, err := getConfig()

	if err != nil {
		panic(err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     config.clientId,
		ClientSecret: config.clientSecret,
		Scopes:       []string{calendar.CalendarScope},
		Endpoint:     google.Endpoint,
	}

	tokenSource := oauthConfig.TokenSource(context.TODO(), &oauth2.Token{
		RefreshToken: config.refreshToken,
	})

	_, err = tokenSource.Token()
	if err != nil {
		panic(err)
	}

	client := oauth2.NewClient(context.TODO(), tokenSource)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}

	bakalariCalendar, err := getBakalariCalendar(config)

	if err != nil {
		panic(err)
	}

	googleCalendar, err := getGoogleCalendarEvents(srv)

	if err != nil {
		panic(err)
	}

	for _, day := range bakalariCalendar {
		for _, event := range day {
			googleEvent, err := findGoogleEvent(*googleCalendar, event)

			if err != nil {
				panic(err)
			}

			switch event.status {
			case "normal":
				if googleEvent != nil {
					// Check if the event is still correct
				} else {
					_, err := srv.Events.Insert("primary", getClassEvent(event)).Do()
					if err != nil {
						panic(err)
					}
				}
			default:
				if googleEvent != nil {
					err := srv.Events.Delete("primary", googleEvent.Id).Do()
					if err != nil {
						panic(err)
					}
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
	timeLocation, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		panic(err)
	}
	start := time.Date(class.date.Year(), class.date.Month(), class.date.Day(), class.from.Hour(), class.from.Minute(), 0, 0, timeLocation).Format(time.RFC3339)
	end := time.Date(class.date.Year(), class.date.Month(), class.date.Day(), class.to.Hour(), class.to.Minute(), 0, 0, timeLocation).Format(time.RFC3339)

	return &calendar.Event{
		Summary:     summary,
		Location:    location,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start,
			TimeZone: "Europe/Prague",
		},
		End: &calendar.EventDateTime{
			DateTime: end,
			TimeZone: "Europe/Prague",
		},
		Reminders: &calendar.EventReminders{
			UseDefault:      false,
			ForceSendFields: []string{"UseDefault"},
		},
		ExtendedProperties: &calendar.EventExtendedProperties{
			Private: map[string]string{
				"forBakalariCalendarSync": "true",
			},
		},
	}
}

func findGoogleEvent(googleCalendar calendar.Events, class Class) (*calendar.Event, error) {
	for _, event := range googleCalendar.Items {
		parsedTime, err := time.Parse(time.RFC3339, event.Start.DateTime)

		if err != nil {
			return nil, err
		}

		h1, m1, s1 := parsedTime.Clock()
		h2, m2, s2 := class.from.Clock()
		y1, M1, d1 := parsedTime.Date()
		y2, M2, d2 := class.date.Date()

		if h1 == h2 && m1 == m2 && s1 == s2 && y1 == y2 && M1 == M2 && d1 == d2 {
			return event, nil
		}
	}

	return nil, nil
}
