package handler

import (
	"fmt"
	"hash/crc32"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	proto "github.com/micro/go-os/event/proto"
)

var (
	colours = []string{"blue", "green", "yellow", "purple", "orange"}
)

type sortedRecords struct {
	records []*proto.Record
}

func (s sortedRecords) Len() int {
	return len(s.records)
}

func (s sortedRecords) Less(i, j int) bool {
	return s.records[i].Timestamp > s.records[j].Timestamp
}

func (s sortedRecords) Swap(i, j int) {
	s.records[i], s.records[j] = s.records[j], s.records[i]
}

func colour(s string) string {
	return colours[crc32.ChecksumIEEE([]byte(s))%uint32(len(colours))]
}

func distanceOfTime(minutes float64) string {
	switch {
	case minutes < 1:
		return fmt.Sprintf("%d secs", int(minutes*60))
	case minutes < 59:
		return fmt.Sprintf("%d minutes", int(minutes))
	case minutes < 90:
		return "about an hour"
	case minutes < 120:
		return "almost 2 hours"
	case minutes < 1080:
		return fmt.Sprintf("%d hours", int(minutes/60))
	case minutes < 1680:
		return "about a day"
	case minutes < 2160:
		return "more than a day"
	case minutes < 2520:
		return "almost 2 days"
	case minutes < 2880:
		return "about 2 days"
	default:
		return fmt.Sprintf("%d days", int(minutes/1440))
	}

	return ""
}

func timeAgo(t int64) string {
	d := time.Unix(t, 0)
	timeAgo := ""
	startDate := time.Now().Unix()
	deltaMinutes := float64(startDate-d.Unix()) / 60.0
	if deltaMinutes <= 523440 { // less than 363 days
		timeAgo = fmt.Sprintf("%s ago", distanceOfTime(deltaMinutes))
	} else {
		timeAgo = d.Format("2 Jan")
	}

	return timeAgo
}

func hostPath(r *http.Request) string {
	if path := r.Header.Get("X-Micro-Web-Base-Path"); len(path) > 0 {
		return path
	}
	return "/"
}

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/search", Search)
	r.HandleFunc("/latest", Latest)
	r.HandleFunc("/event/{id}", Event)
	return r
}
