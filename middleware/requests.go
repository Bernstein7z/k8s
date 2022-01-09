package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// IsAvailable
func (user User) IsAvailable() (bool, error) {
	uriExtension := "/_matrix/client/v3/register/available"
	queryParams := "?username=" + user.Username
	uri := Server.BaseURL + uriExtension + queryParams

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return false, errors.New("could not create the GET request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, errors.New("could not make a request to the HomeServer: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response SynapseErr
		_ = json.NewDecoder(resp.Body).Decode(&response)
		return false, errors.New(response.ErrCode + ": " + response.Error)
	}

	return true, nil
}

// RegisterGET
func RegisterGET() ([2]string, error) {
	uriExtension := "/_matrix/client/v3/register"
	uri := Server.BaseURL + uriExtension

	var payload struct {
		InitialDeviceDisplayName string `json:"initial_device_display_name"`
	}
	payload.InitialDeviceDisplayName = DeviceName
	body, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return [2]string{}, errors.New("could not create the GET request: " + err.Error())
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return [2]string{}, errors.New("could not make a request to the HomeServer: " + err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		var sr SynapseErr
		_ = json.NewDecoder(resp.Body).Decode(&sr)
		return [2]string{}, errors.New(sr.ErrCode + ": " + sr.Error)
	}
	defer resp.Body.Close()

	var data struct {
		Session string `json:"session"`
		Flows   []struct {
			Stages []string `json:"stages"`
		} `json:"flows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [2]string{}, errors.New("could not parse the body: " + err.Error())
	}

	return [2]string{data.Session, data.Flows[0].Stages[0]}, nil
}

// Register
func (user User) Register() (Response, error) {
	uriExtension := "/_matrix/client/v3/register"
	uri := Server.BaseURL + uriExtension

	initial, err := RegisterGET()
	if err != nil {
		return Response{}, errors.New(err.Error())
	}

	payload := Register{
		Auth: RegisterAuthData{
			Session: initial[0],
			Type:    initial[1],
		},
		InhibitLogin:             false,
		InitialDeviceDisplayName: DeviceName,
		Password:                 user.Password,
		Username:                 user.Username,
	}
	body, _ := json.Marshal(payload)

	request, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return Response{}, errors.New(err.Error())
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		var sr SynapseErr
		_ = json.NewDecoder(resp.Body).Decode(&sr)
		return Response{}, errors.New(sr.ErrCode + ": " + sr.Error)
	}
	defer resp.Body.Close()

	var matrixData Response
	if err := json.NewDecoder(resp.Body).Decode(&matrixData); err != nil {
		return Response{}, errors.New(err.Error())
	}

	return matrixData, nil
}

func (user User) Login() (Response, error) {
	uriExtension := "/_matrix/client/v3/login"
	uri := Server.BaseURL + uriExtension

	payload := LoginRequest{
		InitialDeviceDisplayName: DeviceName,
		Password:                 user.Password,
		Type:                     "m.login.password",
		Identifier: UserIdentifier{
			Type: "m.id.user",
			User: user.Username,
		},
	}
	body, _ := json.Marshal(payload)

	request, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return Response{}, errors.New("could not request a login: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		var sr SynapseErr
		_ = json.NewDecoder(resp.Body).Decode(&sr)
		return Response{}, errors.New(sr.ErrCode + ": " + sr.Error)
	}

	var result Response
	_ = json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

func (admin Admin) CreateRoom(name string, topic string, users *[]string) (Room, error) {
	uriExtension := "/_matrix/client/v3/createRoom"
	uri := Server.BaseURL + uriExtension

	usersMap := make(map[string]int8)

	if len(*users) != 0 {
		for _, user := range *users {
			usersMap[user] = 0
		}
		usersMap["@admin:localhost"] = 100
	}

	payload := RoomRequest{
		Name:          name,
		Preset:        "private_chat",
		RoomAliasName: name,
		Topic:         topic,
		Visibility:    "private",
		Invite:        *users,
		CreationContent: CreationContent{
			MFederate: false,
		},
		InitialState: []StateEvent{
			{
				Type:     "m.room.guest_access",
				StateKey: "",
				Content: Content{
					GuestAccess: "can_join",
				},
			},
		},
		PowerLevelContentOverride: PowerLevelEventContent{
			EventsDefault: 50,
			Events: map[string]int8{
				"m.room.name":         50,
				"m.room.power_levels": 100,
			},
			Users: usersMap,
		},
	}
	body, _ := json.Marshal(payload)

	request, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", admin.AccessToken))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return Room{}, errors.New(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		var sr SynapseErr
		_ = json.NewDecoder(resp.Body).Decode(&sr)
		return Room{}, errors.New(sr.ErrCode + ": " + sr.Error)
	}
	defer resp.Body.Close()

	var result Room
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Room{}, errors.New(err.Error())
	}

	return result, nil
}
