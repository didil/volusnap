package api

import (
	"net/http"
)

func newAccountController() *accountController {
	return &accountController{}
}

type accountController struct {
}

func (ctrl *accountController) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, &struct{}{})
}
