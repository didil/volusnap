package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didil/volusnap/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountSvc struct {
	mock.Mock
}

func (m *mockAccountSvc) List(userID int) (models.AccountSlice, error) {
	args := m.Called(userID)
	return args.Get(0).(models.AccountSlice), args.Error(1)
}

func (m *mockAccountSvc) Create(userID int, provider string, name string, token string) (int, error) {
	args := m.Called(userID, provider, name, token)
	return args.Int(0), args.Error(1)
}

func (m *mockAccountSvc) GetForUser(userID, accountID int) (*models.Account, error) {
	args := m.Called(userID, accountID)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *mockAccountSvc) Get(accountID int) (*models.Account, error) {
	args := m.Called(accountID)
	return args.Get(0).(*models.Account), args.Error(1)
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

	var jErr JSONErr
	err = json.NewDecoder(resp.Body).Decode(&jErr)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid Authorization Header", jErr.Err)

	accountSvc.AssertExpectations(t)
}

func Test_handleListAccountsOK(t *testing.T) {
	userID := 5
	token, err := signJWT(userID)
	assert.NoError(t, err)

	accountSvc := new(mockAccountSvc)
	accountCtrl := newAccountController(accountSvc)

	accounts := models.AccountSlice{
		&models.Account{Provider: "DigitalOcean", Name: "DO 1"},
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

	var listResp listAccountsResp
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	assert.NoError(t, err)

	assert.ElementsMatch(t, listResp.Accounts, accounts)

	accountSvc.AssertExpectations(t)
}

func Test_handleCreateAccountOK(t *testing.T) {
	userID := 5
	token, err := signJWT(userID)
	assert.NoError(t, err)

	accountSvc := new(mockAccountSvc)
	accountCtrl := newAccountController(accountSvc)

	accountID := 105

	accountSvc.On("Create", userID, "my-provider", "account-name", "my-token").Return(accountID, nil)

	r := buildRouter(&appController{accountCtrl: accountCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&createAccountReq{Provider: "my-provider", Name: "account-name", Token: "my-token"})

	req, err := http.NewRequest(http.MethodPost, s.URL+"/api/v1/account/", &b)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/JSON")

	var createResp createAccountResp
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	assert.NoError(t, err)

	assert.Equal(t, createResp.ID, accountID)

	accountSvc.AssertExpectations(t)
}
