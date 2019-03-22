package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/didil/volusnap/pkg/api"
)

// Signup a new user
func (c *Client) Signup(email string, password string) (string, error) {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&api.SignupReq{Email: email, Password: password})

	resp, err := http.Post(c.serverURL+"/api/v1/auth/signup", "application/JSON", &b)
	if err != nil {
		return "", fmt.Errorf("request error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var jErr api.JSONErr
		err = json.NewDecoder(resp.Body).Decode(&jErr)
		if err != nil {
			return "", fmt.Errorf("err json error: %v", err)
		}

		return "", fmt.Errorf("statuscode error: %v - %v", resp.Status, jErr.Err)
	}

	var sResp api.SignupResp
	err = json.NewDecoder(resp.Body).Decode(&sResp)
	if err != nil {
		return "", fmt.Errorf("json error: %v", err)
	}

	return fmt.Sprintf("Signup Successful ID: %v", sResp.ID), nil
}
