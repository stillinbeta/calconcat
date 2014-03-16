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
)

type NoMoreVevents struct{}

func (_ NoMoreVevents) Error() string {
    return "That's all the vevents we got"
}

type Vevent struct {
    Err error
    Vevent string
}

func GetVevents (stream io.Reader, c chan Vevent) {
    in_vevent := false
    buf := make([]string, 0, vevent_probable_size)
    scanner := bufio.NewScanner(stream)
    for scanner.Scan() {
        line := scanner.Text()
        if in_vevent {
            buf = append(buf, line)
            if line == vevent_end {
                c <- Vevent{nil, strings.Join(buf, "\r\n")}

                // Start all over again
                buf = make([]string, 0, vevent_probable_size)
                in_vevent = false
            }
        } else if line == vevent_start {
            in_vevent = true
            buf = append(buf, line)
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
        c <- Vevent{NoMoreVevents{}, ""}
    }
    c <- Vevent{NoMoreVevents{}, ""}

}

