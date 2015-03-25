package main

import (
	"flag"
	"fmt"

	"github.com/oal/timeplan-uia/utils"
)

var serve = flag.Bool("serve", false, "Start server")
var update = flag.Bool("update", false, "Update time tables")

func main() {
	flag.Parse()

	if *serve {
		startServer()
	} else if *update {
		utils.UpdateTimetables()
	}

	ical, err := utils.ToICal("timeplaner/v2015/Barnehagelærerutdanning, bachelor 1. år.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println(ical)
}
