package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/didil/volusnap/pkg/api"
)

// CreateAccount creates a new volusnap account for a cloud provider
func (c *Client) CreateAccount(authToken, provider, name, providerToken string) (string, error) {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&api.CreateAccountReq{
		Provider: provider,
		Name:     name,
		Token:    providerToken,
	})

	req, err := http.NewRequest(http.MethodPost, c.serverURL+"/api/v1/account/", &b)
	if err != nil {
		return "", fmt.Errorf("newrequest error: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
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

	var cAResp api.CreateAccountResp
	err = json.NewDecoder(resp.Body).Decode(&cAResp)
	if err != nil {
		return "", fmt.Errorf("json error: %v", err)
	}

	return fmt.Sprintf("Created Account Successfully, ID: %v\n", cAResp.ID), nil
}
