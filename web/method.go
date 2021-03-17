package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	if statusCode != http.StatusOK {
		return resp.StatusCode, json.NewDecoder(resp.Body).Decode(&objErr)
	}
	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(&objRes)
}
