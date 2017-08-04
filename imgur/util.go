package imgur

import (
	"net/http"
	"io"
	"github.com/ishanjain28/imgur-bot/log"
)

func createError(statusCode int, method, error, request string) *iError {
	ierr := &iError{
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

func (i *iError) Error() string {
	return i.Data.Method + "-" + i.Data.Error
}

func makeAuthorisedRequest(method, url string, data io.Reader) (string, *iError) {

	client := &http.Client{}

	req, err := http.NewRequest(method, hostaddr+url, data)
	if err != nil {
		log.Warn.Println("error in creating upload image req", err.Error())
		return nil, createError(0, method, "error in creating request", url)
	}
}
