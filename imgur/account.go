package imgur

import (
	"net/url"
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/getlantern/errors"
	"strconv"
	"encoding/json"
)

type AccountBase struct {
	Data struct {
		ID             int         `json:"id"`
		URL            string      `json:"url"`
		Bio            interface{} `json:"bio"`
		Avatar         interface{} `json:"avatar"`
		Reputation     int         `json:"reputation"`
		ReputationName string      `json:"reputation_name"`
		Created        int         `json:"created"`
		ProExpiration  bool        `json:"pro_expiration"`
		UserFollow struct {
			Status bool `json:"status"`
		} `json:"user_follow"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

func (i *Imgur) GenerateAccessToken(refreshToken string) error {
	//TOOD:Returns, grant_type invalid error
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	form.Add("client_id", i.Config.ClientID)
	form.Add("client_secret", i.Config.ClientSecret)

	client := &http.Client{}

	fmt.Println(form.Encode())
	req, err := http.NewRequest("POST", hostaddr+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Status Code from imgur is", resp.StatusCode)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(b))

	return nil

}

func (i *Imgur) AccountBase(username, accountID string) (base *AccountBase, err error) {
	var req *http.Request

	//Create a new Request
	if username == "" && accountID != "" {
		req, err = http.NewRequest("GET", hostaddr+"/3/account?account_id="+accountID, nil)
	}
	req, err = http.NewRequest("GET", hostaddr+"/3/account/"+username, nil)
	if err != nil {
		return nil, err
	}

	//Add Authorization header
	req.Header.Add("authorization", "Client-ID "+i.Config.ClientID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Status code is not 200 (" + strconv.Itoa(resp.StatusCode) + ")")
	}

	a := &AccountBase{}

	err = json.NewDecoder(resp.Body).Decode(a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
