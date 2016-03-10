package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
	"golang.org/x/net/context"

	event "github.com/micro/event-srv/proto/event"
)

var (
	templateDir = "templates"
	opts        *ace.Options

	eventClient event.EventClient
)

func Init(dir string, e event.EventClient) {
	eventClient = e

	opts = ace.InitializeOptions(nil)
	opts.BaseDir = dir
	opts.DynamicReload = true
	opts.FuncMap = template.FuncMap{
		"TimeAgo": func(t int64) string {
			return timeAgo(t)
		},
		"Timestamp": func(t int64) string {
			return time.Unix(t, 0).Format("02 Jan 06 15:04:05 MST")
		},
		"Colour": func(s string) string {
			return colour(s)
		},
	}
}

func render(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	basePath := hostPath(r)

	opts.FuncMap["URL"] = func(path string) string {
		return filepath.Join(basePath, path)
	}

	tpl, err := ace.Load("layout", tmpl, opts)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", 302)
		return
	}

	if err := tpl.Execute(w, data); err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", 302)
	}
}

// The index page
func Index(w http.ResponseWriter, r *http.Request) {
	rsp, err := eventClient.Search(context.TODO(), &event.SearchRequest{
		Reverse: true,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	sort.Sort(sortedRecords{rsp.Records})

	render(w, r, "index", map[string]interface{}{
		"Latest": rsp.Records,
	})
}

func Latest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limit := 15

	page, err := strconv.Atoi(r.Form.Get("p"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	offset := (page * limit) - limit

	rsp, err := eventClient.Search(context.TODO(), &event.SearchRequest{
		Reverse: true,
		Limit:   int64(limit),
		Offset:  int64(offset),
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	var less, more int
	if len(rsp.Records) == limit {
		more = page + 1
	}

	if page > 1 {
		less = page - 1
	}

	sort.Sort(sortedRecords{rsp.Records})

	render(w, r, "latest", map[string]interface{}{
		"Latest": rsp.Records,
		"Less":   less,
		"More":   more,
	})
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.Form.Get("id")

		if len(id) > 0 {
			http.Redirect(w, r, filepath.Join(hostPath(r), "event/"+id), 302)
			return
		}

		rid := r.Form.Get("rid")
		typ := r.Form.Get("type")

		if len(rid) == 0 && len(typ) == 0 {
			http.Redirect(w, r, filepath.Join(hostPath(r), "search"), 302)
			return
		}

		rsp, err := eventClient.Search(context.TODO(), &event.SearchRequest{
			Id:      rid,
			Type:    typ,
			Reverse: true,
		})
		if err != nil {
			http.Redirect(w, r, filepath.Join(hostPath(r), "search"), 302)
			return
		}

		query := "ID: " + rid + " Type: " + typ

		render(w, r, "results", map[string]interface{}{
			"Query":   query,
			"Results": rsp.Records,
		})
		return
	}
	render(w, r, "search", map[string]interface{}{})
}

func Event(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		http.Redirect(w, r, "/", 302)
		return
	}
	// TODO: limit/offset
	rsp, err := eventClient.Read(context.TODO(), &event.ReadRequest{
		Id: id,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	render(w, r, "event", map[string]interface{}{
		"Id":     id,
		"Record": rsp.Record,
	})
}
