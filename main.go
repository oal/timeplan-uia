package main

import (
	"flag"

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
}
