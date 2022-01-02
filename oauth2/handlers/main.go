package handlers

import (
	"fmt"
	"log"
	"net/http"
	"oauth2/types"
	"oauth2/util"
)

func Index(w http.ResponseWriter, r *http.Request) {

}

// OAuth2 handles the web application flow by OAuth 2.0 with GitHub OAuth Apps.
// It flows the OAuth spec.
func OAuth2(w http.ResponseWriter, r *http.Request) {
	code, err := util.Code(r)
	if err != nil {
		log.Fatal("could not get the code: ", err)
	}

	token, err := util.AccessToken(code)
	if err != nil || token.Value == "" {
		//	TODO: redirect to error page
	}

	data, err := util.Data(token)
	if err != nil {

	}

	fmt.Println("data", data)
}

func OpenID(w http.ResponseWriter, r *http.Request) {
	code, err := util.Code(r)
	if err != nil {
		log.Fatal("could not get the code: ", err)
	}

	idToken, err := util.IdToken(code)
	if err != nil {
		log.Fatal("could not get the id token: ", err)
	}

	// openid connect, the result is the receiving user for session management
	_, err = util.ParseJWT(idToken)
	if err != nil {
		log.Fatal(err)
	}

	token := types.Token{
		Value: idToken.AccessToken,
		Type:  idToken.TokenType,
		Scope: idToken.Scope,
	}
	data, err := util.Data(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("data", data)
}
