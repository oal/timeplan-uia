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

const URL = "http://timeplan.uia.no/swsuiav/public/no/default.aspx"
const FOLDER = "timeplaner/v2015"

type Department struct {
	Name string
	Code string
}

var departments = []*Department{}

func UpdateTimetables() {
	os.MkdirAll(FOLDER, 0777)

	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}

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
		"__VIEWSTATE":          []string{viewstate},
		"__VIEWSTATEGENERATOR": []string{viewstateGen},
		"__EVENTVALIDATION":    []string{eventValidation},
		"tLinkType":            []string{linkType},
		"dlObject":             []string{department.Code},
		"lbWeeks":              []string{lbWeeks},
		"lbDays":               []string{lbDays},
		"RadioType":            []string{radioType},
		"bGetTimetable":        []string{"Vis+timeplan"},
		"tWildcard":            []string{""},
		"__EVENTTARGET":        []string{""},
		"__EVENTARGUMENT":      []string{""},
		"__LASTFOCUS":          []string{""},
	}

	resp, err := client.PostForm(URL, data)
	if err != nil {
		return err
	}

	generateCSV(resp.Body)
	return nil
}

func generateCSV(r io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
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
		return nil
	}

	title := doc.Find("p.title i").Text()
	csv := []byte(strings.Join(lines, "\r\n"))
	err = ioutil.WriteFile(FOLDER+"/"+title+".csv", csv, 0664)
	if err != nil {
		return err
	}

	fmt.Println("Created timetable for", title)
	return nil
}
