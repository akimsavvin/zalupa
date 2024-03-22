package main

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Role        string `json:"role"`
	Scope       string `json:"scope"`
	Domain      string `json:"domain"`
	TokenType   string `json:"token_type"`
	Env         string `json:"env"`
	ExpiresIn   int    `json:"expires_in"`
	UserId      int    `json:"userId"`
	Customer    string `json:"customer"`
}

type CreatePolicyRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Filters []struct {
		Column   string   `json:"column"`
		Values   []string `json:"values"`
		Operator string   `json:"operator"`
		Not      bool     `json:"not"`
	} `json:"filters"`
	Users        []int `json:"users"`
	VirtualUsers []int `json:"virtualUsers"`
	Groups       []int `json:"groups"`
}

type Policy struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Filters []struct {
		Column   string   `json:"column"`
		Values   []string `json:"values"`
		Operator string   `json:"operator"`
		Not      bool     `json:"not"`
	} `json:"filters"`
	Users        []int `json:"users"`
	VirtualUsers []int `json:"virtualUsers"`
	Groups       []int `json:"groups"`
}
