package main

type HomeServer struct {
	BaseURL string
}

type SynapseErr struct {
	ErrCode string `json:"errcode"`
	Error   string `json:"error"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterAuthData struct {
	Session string `json:"session"`
	Type    string `json:"type"`
}

type registerType struct {
	SSO struct {
		IdentityProviders []struct {
			Id    string
			Name  string
			Brand string
		}
	}
	Token              string
	Password           string
	ApplicationService string
}

var RegisterType = registerType{
	Token:              "m.login.token",
	Password:           "m.login.password",
	ApplicationService: "m.login.application_service",
}

type Register struct {
	Auth                     RegisterAuthData `json:"auth"`
	InhibitLogin             bool             `json:"inhibit_login"`
	InitialDeviceDisplayName string           `json:"initial_device_display_name"`
	Password                 string           `json:"password"`
	Username                 string           `json:"username"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
	DeviceId    string `json:"device_id"`
	UserId      string `json:"user_id"`
}
