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
	err := tplIndex.ExecuteWriter(pongo2.Context{
		"timetables": timetables,
	}, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) CSV(w web.ResponseWriter, r *web.Request) {
	file := r.PathParams["file"]
	lastSlash := strings.LastIndex(file, "/")
	if lastSlash != -1 {
		file = file[lastSlash:]
	}

	f, err := os.Open("timeplaner/v2015/" + file)
	if err != nil {
		return
	}

	defer f.Close()
	io.Copy(w, f)
}

func (c *Context) ICal(w web.ResponseWriter, r *web.Request) {
	file := r.PathParams["file"]
	lastSlash := strings.LastIndex(file, "/")
	if lastSlash != -1 {
		file = file[lastSlash:]
	}

	file = file[0 : len(file)-5]

	data, err := utils.ToICal(fmt.Sprintf("timeplaner/v2015/%v.csv", file))
	if err != nil {
		return
	}

	w.Write([]byte(data))
}

func startServer() {
	files, err := ioutil.ReadDir("timeplaner/v2015")
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
		Get("/csv/:file", (*Context).CSV).
		Get("/ical/:file", (*Context).ICal)
	http.ListenAndServe("localhost:15103", router) // Start the server!
}
