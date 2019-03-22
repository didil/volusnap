package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func newAuthController(authSvc authSvcer) *authController {
	return &authController{authSvc: authSvc}
}

type authController struct {
	authSvc authSvcer
}

// SignupReq signup request format
type SignupReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// SignupResp signup response format
type SignupResp struct {
	ID int `json:"id,omitempty"`
}

func (ctrl *authController) handleSignup(w http.ResponseWriter, r *http.Request) {
	signup := &SignupReq{}

	err := json.NewDecoder(r.Body).Decode(signup)
	if err != nil {
		jsonError(w, fmt.Sprintf("JSON err: %v", err), http.StatusInternalServerError)
		return
	}

	id, err := ctrl.authSvc.Signup(signup.Email, signup.Password)
	if err != nil {
		jsonError(w, fmt.Sprintf("signup err: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &SignupResp{ID: id})
}

// LoginReq login request json
type LoginReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// LoginResp login response json
type LoginResp struct {
	Token string `json:"token,omitempty"`
}

func (ctrl *authController) handleLogin(w http.ResponseWriter, r *http.Request) {
	login := &LoginReq{}

	err := json.NewDecoder(r.Body).Decode(login)
	if err != nil {
		jsonError(w, fmt.Sprintf("JSON err: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := ctrl.authSvc.Login(login.Email, login.Password)
	if err != nil {
		jsonError(w, fmt.Sprintf("login err: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &LoginResp{Token: token})
}
