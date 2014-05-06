package main

import (
	"fmt"
	"github.com/stillinbeta/calconcat/calconcat"
	"log"
	"net/http"
)

const (
	CONFIG_FILE = "config.json.example"
)

func getVevents(url string, c chan calconcat.Vevent) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		//w.WriteHeader(http.StatusBadGateway)
		//fmt.Fprintln(w, "Sorry, failed to retrieve calendar")
		log.Printf("Couldn't retrieve upstream URL: %v", err)
		c <- calconcat.Vevent{err, ""}
	} else {
		calconcat.GetVevents(resp.Body, c)
	}
}

type iCalendarHandler struct {
	Config *calconcat.CalenderConfigMap
}

func (ich iCalendarHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {

	key := r.URL.Path[1:] // Strip leading slash
	calList, ok := (*ich.Config)[key]
	if !ok {
		log.Printf("Couldn't find a calendar list for URL %v", key)
		w.WriteHeader(404)
		w.Header().Add("Content-type", "text/plain")
		fmt.Fprintf(w, "NOT FOUND")
		return
	}

	outstanding := len(calList.CalendarList)
	c := make(chan calconcat.Vevent, 1)

	for _, url := range calList.CalendarList {
		go getVevents(url, c)
	}

	w.Header().Add("Content-type", "text/calendar")

	for vevent := range c {
		if vevent.Err != nil {
			outstanding--
			if outstanding == 0 {
				break
			}
		} else {
			fmt.Fprintf(w, vevent.Vevent)
		}
	}
}


func main() {
	config, err := calconcat.ParseConfig(CONFIG_FILE)
	if err != nil {
		log.Fatal("Failed to read config file", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", iCalendarHandler{
		&(config.Calendars),
	})
	listenOn := fmt.Sprintf("%v:%v", config.ListenOn, config.Port)
	log.Printf("Listening on %v", listenOn)
	err = http.ListenAndServe(listenOn, mux)
	if err != nil {
		log.Fatal("Couldn't start server: ", err)
	}
}
