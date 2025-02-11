package utils

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

func getBakalariCalendar(config *Config) ([]Day, error) {
	loginUrl := config.Url + "/bakaweb/Login"

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	_, err = client.PostForm(loginUrl, url.Values{
		"username": {config.Username},
		"password": {config.Password},
	})

	if err != nil {
		return nil, err
	}

	urls := []string{
		config.Url + "/bakaweb/next/rozvrh.aspx",
		config.Url + "/bakaweb/next/rozvrh.aspx?s=next",
	}

	events := []Day{}

	for _, url := range urls {
		weekEvents, err := getEvents(client, url)

		if err != nil {
			return nil, err
		}

		events = append(events, weekEvents...)
	}

	return events, nil
}

func getEvents(client *http.Client, calendarUrl string) ([]Day, error) {
	resp, err := client.Get(calendarUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return nil, err
	}

	timetable := dom.GetElementByID(doc, "schedule")

	hours, err := getHours(timetable)

	if err != nil {
		return nil, err
	}

	return getDays(timetable, hours)
}

type Class struct {
	num       int
	from      time.Time
	to        time.Time
	date      time.Time
	teacher   string
	room      string
	homeworks []string
	name      string
	status    string
	theme     string
}

type Day []Class

func getDays(timetable *html.Node, hours []Hour) ([]Day, error) {
	days := []Day{}

	for _, dayContainer := range dom.GetElementsByClassName(timetable, "day-row") {
		day, err := getDay(dayContainer, hours)

		if err != nil {
			return nil, err
		}

		days = append(days, *day)
	}

	return days, nil
}

func getDay(doc *html.Node, hours []Hour) (*Day, error) {
	day := Day{}

	dateNode := dom.QuerySelector(doc, ".day-name > div > span")

	currentYear := time.Now().Year()
	parsedDate, err := time.Parse("2/1", dom.InnerText(dateNode))
	if err != nil {
		return nil, err
	}

	date := time.Date(currentYear, parsedDate.Month(), parsedDate.Day(), 0, 0, 0, parsedDate.Nanosecond(), parsedDate.Location())

	for index, classContainer := range dom.QuerySelectorAll(doc, ".day-row > div > div > span") {
		class, err := getClass(classContainer, index, hours, date)

		if err != nil {
			continue
		}

		day = append(day, *class)
	}

	return &day, nil
}

func getClass(node *html.Node, index int, hours []Hour, date time.Time) (*Class, error) {
	emptyNode := dom.QuerySelector(node, ".empty")

	dayItemNode := dom.QuerySelector(node, "div[data-detail]")

	type DataDetail struct {
		Type      string   `json:"type"`
		Teacher   string   `json:"teacher"`
		Room      string   `json:"room"`
		Homeworks []string `json:"homeworks"`
		Theme     string   `json:"theme"`
	}
	dataDetail := DataDetail{}

	if dayItemNode != nil {
		for _, attr := range dayItemNode.Attr {
			if attr.Key == "data-detail" {
				if err := json.Unmarshal([]byte(attr.Val), &dataDetail); err != nil {
					return nil, err
				}
			}
		}
	}

	if dataDetail.Type == "removed" {
		return &Class{
			num:    hours[index].num,
			from:   hours[index].from,
			to:     hours[index].to,
			date:   date,
			status: "removed",
		}, nil
	}

	if emptyNode != nil {
		return &Class{
			num:    hours[index].num,
			from:   hours[index].from,
			to:     hours[index].to,
			date:   date,
			status: "empty",
		}, nil
	}

	nameNode := dom.QuerySelector(node, ".day-item > div > div > div:nth-child(2)")

	if nameNode == nil {
		return nil, errors.New("Classname not found")
	}

	name := dom.InnerText(nameNode)

	return &Class{
		num:       hours[index].num,
		from:      hours[index].from,
		to:        hours[index].to,
		date:      date,
		teacher:   dataDetail.Teacher,
		room:      dataDetail.Room,
		homeworks: dataDetail.Homeworks,
		theme:     dataDetail.Theme,
		name:      name,
		status:    "normal",
	}, nil
}

func getHours(timetable *html.Node) ([]Hour, error) {
	hoursContainer := dom.GetElementByID(timetable, "hours")

	hours := []Hour{}

	for _, hourContainer := range dom.GetElementsByClassName(hoursContainer, "item") {
		hour, err := getHourData(hourContainer)

		if err != nil {
			return nil, err
		}

		hours = append(hours, *hour)
	}

	return hours, nil
}

type Hour struct {
	num  int
	from time.Time
	to   time.Time
}

func getHourData(doc *html.Node) (*Hour, error) {
	numNode := dom.QuerySelector(doc, ".num")

	num, err := strconv.Atoi(dom.InnerText(numNode))

	if err != nil {
		return nil, err
	}

	const layout = "15:04"

	fromNode := dom.QuerySelector(doc, ".from")

	from, err := time.Parse(layout, dom.InnerText(fromNode))

	if err != nil {
		return nil, err
	}

	toNode := dom.QuerySelector(doc, ".to")

	to, err := time.Parse(layout, dom.InnerText(toNode))

	if err != nil {
		return nil, err
	}

	return &Hour{
		num,
		from,
		to,
	}, nil
}
