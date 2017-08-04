package imgur

type Config struct {
	ClientID     string
	ClientSecret string
	UseFreeAPI   bool
	XMashapeKey  string
}

type Imgur struct {
	Config Config
}

type Basic struct {
	Data    int64 `json:"data"`
	Status  int  `json:"status"`
	Success bool `json:"success"`
}

type IError struct {
	Data struct {
		Error   string `json:"error"`
		Request string `json:"request"`
		Method  string `json:"method"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

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

type uploadImage struct {
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
