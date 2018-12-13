package health

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type WebCheck struct {
	Name string
	URL  string
	Type string
}

func (wc WebCheck) Check() *DependencyInfo {
	var err error
	var t float64
	var ver, rev string

	state := DependencyState{}
	sTime := time.Now()
	resp, err := http.Get(wc.URL)
	t = time.Since(sTime).Seconds()

	if err != nil {
		state.Status = CRIT
		state.Details = "Error connecting to `" + wc.URL + "`: " + err.Error()
	} else {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			state.Status = WARN
			state.Details = "Error reading response body: " + err.Error()
		} else {
			var jsonData map[string]string
			if err = json.Unmarshal(data, &jsonData); err != nil {
				state.Status = WARN
				state.Details = "Error parsing response body: " + err.Error()
			} else {
				state.Status = jsonData["status"]
				ver = jsonData["version"]
				rev = jsonData["revision"]
			}
		}
	}
	return &DependencyInfo{
		Name:         wc.Name,
		Type:         wc.Type,
		Version:      ver,
		Revision:     rev,
		State:        state,
		ResponseTime: t,
	}
}

func (wc WebCheck) GetName() string {
	return wc.Name
}

func (wc WebCheck) GetType() string {
	return wc.Type
}
