package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountSvc struct {
	mock.Mock
}

func (m *mockAccountSvc) List(userID uint) ([]Account, error) {
	args := m.Called(userID)
	return args.Get(0).([]Account), args.Error(1)
}

func Test_handleListAccountsAuthErr(t *testing.T) {
	accountSvc := new(mockAccountSvc)
	accountCtrl := newAccountController(accountSvc)

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

	accountSvc.AssertExpectations(t)
}

func Test_handleListAccountsOK(t *testing.T) {
	userID := uint(5)
	token, err := signJWT(userID)
	assert.NoError(t, err)

	accountSvc := new(mockAccountSvc)
	accountCtrl := newAccountController(accountSvc)

	accounts := []Account{
		Account{Provider: "DigitalOcean", Name: "DO 1"},
	}

	accountSvc.On("List", userID).Return(accounts, nil)

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

	assert.Equal(t, resp.Header.Get("Content-Type"), "application/JSON")

	var lAResp listAccountsResp
	err = json.NewDecoder(resp.Body).Decode(&lAResp)
	assert.NoError(t, err)

	assert.ElementsMatch(t, lAResp.Accounts, accounts)

	accountSvc.AssertExpectations(t)
}
