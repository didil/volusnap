package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/didil/volusnap/pkg/api"
)

// Login a new user
func (c *Client) Login(email string, password string) (string, error) {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&api.LoginReq{Email: email, Password: password})

	resp, err := http.Post(c.serverURL+"/api/v1/auth/login", "application/JSON", &b)
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

	var lResp api.LoginResp
	err = json.NewDecoder(resp.Body).Decode(&lResp)
	if err != nil {
		return "", fmt.Errorf("json error: %v", err)
	}

	return fmt.Sprintf("Login Successful Token:\n%v", lResp.Token), nil
}
