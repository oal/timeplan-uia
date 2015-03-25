package utils

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

func ToICal(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	calName := path[strings.LastIndex(path, "/")+1:]
	calName = calName[0:strings.LastIndex(calName, ".")]
	cal := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//oal//timeplan//EN",
		fmt.Sprintf("X-WR-CALNAME:%v", calName),
		"CALSCALE:GREGORIAN",
		"BEGIN:VTIMEZONE",
		"TZID:Europe/Oslo",
		"END:VTIMEZONE",
	}

	r := csv.NewReader(f)

	// Read first line, as it's just the headers. It is an error if nothing is read.
	_, err = r.Read()
	if err != nil {
		return "", err
	}

	record, err := r.Read()
	for err == nil {
		date, _ := time.Parse("2006-01-02", record[1])
		startTime, _ := time.Parse("15:04", record[2])
		dateStart := date.Add(time.Duration(startTime.Hour())*time.Hour + time.Duration(startTime.Minute())*time.Minute)

		endTime, _ := time.Parse("15:04", record[4])
		dateEnd := date.Add(time.Duration(endTime.Hour())*time.Hour + time.Duration(endTime.Minute())*time.Minute)

		hash := sha1.Sum([]byte(record[3] + record[5] + record[6]))
		uid := fmt.Sprintf("%vZ-%x@timeplan.uia.no", dateStart.Format("20060102T1504"), hash[0:4])

		summary := strings.Replace(record[0], "\n", "\n ", -1)
		description := strings.Replace(record[5], "\n", "\n ", -1)
		location := strings.Replace(record[6], "\n", "\n ", -1)
		cal = append(cal, []string{
			"BEGIN:VEVENT",
			fmt.Sprintf("UID:%v", uid),
			fmt.Sprintf("DTSTART:%v00Z", dateStart.Format("20060102T1504")),
			fmt.Sprintf("DTEND:%v00Z", dateEnd.Format("20060102T1504")),
			fmt.Sprintf("SUMMARY:%v", summary),
			fmt.Sprintf("DESCRIPTION:%v", description),
			fmt.Sprintf("LOCATION:%v", location),
			"END:VEVENT",
		}...)

		record, err = r.Read()
	}
	cal = append(cal, "END:VCALENDAR\n")

	return strings.Join(cal, "\r\n"), nil
}
