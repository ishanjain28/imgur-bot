package imgur

import (
	"net/url"
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/ishanjain28/imgur-bot/log"
)

func (i *Imgur) GenerateAccessToken(refreshToken string) *error {

	//TODO:Returns, grant_type invalid error, most likely a bug in imgur's API
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	form.Add("client_id", i.Config.ClientID)
	form.Add("client_secret", i.Config.ClientSecret)

	client := &http.Client{}

	req, err := http.NewRequest("POST", hostaddr+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return createError(0, "POST", err.Error(), "/oauth2/token")
	}
	req.Header.Add("Content-Type", " multipart/form-data; boundary=------------------------e83a7963e97655ab")

	resp, err := client.Do(req)
	if err != nil {
		return createError(resp.StatusCode, "POST", err.Error(), "/oauth2/token")
	}

	if resp.StatusCode != 200 {
		return createError(resp.StatusCode, "POST", "", "/oauth2/token")
	}

	b, _ := ioutil.ReadAll(resp.Body)

	//TODO:Complete this
	fmt.Println(string(b))

	return nil

}

func (i *Imgur) AccountBase(username, accountID string) (base *accountBase, ierr *iError) {
	var req *http.Request
	var err error
	//Create ab new Request
	if username == "" && accountID != "" {
		req, err = http.NewRequest("GET", hostaddr+"/3/account?account_id="+accountID, nil)
	} else {
		req, err = http.NewRequest("GET", hostaddr+"/3/account/"+username, nil)
	}
	if err != nil {
		return nil, createError(0, "GET", err.Error(), "/3/account/{{username}}")
	}

	//Add Authorization header
	req.Header.Add("authorization", "Client-ID "+i.Config.ClientID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, createError(resp.StatusCode, "GET", err.Error(), "/3/account/{{username}}")
	}
	if resp.StatusCode != 200 {
		return nil, createError(resp.StatusCode, "GET", "", "/3/account/{{username}}")
	}

	ab := &accountBase{}

	// First try unmarshalling response into accountBase struct, If it fails, then unmarshal it into
	// error and return because there are only two possible responses from imgur
	err = json.NewDecoder(resp.Body).Decode(ab)
	if err != nil {

		ierr := &iError{}
		err = json.NewDecoder(resp.Body).Decode(&ierr)
		if err != nil {
			//	what is wrong with imgur!!!!?????
			log.Warn.Println("error in unmarshalling", err.Error())
		}

		return nil, ierr
	}

	return ab, nil
}

func ImageCount(username string) (ierr *error) {

	req, err := http.NewRequest("GET", hostaddr+"/3/account/"+username+"/images/count", nil)

}

func CommentCount(username string) {
	req, err := http.NewRequest("GET", hostaddr+"/3/account/"+username+"/comments/count", nil)
	if err != nil {

	}

}
