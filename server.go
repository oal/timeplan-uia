package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/oal/timeplan-uia/utils"

	"github.com/flosch/pongo2"
	"github.com/gocraft/web"
)

var timetables = []string{}

var tplIndex = pongo2.Must(pongo2.FromFile("templates/index.html"))

type Context struct{}

func (c *Context) Index(w web.ResponseWriter, r *web.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tplIndex.ExecuteWriter(pongo2.Context{
		"timetables": timetables,
	}, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) CSV(w web.ResponseWriter, r *web.Request) {
	file := r.PathParams["file"]
	if len(file) < 5 || file[len(file)-4:] != ".csv" {
		return
	}

	semester := r.PathParams["semester"]
	f, err := os.Open(fmt.Sprintf("timeplaner/%s/%s", semester, file))
	if err != nil {
		return
	}

	defer f.Close()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	io.Copy(w, f)
}

func (c *Context) ICal(w web.ResponseWriter, r *web.Request) {
	file := r.PathParams["file"]
	if len(file) < 5 || file[len(file)-4:] != ".ics" {
		return
	}

	file = file[0 : len(file)-4]
	fmt.Println(file)

	semester := r.PathParams["semester"]
	data, err := utils.ToICal(fmt.Sprintf("timeplaner/%v/%v.csv", semester, file))
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Write([]byte(data))
}

func startServer() {
	files, err := ioutil.ReadDir("timeplaner/h2015")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		name := f.Name()
		name = name[0:strings.LastIndex(name, ".")]
		timetables = append(timetables, name)
	}

	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Get("/", (*Context).Index).
		Get("/csv/:semester/:file", (*Context).CSV).
		Get("/ical/:semester/:file", (*Context).ICal)
	http.ListenAndServe("0.0.0.0:15103", router) // Start the server!
}
