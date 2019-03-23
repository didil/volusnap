package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/didil/volusnap/pkg/api"
)

// ListAccounts list a user's accounts
func (c *Client) ListAccounts(authToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, c.serverURL+"/api/v1/account/", nil)
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

	var lAResp api.ListAccountsResp
	err = json.NewDecoder(resp.Body).Decode(&lAResp)
	if err != nil {
		return "", fmt.Errorf("json error: %v", err)
	}

	out := "Accounts List:\nID\tProvider\tName\n"

	for _, acc := range lAResp.Accounts {
		out += fmt.Sprintf("%v\t%v\t%v\t\n", acc.ID, acc.Provider, acc.Name)
	}

	return out, nil
}
