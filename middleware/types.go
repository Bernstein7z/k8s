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

type UserIdentifier struct {
	Type string `json:"type"` // m.id.user, m.id.thirdparty, m.id.phone
	User string `json:"user"` // username
}

type LoginRequest struct {
	InitialDeviceDisplayName string         `json:"initial_device_display_name"`
	Password                 string         `json:"password"`
	Type                     string         `json:"type"` // m.login.password, m.login.token
	Identifier               UserIdentifier `json:"identifier"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	DeviceId    string `json:"device_id"`
	UserId      string `json:"user_id"`
}

type Admin struct {
	AccessToken string `json:"access_token"`
	UserId      string `json:"user_id"`
	DeviceId    string `json:"device_id"`
}

type CreationContent struct {
	MFederate bool `json:"m.federate"`
}

type Content struct {
	GuestAccess string `json:"guest_access"`
}

type StateEvent struct {
	Type     string  `json:"type"`
	StateKey string  `json:"state_key"`
	Content  Content `json:"content"`
}

type RoomRequest struct {
	Name            string          `json:"name"`
	Preset          string          `json:"preset"` // private_chat, public_chat, trusted_private_chat
	RoomAliasName   string          `json:"room_alias_name"`
	Topic           string          `json:"topic"`
	Visibility      string          `json:"visibility"` // public, private
	Invite          []string        `json:"invite"`     // users who should be invited
	CreationContent CreationContent `json:"creation_content"`
	InitialState    []StateEvent    `json:"initial_state"`
}

type RoomResponse struct {
	RoomId string `json:"room_id"`
}
