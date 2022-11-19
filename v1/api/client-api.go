package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Private struct {
	Api       string `json:"api"`
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
	OrgID     string `json:"org_id"`
}

func (private *Private) PrivateClientApi(command []byte) []byte {
	client := http.Client{}
	sign := hashBySegments(private.SecretKey, command)
	req, errRequest := http.NewRequest("POST", private.Api, bytes.NewBuffer(command))
	if errRequest != nil {
		log.Println(errRequest)
	}
	req.Header.Add("Key", private.ApiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Sign", sign)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	rbody, _ := io.ReadAll(req.Body)
	fmt.Println(rbody)
	body, err := io.ReadAll(resp.Body)
	return body
}

type PublicClientApi struct {
	Command string
}

func (public *PublicClientApi) Return() []byte {
	client := http.Client{}
	req, errRequest := http.NewRequest("GET", public.Command, nil)
	if errRequest != nil {
		log.Println(errRequest)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	return body
}

func hashBySegments(secret string, command []byte) string {

	sig := hmac.New(sha512.New, []byte(secret))
	sig.Write(command)
	return hex.EncodeToString(sig.Sum(nil))

}
