package calconcat

import (
    "encoding/json"
    "os"
    "log"
)

type calendarConfig struct {
    CalendarList []string `json:"calendars"`
    Title string          `json:"title"`
    Timezone string       `json:"tz"`
}

type CalenderConfigMap map[string]calendarConfig

type configFile struct {
	Calendars CalenderConfigMap `json:"calendars"`
	Port      int                 `json:"port"`
	ListenOn  string              `json:"listen_on"`
}

func ParseConfig(filename string) (*configFile, error) {
	file, err := os.Open(filename)
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
