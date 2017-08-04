package imgur

import (
	"net/http"
	"io"
	"github.com/ishanjain28/imgur-bot/log"
	"strings"
	"fmt"
)

func createError(statusCode int, method, error, request string) *IError {
	ierr := &IError{
		Success: false,
		Status:  statusCode,

	}

	if statusCode != 200 && statusCode != 0 && error == "" {
		ierr.Data.Error = "Status code is not 200"
	} else {
		ierr.Data.Error = error
	}

	ierr.Data.Request = request
	ierr.Data.Method = method

	return ierr
}

func (i *IError) Error() string {

	if i != nil {
		return i.Data.Method + ": " + i.Data.Error
	}

	return ""
}

func makeAuthorisedRequest(method, url, token, tokenType, data string) (*http.Response, *IError) {

	fmt.Println(method, url, token, data)
	client := &http.Client{}

	var dataReader io.Reader
	if data != "" {
		dataReader = strings.NewReader(data)
	}

	req, err := http.NewRequest(method, hostaddr+url, dataReader)
	if err != nil {
		log.Warn.Println("error in creating upload image req", err.Error())
		return nil, createError(0, method, "error in creating request", url)
	}

	if tokenType == "client_id" {
		req.Header.Add("authorization", "Client-ID "+token)
	} else {
		req.Header.Add("authorization", "Bearer "+token)
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, createError(resp.StatusCode, method, err.Error(), url)
	}

	if resp.StatusCode != 200 {
		return nil, createError(resp.StatusCode, method, "", url)
	}

	return resp, nil

}

func makeUnauthorisedRequest(method, url, data string) (io.ReadCloser, *IError) {

	client := &http.Client{}

	req, err := http.NewRequest(method, hostaddr+url, strings.NewReader(data))
	if err != nil {
		log.Warn.Println("error in creating upload image req", err.Error())
		return nil, createError(0, method, "error in creating request", url)
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, createError(resp.StatusCode, method, err.Error(), url)
	}

	if resp.StatusCode != 200 {
		return nil, createError(resp.StatusCode, method, "", url)
	}

	return resp.Body, nil

}
