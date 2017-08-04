package imgur

import (
	"net/http"
	"net/url"
	"github.com/ishanjain28/imgur-bot/log"
	"github.com/ishanjain28/imgur-bot/common"
	"strings"
	"encoding/json"
)

func (i *Imgur) UploadImage(imgLink string, user *common.User) (image *uploadImage, ierr *iError) {

	client := &http.Client{}

	form := url.Values{}
	form.Add("image", imgLink)
	form.Add("type", "URL")

	req, err := http.NewRequest("POST", hostaddr+"/3/image", strings.NewReader(form.Encode()))
	if err != nil {
		log.Warn.Println("error in creating upload image req", err.Error())

		return nil, createError(0, "POST", "error in creating request", "/3/image")
	}

	//Add Authorization header
	req.Header.Add("authorization", "Bearer "+user.AccessToken)
	//req.Header.Add("Content-Type", "multipart/form-data; boundary=------------------------e83a7963e97655ab")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, createError(resp.StatusCode, "POST", err.Error(), "/3/image")
	}

	if resp.StatusCode != 200 {
		return nil, createError(resp.StatusCode, "POST", "", "/3/image")
	}

	a := &uploadImage{}

	// First try unmarshalling response into uploadImage struct, If it fails, then unmarshal it into
	// error and return because there are only two possible responses from imgur
	err = json.NewDecoder(resp.Body).Decode(a)
	if err != nil {

		ierr := &error{}
		err = json.NewDecoder(resp.Body).Decode(&ierr)
		if err != nil {
			//	what is wrong with imgur!!!!?????
			log.Warn.Println("error in unmarshalling", err.Error())
		}

		return nil, ierr
	}

	return a, nil
}
