package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func Sync(config *Config) error {
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{calendar.CalendarScope},
		Endpoint:     google.Endpoint,
	}

	tokenSource := oauthConfig.TokenSource(context.TODO(), &oauth2.Token{
		RefreshToken: config.RefreshToken,
	})

	client := oauth2.NewClient(context.TODO(), tokenSource)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	bakalariCalendar, err := getBakalariCalendar(config)

	if err != nil {
		return err
	}

	googleCalendar, err := getGoogleCalendarEvents(srv)

	if err != nil {
		return err
	}

	for _, day := range bakalariCalendar {
		for _, event := range day {
			googleEvent, err := findGoogleEvent(*googleCalendar, event, srv)

			if err != nil {
				return err
			}

			switch event.status {
			case "normal":
				bakalariEvent, err := getClassEvent(event)

				if err != nil {
					return err
				}

				if googleEvent != nil {
					if isEventDifferent(googleEvent, bakalariEvent) {
						if _, err := srv.Events.Patch("primary", googleEvent.Id, bakalariEvent).Do(); err != nil {
							return err
						}
					}
				} else {
					if _, err := srv.Events.Insert("primary", bakalariEvent).Do(); err != nil {
						return err
					}
				}
			default:
				if googleEvent != nil {
					if err := srv.Events.Delete("primary", googleEvent.Id).Do(); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func isEventDifferent(googleEvent *calendar.Event, bakalariEvent *calendar.Event) bool {
	if googleEvent.Summary == bakalariEvent.Summary && googleEvent.Description == bakalariEvent.Description && googleEvent.Location == bakalariEvent.Location {
		return false
	}

	return true
}

func getClassEvent(class Class) (*calendar.Event, error) {
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
		return nil, err
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
	}, nil
}

func findGoogleEvent(googleCalendar calendar.Events, class Class, srv *calendar.Service) (*calendar.Event, error) {
	googleEvents := []*calendar.Event{}

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
			googleEvents = append(googleEvents, event)
		}
	}

	for index, event := range googleEvents {
		if index == len(googleEvents)-1 {
			return event, nil
		}

		if err := srv.Events.Delete("primary", event.Id).Do(); err != nil {
			return event, nil
		}
	}

	return nil, nil
}
