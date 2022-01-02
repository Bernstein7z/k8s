package types

type Token struct {
	Value string `json:"access_token"`
	Type  string `json:"token_type"`
	Scope string `json:"scope"`
}

type IDToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

type Payload struct {
	Iss           string `json:"iss"` // Iss (Issuer)
	Azp           string `json:"azp"` // Azp
	Aud           string `json:"aud"` // Aud (audience) == client_id
	Sub           string `json:"sub"` // Sub (subject) user id
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Nonce         string `json:"nonce"`
	Iat           int64  `json:"iat"` // Iat Issued at in unix time in seconds
	Exp           int64  `json:"exp"` // Exp expiration time in unix time in seconds
}
