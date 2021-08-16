package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"x/mlog"
)

var logMarshal = mlog.NewTagLog("rest_cetm")

func ResParamArrUrlClient(url string, objArray interface{}, objRes interface{}) error {
	result, err := json.Marshal(objArray)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(result))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logMarshal.Errorf(url, "error marshal post")
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&objRes)
}

func GetRequestHeader(url string, header map[string][]string, objRes interface{}) (error, int) {
	req, err := http.NewRequest("GET", url, nil)
	var statusCode int
	if err != nil {
		return err, statusCode
	}
	for key, val := range header {
		req.Header[key] = val
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		if resp != nil {
			statusCode = resp.StatusCode
		}
		return err, statusCode
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(objRes), resp.StatusCode
}

func ResUrlClientGet(url string, objRes interface{}) error {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		logMarshal.Errorf(url, "error marshal post")
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&objRes)
}

func MethodGet(url string, objRes interface{}) (statusCode int, err error) {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		logMarshal.Errorf(url, "error marshal get")
		if resp != nil {
			return resp.StatusCode, err
		}
		return 404, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(&objRes)
}

func MethodGetNew(url string, objRes interface{}, objErr interface{}) (statusCode int, err error) {
	resp, err := http.Get(url)
	fmt.Println("==== ĐẾN ĐÂY GEO", err)
	if err != nil {
		logMarshal.Errorf(url, "error marshal get")
		if resp != nil {
			return resp.StatusCode, err
		}
		return 404, err
	}

	defer resp.Body.Close()
	statusCode = resp.StatusCode
	if statusCode != http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&objErr)
		return resp.StatusCode, err
	}
	err = json.NewDecoder(resp.Body).Decode(&objRes)
	fmt.Println("==== ĐẾN ĐÂY", objRes)
	return resp.StatusCode, err
}
