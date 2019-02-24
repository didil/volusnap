package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

type mockAuthSvc struct {
	mock.Mock
}

func (m *mockAuthSvc) Signup(email string, password string) (uint, error) {
	args := m.Called(email, password)
	return uint(args.Int(0)), args.Error(1)
}

func (m *mockAuthSvc) Login(email string, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func Test_handleSignupWithErr(t *testing.T) {
	email := "email@example.com"
	password := "123456"

	authSvc := new(mockAuthSvc)
	authCtrl := newAuthController(authSvc)

	authSvc.On("Signup", email, password).Return(0, fmt.Errorf("some err"))

	r := buildRouter(&appController{authCtrl: authCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&signupReq{Email: email, Password: password})

	resp, err := http.Post(s.URL+"/api/v1/auth/signup", "application/JSON", &b)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 500)

	var jErr jsonErr
	err = json.NewDecoder(resp.Body).Decode(&jErr)
	assert.NoError(t, err)

	assert.Equal(t, "signup err: some err", jErr.Err)

	authSvc.AssertExpectations(t)
}

func Test_handleSignupOk(t *testing.T) {
	email := "email@example.com"
	password := "123456"

	authSvc := new(mockAuthSvc)
	authCtrl := newAuthController(authSvc)

	authSvc.On("Signup", email, password).Return(1, nil)

	r := buildRouter(&appController{authCtrl: authCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&signupReq{Email: email, Password: password})

	resp, err := http.Post(s.URL+"/api/v1/auth/signup", "application/JSON", &b)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 200)

	var sResp signupResp
	err = json.NewDecoder(resp.Body).Decode(&sResp)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), sResp.ID)

	authSvc.AssertExpectations(t)
}

func Test_handleLoginWithErr(t *testing.T) {
	email := "email@example.com"
	password := "123456"

	authSvc := new(mockAuthSvc)
	authCtrl := newAuthController(authSvc)

	authSvc.On("Login", email, password).Return("", fmt.Errorf("some err"))

	r := buildRouter(&appController{authCtrl: authCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&loginReq{Email: email, Password: password})

	resp, err := http.Post(s.URL+"/api/v1/auth/login", "application/JSON", &b)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 500)

	var jErr jsonErr
	err = json.NewDecoder(resp.Body).Decode(&jErr)
	assert.NoError(t, err)

	assert.Equal(t, "login err: some err", jErr.Err)

	authSvc.AssertExpectations(t)
}

func Test_handleLoginOk(t *testing.T) {
	email := "email@example.com"
	password := "123456"

	authSvc := new(mockAuthSvc)
	authCtrl := newAuthController(authSvc)

	authSvc.On("Login", email, password).Return("my-token", nil)

	r := buildRouter(&appController{authCtrl: authCtrl})
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&loginReq{Email: email, Password: password})

	resp, err := http.Post(s.URL+"/api/v1/auth/login", "application/JSON", &b)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 200)

	var lResp loginResp
	err = json.NewDecoder(resp.Body).Decode(&lResp)
	assert.NoError(t, err)

	assert.Equal(t, "my-token", lResp.Token)

	authSvc.AssertExpectations(t)
}
