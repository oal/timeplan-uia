package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client *http.Client

const URL = "http://timeplan.uia.no/swsuiah/public/no/default.aspx"
const FOLDER = "timeplaner/h2015"

type Department struct {
	Name string
	Code string
}

var departments = []*Department{}

func UpdateTimetables() {
	os.MkdirAll(FOLDER, 0777)

	doc, err := loadDepartments()
	if err != nil {
		panic(err)
	}

	for _, dep := range departments {
		err = loadTimetable(doc, dep)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second) // Give the server some room to breathe
	}
}

// UpdateSingleURL allows users to pass in a URL to a timetable on UiA's website and export only this file to .csv.
func UpdateSingleURL(url string) {
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}

	title, csv, err := generateCSV(resp.Body)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(title+".csv", []byte(csv), 0664)
	if err != nil {
		panic(err)
	}
}

func loadDepartments() (*goquery.Document, error) {
	resp, err := client.Get(URL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	doc.Find("#dlObject option").Each(func(i int, sel *goquery.Selection) {
		val, _ := sel.Attr("value")
		departments = append(departments, &Department{sel.Text(), val})
	})

	return doc, nil
}

func loadTimetable(doc *goquery.Document, department *Department) error {
	viewstate, _ := doc.Find("#__VIEWSTATE").Attr("value")
	viewstateGen, _ := doc.Find("#__VIEWSTATEGENERATOR").Attr("value")
	eventValidation, _ := doc.Find("#__EVENTVALIDATION").Attr("value")
	linkType, _ := doc.Find("#tLinkType").Attr("value")
	lbWeeks, _ := doc.Find("#lbWeeks option").Eq(1).Attr("value")
	lbDays, _ := doc.Find("#lbDays option").Eq(0).Attr("value")
	radioType, _ := doc.Find("#RadioType_0").Attr("value")

	data := url.Values{
		"__VIEWSTATE":          {viewstate},
		"__VIEWSTATEGENERATOR": {viewstateGen},
		"__EVENTVALIDATION":    {eventValidation},
		"tLinkType":            {linkType},
		"dlObject":             {department.Code},
		"lbWeeks":              {lbWeeks},
		"lbDays":               {lbDays},
		"RadioType":            {radioType},
		"bGetTimetable":        {"Vis+timeplan"},
		"tWildcard":            {""},
		"__EVENTTARGET":        {""},
		"__EVENTARGUMENT":      {""},
		"__LASTFOCUS":          {""},
	}

	resp, err := client.PostForm(URL, data)
	if err != nil {
		return err
	}

	title, csv, err := generateCSV(resp.Body)
	if err != nil {
		panic(err)
	}
	return writeTimetable(title, csv)
}

func generateCSV(r io.Reader) (string, string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", "", err
	}

	lines := []string{"Subject,Start Date,Start Time,End Date,End Time,Description,Location"}

	rows := doc.Find("tr.tr2")
	rows.Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")

		dateString := ""
		var from time.Time
		var to time.Time
		var subject string
		var location string
		var description string

		cells.Each(func(j int, s *goquery.Selection) {
			switch j {
			case 1:
				dateString = s.Text() + " 2015"
			case 2:
				times := strings.Split(strings.TrimSpace(s.Text()), "-")
				from, _ = time.Parse("02 Jan 2006 15.04", dateString+" "+times[0])
				to, _ = time.Parse("02 Jan 2006 15.04", dateString+" "+times[1])
			case 3:
				subject = `"` + strings.Replace(s.Text(), `"`, `""`, -1) + `"`
			case 4:
				location = `"` + strings.Replace(s.Text(), `"`, `""`, -1) + `"`
			case 5:
				description = `"` + strings.Replace(s.Text(), `"`, `""`, -1) + `"`
			}
		})

		line := []string{
			strings.Replace(subject, ",", ".", -1),
			from.Format("2006-01-02"),
			from.Format("15:04"),
			to.Format("2006-01-02"),
			to.Format("15:04"),
			description,
			location,
		}

		lines = append(lines, strings.Join(line, ","))
	})

	if len(lines) < 5 {
		return "", "", nil
	}

	title := doc.Find("p.title i").Text()
	csv := strings.Join(lines, "\r\n")
	return title, csv, nil
}

func writeTimetable(title, csv string) error {
	err := ioutil.WriteFile(FOLDER+"/"+title+".csv", []byte(csv), 0664)
	if err != nil {
		return err
	}

	fmt.Println("Created timetable for", title)
	return nil
}

func init() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}
}