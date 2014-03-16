package calconcat

import (
    "strings"
    "testing"
    "os"
    "time"
)

//import "calconcat/vevent"

const ical_file = "../example.ical"

func checkIsVevent(s string) bool {
    return (strings.HasPrefix(s, "BEGIN:VEVENT") &&
            strings.HasSuffix(s, "END:VEVENT"))
}

func Test_GetVevents(t *testing.T) {

    timeout := make(chan bool, 1)
    go func() {
            time.Sleep(1 * time.Second)
            timeout <- true
    }()

    file, err := os.Open(ical_file)
    if err != nil {
        t.Fatalf("Couldn't open ical file: %v", err)
    }

    channel := make(chan Vevent, 1)

    go GetVevents(file, channel)

    vevent := <-channel
    if vevent.Err != nil {
        t.Fatalf("First vevent had an error")
    }
    if strings.Index(vevent.Vevent, "Bastille") == -1 {
        t.Fatalf("First event isnt' Bastille day (got %v)", vevent.vevent)
    }
    if !checkIsVevent(vevent.Vevent) {
        t.Fatalf("First event isnt' a VEVENT")
    }

    vevent = <-channel
    if vevent.Err != nil {
        t.Fatalf("Second vevent had an error")
    }
    if strings.Index(vevent.Vevent, "Networld+Interop") == -1 {
        t.Fatalf("Second event isn't Networld+Interop conference")
    }
    if !checkIsVevent(vevent.Vevent) {
        t.Fatalf("Second event isnt' a VEVENT")
    }

    select {
    case vevent = <-channel :
        if vevent.Err == nil {
            t.Fatalf("Third vevent didn't error")
        }

    case <-timeout:
        t.Fatalf("Timout waiting for last element")
    }
}
