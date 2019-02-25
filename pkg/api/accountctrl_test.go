package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handleListAccountsAuthErr(t *testing.T) {
	accountCtrl := newAccountController()

	r := buildRouter(&appController{accountCtrl: accountCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	resp, err := http.Get(s.URL + "/api/v1/account/")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var jErr jsonErr
	err = json.NewDecoder(resp.Body).Decode(&jErr)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid Authorization Header", jErr.Err)
}

func Test_handleListAccountsOK(t *testing.T) {
	userID := uint(5)
	token, err := signJWT(userID)
	assert.NoError(t, err)

	accountCtrl := newAccountController()

	r := buildRouter(&appController{accountCtrl: accountCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL+"/api/v1/account/", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/JSON")
}
