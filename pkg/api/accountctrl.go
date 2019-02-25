package api

import (
	"fmt"
	"net/http"
)

func newAccountController(accountSvc accountSvcer) *accountController {
	return &accountController{accountSvc: accountSvc}
}

type accountController struct {
	accountSvc accountSvcer
}

type listAccountsResp struct {
	Accounts []Account `json:"accounts,omitempty"`
}

func (ctrl *accountController) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxKey("userID")).(uint)
	accounts, err := ctrl.accountSvc.List(userID)
	if err != nil {
		jsonError(w, fmt.Sprintf("ListAccounts err: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &listAccountsResp{Accounts: accounts})
}