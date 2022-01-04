package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

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

	fmt.Printf("user with the username %s is available\n", user.Username)
	return true, nil
}

func initialRegister() (string, string, error) {
	uriExtension := "/_matrix/client/v3/register"
	uri := Server.BaseURL + uriExtension

	var payload struct {
		InitialDeviceDisplayName string `json:"initial_device_display_name"`
	}
	payload.InitialDeviceDisplayName = DeviceName
	body, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return "", "", errors.New("could not create the GET request: " + err.Error())
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", "", errors.New("could not make a request to the HomeServer: " + err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		var sr SynapseErr
		_ = json.NewDecoder(resp.Body).Decode(&sr)
		return "", "", errors.New(sr.ErrCode + ": " + sr.Error)
	}
	defer resp.Body.Close()

	var data struct {
		Session string `json:"session"`
		Flows   []struct {
			Stages []string `json:"stages"`
		} `json:"flows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", errors.New("could not parse the body: " + err.Error())
	}

	fmt.Println("session: ", data.Session)
	return data.Session, data.Flows[0].Stages[0], nil
}

func (user User) Register() (RegisterResponse, error) {
	uriExtension := "/_matrix/client/v3/register"
	uri := Server.BaseURL + uriExtension

	session, flow, err := initialRegister()
	if err != nil {
		return RegisterResponse{}, errors.New(err.Error())
	}

	type Payload struct {
		InhibitLogin             bool   `json:"inhibit_login"`
		InitialDeviceDisplayName string `json:"initial_device_display_name"`
		Password                 string `json:"password"`
		Username                 string `json:"username"`
	}
	payload := Payload{
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
		return RegisterResponse{}, errors.New(err.Error())
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		var sr SynapseErr
		_ = json.NewDecoder(resp.Body).Decode(&sr)
		return RegisterResponse{}, errors.New(fmt.Sprintf("first attempt: %d", resp.StatusCode) + " " + sr.ErrCode + ": " + sr.Error)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		var data struct {
			Session string `json:"session"`
			Flows   []struct {
				Stages []string `json:"stages"`
			} `json:"flows"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&data)
		payload := Register{
			Auth: RegisterAuthData{
				Session: session,
				Type:    flow,
			},
			InhibitLogin:             false,
			InitialDeviceDisplayName: DeviceName,
			Password:                 user.Password,
			Username:                 user.Username,
		}
		body, _ = json.Marshal(payload)
		request, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(request)
		defer resp.Body.Close()

		if err != nil {
			return RegisterResponse{}, errors.New(err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			var sr SynapseErr
			_ = json.NewDecoder(resp.Body).Decode(&sr)
			return RegisterResponse{}, errors.New(sr.ErrCode + ": " + sr.Error)
		}

		var matrixData RegisterResponse
		if err := json.NewDecoder(resp.Body).Decode(&matrixData); err != nil {
			return RegisterResponse{}, errors.New(err.Error())
		}
		return matrixData, nil
	}

	return RegisterResponse{}, nil
}
