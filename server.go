package main

import (
    "os"
    "encoding/json"
    "fmt"
    "net/http"
    "log"
    "github.com/stillinbeta/calconcat/calconcat"
)

const (
    CONFIG_FILE = "config.json"
    LISTEN_ON = "localhost"
    PORT = 8080
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
    Config *map[string][]string
}

func (ich iCalendarHandler) ServeHTTP (
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

    outstanding := len(calList)
    c := make(chan calconcat.Vevent, 1)

    for _, url := range calList {
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

type configFile struct {
    Calendars map[string][]string `json:"calendars"`
    Port int `json:"port"`
    ListenOn string `json:"listen_on"`
}

func parseConfig() (*configFile, error) {
    file, err := os.Open(CONFIG_FILE)
    if err != nil {
        return nil, err
    }

    defer file.Close()

    result := new(configFile)
    decoder := json.NewDecoder(file)
    err = decoder.Decode(result)
    if err != nil {
        log.Printf("Error decoding config file! %v", err)
        return nil, err
    }

    return result, nil

}

func main() {
    config, err := parseConfig()
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
