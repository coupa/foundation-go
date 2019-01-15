package health

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type WebCheck struct {
	Name string
	URL  string
	Type string
	//ExpectedStatusCode should only be set when you expect a specific response code from
	//the health check. If it is set to a code > 0, when the response code is the same,
	//the dependency state status will be OK; otherwise it will be CRIT. I.e., if
	//you set it to 404, then when the response code is 404, the status will be OK.
	//This is mainly for checking some web server that does not implement a proper
	//health endpoint and you have to compromise the approach to check it.
	ExpectedStatusCode int
}

func (wc WebCheck) Check() *DependencyInfo {
	var err error
	var t float64
	var ver, rev string

	state := DependencyState{Status: OK}
	sTime := time.Now()
	resp, err := http.Get(wc.URL)
	t = time.Since(sTime).Seconds()

	if err != nil {
		state.Status = CRIT
		state.Details = "Error connecting to `" + wc.URL + "`: " + err.Error()
	} else {
		code := resp.StatusCode
		if code < 300 {
			defer resp.Body.Close()
			ver, rev = wc.parseBody(resp, &state)
		} else if code < 400 {
			//3xx redirect code. Set WARN
			state.Status = WARN
			state.Details = fmt.Sprintf("Status code: %d, %s redirected", code, wc.URL)
		} else {
			state.Status = CRIT
			state.Details = fmt.Sprintf("Status code: %d, error checking %s", code, wc.URL)
		}
		wc.verifyStatusCode(code, &state)
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

//parseBody checks the response body and set the DependencyState data. It returns
//the version and revision if it is able to find it from the response body.
func (wc WebCheck) parseBody(resp *http.Response, state *DependencyState) (string, string) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if wc.Type != TypeThirdParty {
			//When the server type is "internal" or "service", assume that it would
			//follow the microservice standard, so WARN on error
			state.Status = WARN
		}
		state.Details = "Unable to read the response body: " + err.Error()
		return "", ""
	}
	var jsonData map[string]string
	if err = json.Unmarshal(data, &jsonData); err != nil {
		if wc.Type != TypeThirdParty {
			//When the server type is "internal" or "service", assume that it would
			//follow the microservice standard, so WARN on error
			state.Status = WARN
		}
		state.Details = "Response body is not a key-value JSON object: " + err.Error()
		return "", ""
	}
	state.Status = jsonData["status"]
	return jsonData["version"], jsonData["revision"]
}

//verifyStatusCode checks if ExpectedStatusCode
func (wc WebCheck) verifyStatusCode(statusCode int, state *DependencyState) {
	if wc.ExpectedStatusCode <= 0 {
		return
	}
	if wc.ExpectedStatusCode == statusCode {
		state.Status = OK
		return
	}
	state.Status = CRIT
	state.Details = fmt.Sprintf("Expected status code %d but got %d from %s; ", wc.ExpectedStatusCode, statusCode, wc.URL) + state.Details
}
