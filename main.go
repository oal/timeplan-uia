package main

import (
	"flag"

	"github.com/oal/timeplan-uia/utils"
)

var serve = flag.Bool("serve", false, "Start server")
var update = flag.Bool("update", false, "Update time tables")
var url = flag.String("url", "", "Download a single calendar based on a URL")

func main() {
	flag.Parse()

	if *serve {
		startServer()
	} else if *update {
		utils.UpdateTimetables()
	} else if *url != "" {
		utils.UpdateSingleURL(*url)
	}
}
