package imgur

import (
	"net/url"
	"github.com/ishanjain28/imgur-bot/log"
	"encoding/json"
)

func (i *Imgur) UploadImage(imgLink, accessToken string) (image *Image, ierr *IError) {

	form := url.Values{}
	form.Add("image", imgLink)
	form.Add("type", "URL")

	resp, ierr := makeAuthorisedRequest("POST", "/3/image", accessToken, "access_token", form.Encode())
	if ierr != nil {
		return nil, ierr
	}

	a := &Image{}

	// First try unmarshalling response into Image struct, If it fails, then unmarshal it into
	// error and return because there are only two possible responses from imgur
	err := json.NewDecoder(resp.Body).Decode(a)
	if err != nil {

		ierr := &IError{}
		err = json.NewDecoder(resp.Body).Decode(&ierr)
		if err != nil {
			//	what is wrong with imgur!!!!?????
			log.Warn.Println("error in unmarshalling", err.Error())
			return nil, createError(resp.StatusCode, "POST", err.Error(), "/3/image")
		}

		return nil, ierr
	}

	return a, nil
}
