package imgur

import (
	"net/http"
	"net/url"
	"github.com/ishanjain28/imgur-bot/log"
	"github.com/ishanjain28/imgur-bot/common"
	"strings"
	"encoding/json"
)

type UploadImage struct {
	Data struct {
		ID          string `json:"id"`
		Title       interface{} `json:"title"`
		Description interface{} `json:"description"`
		Datetime    int `json:"datetime"`
		Type        string `json:"type"`
		Animated    bool `json:"animated"`
		Width       int `json:"width"`
		Height      int `json:"height"`
		Size        int `json:"size"`
		Views       int `json:"views"`
		Bandwidth   int `json:"bandwidth"`
		Vote        interface{} `json:"vote"`
		Favorite    bool `json:"favorite"`
		Nsfw        interface{} `json:"nsfw"`
		Section     interface{} `json:"section"`
		AccountURL  interface{} `json:"account_url"`
		AccountID   int `json:"account_id"`
		IsAd        bool `json:"is_ad"`
		InMostViral bool `json:"in_most_viral"`
		Tags        []interface{} `json:"tags"`
		AdType      int `json:"ad_type"`
		AdURL       string `json:"ad_url"`
		InGallery   bool `json:"in_gallery"`
		Deletehash  string `json:"deletehash"`
		Name        string `json:"name"`
		Link        string `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int `json:"status"`
}

func (i *Imgur) UploadImage(imgLink string, user *common.User) (image *UploadImage, err error) {

	client := &http.Client{}

	form := url.Values{}
	form.Add("image", imgLink)
	form.Add("type", "URL")

	req, err := http.NewRequest("POST", hostaddr+"/3/image", strings.NewReader(form.Encode()))
	if err != nil {
		log.Warn.Println("Error in creating upload image req", err.Error())
		return nil, err
	}

	//Add Authorization header
	req.Header.Add("authorization", "Bearer "+user.AccessToken)
	//req.Header.Add("Content-Type", "multipart/form-data; boundary=------------------------e83a7963e97655ab")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	a := &UploadImage{}

	//TODO:Improve this, By handling error and success responses correctly
	json.NewDecoder(resp.Body).Decode(a)

	return a, nil
}
