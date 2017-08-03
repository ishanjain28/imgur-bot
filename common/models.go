package common

type User struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	TUsername    string `json:"t_username"`
	TChatID      string `json:"t_chat_id"`
	Username     string `json:"username"`
}
