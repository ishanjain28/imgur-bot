package imgur

import (
	"net/url"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ishanjain28/imgur-bot/log"
	"net/http"
)

func (i *Imgur) GenerateAccessToken(refreshToken string) *IError {

	//TODO:Returns, grant_type invalid error, most likely a bug in imgur's API
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	form.Add("client_id", i.Config.ClientID)
	form.Add("client_secret", i.Config.ClientSecret)

	body, err := makeUnauthorisedRequest("POST", "/oauth2/token", form.Encode())

	if err != nil {
		log.Warn.Println("Error in generating token", err.Error())
	}
	//req.Header.Add("Content-Type", " multipart/form-data; boundary=------------------------e83a7963e97655ab")

	b, _ := ioutil.ReadAll(body)

	//TODO:Complete this
	fmt.Println(string(b))

	return nil

}

func (i *Imgur) AccountBase(username, accountID string) (base *AccountBase, ierr *IError) {

	var resp *http.Response

	if username != "" {
		resp, ierr = makeAuthorisedRequest("GET", "/3/account/"+username, i.Config.ClientID, "client_id", "")
	} else {
		resp, ierr = makeAuthorisedRequest("GET", "/3/account?account_id="+accountID, i.Config.ClientID, "client_id", "")
	}
	ab := &AccountBase{}

	// First try unmarshalling response into AccountBase struct, If it fails, then unmarshal it into
	// error and return because there are only two possible responses from imgur
	err := json.NewDecoder(resp.Body).Decode(ab)

	if err != nil {
		ierr := &IError{}
		err = json.NewDecoder(resp.Body).Decode(&ierr)

		if err != nil {
			//	what is wrong with imgur!!!!?????
			log.Warn.Println("error in unmarshalling", err.Error())
		}

		return nil, ierr
	}

	return ab, nil
}

func (i *Imgur) ImageCount(username, accessToken string) (b *Basic, ierr *IError) {

	resp, ierr := makeAuthorisedRequest("GET", "/3/account/"+username+"/images/count", accessToken, "access_token", "")
	if ierr != nil {
		return nil, ierr
	}

	b = &Basic{}

	err := json.NewDecoder(resp.Body).Decode(b)
	if err != nil {
		ierr := &IError{}
		err := json.NewDecoder(resp.Body).Decode(&ierr)
		if err != nil {
			log.Warn.Println("Error in unmarshalling", err.Error())
		}

		return nil, ierr
	}

	return b, nil
}

func (i *Imgur) CommentCount(username, accessToken string) (b *Basic, ierr *IError) {

	resp, ierr := makeAuthorisedRequest("GET", "/3/account/"+username+"/comments/count", accessToken, "access_token", "")
	if ierr != nil {
		fmt.Println(ierr)
		return nil, ierr
	}

	b = &Basic{}

	err := json.NewDecoder(resp.Body).Decode(b)
	if err != nil {
		ierr := &IError{}
		err := json.NewDecoder(resp.Body).Decode(&ierr)
		if err != nil {
			log.Warn.Println("Error in unmarshalling", err.Error())
		}

		return nil, ierr
	}

	return b, nil
}
