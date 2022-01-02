package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	handler "oauth2/handlers"
	"os"
	"os/exec"

	dotenv "github.com/joho/godotenv"
)

func init() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	_ = os.Setenv("state", genHash())
	_ = os.Setenv("nonce", genHash())
}

func genHash() string {
	h := md5.New()
	return hex.EncodeToString(h.Sum(nil))
}

func openExternBrowser(url string) error {
	log.Printf("url: %s", url)
	err := exec.Command("open", url).Start()
	if err != nil {
		return errors.New("error during opening the browser: " + err.Error())
	}
	return nil
}

func oauthAuthorization(scope string) {
	parameters := fmt.Sprintf("?client_id=%s&redirect_url=%s&allow_signup=%s&scope=%s&state=%s",
		os.Getenv("gh_client_id"), os.Getenv("callback_url"), "false", scope, os.Getenv("state"),
	)
	url := os.Getenv("gh_authorize_url") + parameters
	_ = openExternBrowser(url)
}

func openIDAuthorization(scope string) {
	parameters := fmt.Sprintf("?response_type=%s&client_id=%s&redirect_uri=%s&scope=%s&state=%s&nonce=%s",
		"code", os.Getenv("g_client_id"), os.Getenv("op_callback_url"), scope,
		os.Getenv("state"), os.Getenv("nonce"),
	)
	url := os.Getenv("g_authorization_endpoint") + parameters
	_ = openExternBrowser(url)
}

func main() {
	//oauthAuthorization("user)
	openIDAuthorization("openid+profile+email")

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/oauth2/callback", handler.OAuth2)
	http.HandleFunc("/openid/callback", handler.OpenID)

	log.Println("server is listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
