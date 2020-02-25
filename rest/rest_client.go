package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type RestResponse struct {
	Status int
	Error  error
	Body   []byte
}

type RestClient struct {
	httpClient *http.Client
}

func (self *RestClient) SetHttpClient(httpClient *http.Client) {
	self.httpClient = httpClient
}

func (self *RestClient) GetObject(url string, obj interface{}) RestResponse {
	return self.executeRequest("GET", url, nil, obj)
}

func (self *RestClient) PostObject(url string, obj interface{}) RestResponse {
	return self.executeRequest("POST", url, obj, obj)
}

func (self *RestClient) DeleteObject(url string) RestResponse {
	return self.executeRequest("DELETE", url, nil, nil)
}

func (self *RestClient) PutObject(url string, obj interface{}) RestResponse {
	return self.executeRequest("PUT", url, obj, obj)
}

func (self *RestClient) executeRequest(method string, url string, inputObj interface{}, outputObj interface{}) RestResponse {
	client := self.httpClient
	if client == nil {
		client = &http.Client{}
	}
	ret := RestResponse{}
	var inputBody io.Reader
	if outputObj != nil {
		data, err := json.Marshal(outputObj)
		if err != nil {
			ret.Error = err
			return ret
		}
		inputBody = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, inputBody)
	if err != nil {
		ret.Error = err
		return ret
	}
	resp, err := client.Do(req)
	if err != nil {
		ret.Error = err
		return ret
	}
	ret.Status = resp.StatusCode
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ret.Error = err
		return ret
	}
	ret.Body = data
	if outputObj != nil {
		json.Unmarshal(data, outputObj)
	}
	return ret
}
