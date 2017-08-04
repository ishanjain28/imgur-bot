package imgur

func createError(statusCode int, method, error, request string) *imgurError {
	ierr := &imgurError{
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

func (i *imgurError) Error() string {
	return i.Data.Method + "-" + i.Data.Error
}
