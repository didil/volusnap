package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/didil/volusnap/pkg/models"
)

func newAccountController(accountSvc accountSvcer) *accountController {
	return &accountController{accountSvc: accountSvc}
}

type accountController struct {
	accountSvc accountSvcer
}

type listAccountsResp struct {
	Accounts models.AccountSlice `json:"accounts,omitempty"`
}

func (ctrl *accountController) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxKey("userID")).(int)
	accounts, err := ctrl.accountSvc.List(userID)
	if err != nil {
		jsonError(w, fmt.Sprintf("ListAccounts err: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &listAccountsResp{Accounts: accounts})
}

type createAccountReq struct {
	Provider string `json:"provider,omitempty"`
	Name     string `json:"name,omitempty"`
	Token    string `json:"token,omitempty"`
}

type createAccountResp struct {
	ID int `json:"id,omitempty"`
}

func (ctrl *accountController) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxKey("userID")).(int)

	create := &createAccountReq{}

	err := json.NewDecoder(r.Body).Decode(create)
	if err != nil {
		jsonError(w, fmt.Sprintf("JSON err: %v", err), http.StatusInternalServerError)
		return
	}

	accountID, err := ctrl.accountSvc.Create(userID, create.Provider, create.Name, create.Token)
	if err != nil {
		jsonError(w, fmt.Sprintf("CreateAccount err: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &createAccountResp{ID: accountID})
}
