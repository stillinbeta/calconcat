package calconcat

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"
)

const (
    vevent_probable_size = 20
    vevent_start = "BEGIN:VEVENT" // RFC 5545 ;)
    vevent_end = "END:VEVENT"
    vtimezone_start = "BEGIN:VTIMEZONE"
    vtimezone_end = "END:VTIMEZONE"
)

type NoMoreVevents struct{}

func (_ NoMoreVevents) Error() string {
    return "That's all the vevents we got"
}

type Vevent struct {
    Err error
    Vevent string
}

type Vtimezone {
    TZID string
    VTimezone string
}

func GetVevents (stream io.Reader, vevents chan Vevent, vtimezones chan Vtimezone) {
    in_vevent := false
    in_vtimezone := false
    tzid := ""
    buf := make([]string, 0, vevent_probable_size)
    scanner := bufio.NewScanner(stream)

    // Sax-y Vevent parser
    for scanner.Scan() {
        line := scanner.Text()
        if in_vevent {
            buf = append(buf, line)
            if line == vevent_end {
                vevents <- Vevent{nil, strings.Join(buf, "\r\n")}

                // Start all over again
                buf = make([]string, 0, vevent_probable_size)
                in_vevent = false
            }
        } else if in_vtimezone {
            buf = append(buf, line)
            if line == vtimezone_end {
                vtimezones <- Vtimezone{tzid, strings.Join(buf, "\r\n")}

                tzid := ""
                buf = make([]string, 0, vevent_probable_size)
                in_vtimezone = false
        } else if line == vevent_start {
            in_vevent = true
            buf = append(buf, line)
        } else if line == vtimezone_start {
            in_vtimezone = true
            buf = append(buf, line)
        }

    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
        c <- Vevent{NoMoreVevents{}, ""}
    }
    c <- Vevent{NoMoreVevents{}, ""}

}

