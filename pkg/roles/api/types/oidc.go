package types

type OIDCConfig struct {
	ClientID     string   `json:"clientID"`
	ClientSecret string   `json:"clientSecret"`
	Issuer       string   `json:"issuer"`
	RedirectURL  string   `json:"redirectURL"`
	Scopes       []string `json:"scopes"`
}
