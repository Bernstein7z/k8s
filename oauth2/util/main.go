package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	. "oauth2/types"
	"os"
	"strings"
	"time"
)

// Code gets the authorization code provided by the authorization server from the URL query.
func Code(r *http.Request) (string, error) {
	query := r.URL.Query()
	if query.Has("code") {
		if query.Get("state") != os.Getenv("state") {
			return "", errors.New("state does not match")
		}
	} else {
		return "", errors.New("getCode does not exist in query")
	}

	return query.Get("code"), nil
}

func IdToken(code string) (IDToken, error) {
	if code == "" {
		return IDToken{}, errors.New("empty string is not valid")
	}

	values := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     os.Getenv("g_client_id"),
		"client_secret": os.Getenv("g_client_secret"),
		"redirect_uri":  os.Getenv("op_callback_url"),
		"code":          code,
	}
	data, err := json.Marshal(values)
	if err != nil {
		return IDToken{}, errors.New("json marshal: " + err.Error())
	}

	request, _ := http.NewRequest(http.MethodPost, os.Getenv("g_token_endpoint"), bytes.NewBuffer(data))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return IDToken{}, errors.New("error by response: " + err.Error())
	}
	defer resp.Body.Close()

	var token IDToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return IDToken{}, errors.New("could not parse response body: " + err.Error())
	}

	return token, nil
}

func ParseJWT(idToken IDToken) (Payload, error) {
	payload := strings.Split(idToken.IdToken, ".")[1]

	data, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return Payload{}, errors.New("could not parse the payload: " + err.Error())
	}

	var user Payload
	if err := json.NewDecoder(strings.NewReader(string(data))).Decode(&user); err != nil {
		return Payload{}, errors.New("could not decode the payload: " + err.Error())
	}

	if _, err := verifyIDToken(user); err != nil {
		return Payload{}, errors.New("payload verification failed: " + err.Error())
	}

	return user, nil
}

func verifyIDToken(user Payload) (bool, error) {
	if user.Nonce != os.Getenv("nonce") {
		return false, errors.New("false nonce")
	}
	if user.Iss != "https://accounts.google.com" {
		return false, errors.New("false issuer: " + user.Iss)
	}
	if user.Aud != os.Getenv("g_client_id") {
		return false, errors.New("false audience: " + user.Aud)
	}

	if !user.EmailVerified {
		return false, errors.New("email not verified")
	}

	iat := time.Unix(user.Iat, 0)
	check := iat.Add(time.Second * time.Duration(user.Exp))
	isBefore := time.Now().Before(check)
	if !isBefore {
		return false, errors.New("id token is expired")
	}

	return true, nil
}

// AccessToken first makes with authorization code a post request to authorization server and finally parses the
// token information from the response.
func AccessToken(code string) (Token, error) {
	if code == "" {
		return Token{}, errors.New("empty string is not valid")
	}

	values := map[string]string{
		"client_id":     os.Getenv("gh_client_id"),
		"client_secret": os.Getenv("gh_client_secret"),
		"code":          code,
	}
	data, err := json.Marshal(values)
	if err != nil {
		return Token{}, errors.New("json marshal:" + err.Error())
	}

	request, err := http.NewRequest(http.MethodPost, os.Getenv("gh_token_url"), bytes.NewBuffer(data))
	if err != nil {
		return Token{}, errors.New("could not create a new request: " + err.Error())
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return Token{}, errors.New("could not request the server: " + err.Error())
	}
	defer resp.Body.Close()

	var token Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return Token{}, errors.New("could not parse the json response: " + err.Error())
	}

	return token, nil
}

// Data requests the data within the scope from information server (API) with provided token
func Data(token Token) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, os.Getenv("userinfo_endpoint"), nil)
	if err != nil {
		return map[string]interface{}{}, errors.New("could not create a new request: " + err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type, token.Value))
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return map[string]interface{}{}, errors.New("an error occurred during the request: " + err.Error())
	}
	if resp.StatusCode != 200 {
		return map[string]interface{}{}, errors.New("The status code is not 200. " +
			fmt.Sprintf("response: %v %s", resp.StatusCode, resp.Status))
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return map[string]interface{}{}, errors.New("could not parse the response body: " + err.Error())
	}

	return data, nil
}
