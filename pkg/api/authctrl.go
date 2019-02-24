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

type signupReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type signupResp struct {
	ID uint `json:"id,omitempty"`
}

func (ctrl *authController) handleSignup(w http.ResponseWriter, r *http.Request) {
	signup := &signupReq{}

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

	err = json.NewEncoder(w).Encode(&signupResp{ID: id})
	if err != nil {
		jsonError(w, fmt.Sprintf("signup resp json err: %v", err), http.StatusInternalServerError)
		return
	}
}

type loginReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type loginResp struct {
	Token string `json:"token,omitempty"`
}

func (ctrl *authController) handleLogin(w http.ResponseWriter, r *http.Request) {
	login := &loginReq{}

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

	err = json.NewEncoder(w).Encode(&loginResp{Token: token})
	if err != nil {
		jsonError(w, fmt.Sprintf("login resp json err: %v", err), http.StatusInternalServerError)
		return
	}
}
