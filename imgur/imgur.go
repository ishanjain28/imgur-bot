package imgur

import "errors"

var hostaddr = "https://api.imgur.com"

func Init(c Config) (*Imgur, error) {

	i := &Imgur{}
	i.Config = c

	if !c.UseFreeAPI {
		if c.XMashapeKey == "" {
			return nil, errors.New("X-Mashape-Key is not set")
		}
		hostaddr = "https://imgur-apiv3.p.mashape.com"
	}

	return i, nil
}
