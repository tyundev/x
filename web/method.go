package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/reiwav/x/mlog"
	"github.com/reiwav/x/rest"
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

func PostFormRequest(url string, header, params, paramFiles map[string]string, res interface{}) (int, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	//prepare the reader instances to encode
	values := map[string]io.Reader{}
	for key, val := range paramFiles {
		values[key] = mustOpen(val)
	}
	for key, val := range params {
		values[key] = strings.NewReader(val)
	}
	//res, err := PostFileResty(paramFiles,params)
	code, err := Upload(client, url, header, values, &res)
	return code, err
}

func Upload(client *http.Client, url string, header map[string]string, values map[string]io.Reader, objRes interface{}) (statusCode int, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}

		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return 500, err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	for key, val := range header {
		fmt.Println(key, val)
		req.Header.Set(key, val)
	}

	// Submit the request
	res, err := client.Do(req)
	fmt.Println(err, res)
	if err != nil {
		return
	}
	defer res.Body.Close()
	return res.StatusCode, json.NewDecoder(res.Body).Decode(&objRes)
}

func mustOpen(f string) *os.File {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		panic(rest.BadRequest("file: " + f + " not exist"))
	}
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
