package imgur

import (
	"encoding/json"
	"github.com/ishanjain28/imgur-bot/log"
)

func (i *Imgur) Albums(username, accessToken string) (*Album, *IError) {
	//Skipping a page parameter here, It can be used to get only a limited number of albums
	resp, ierr := makeAuthorisedRequest("GET", "/3/account/"+username+"/albums/", accessToken, "access_token", "")
	if ierr != nil {
		return nil, ierr
	}

	b := &Album{}

	err := json.NewDecoder(resp.Body).Decode(b)

	if err != nil {

		log.Warn.Println("Error in unmarshalling, Response might be an error", err.Error())
		ierr := &IError{}

		err := json.NewDecoder(resp.Body).Decode(&ierr)
		if err != nil {
			log.Warn.Println("Error in unmarshalling", err.Error())

			return nil, createError(resp.StatusCode, "GET", err.Error(), "/3/account/{{username}}/albums")
		}
		return nil, ierr
	}
	return b, nil
}
