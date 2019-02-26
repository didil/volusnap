package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func newVolumeController(accountSvc accountSvcer) *volumeController {
	return &volumeController{accountSvc: accountSvc}
}

type volumeController struct {
	accountSvc accountSvcer
}

type listVolumesResp struct {
	Volumes []Volume `json:"volumes"`
}

func (ctrl *volumeController) handleListVolumes(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get accountID: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(ctxKey("userID")).(int)
	account, err := ctrl.accountSvc.GetForUser(userID, accountID)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get account: %v", err), http.StatusInternalServerError)
		return
	}

	providerSvc, err := getProviderService(account)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get provider service: %v", err), http.StatusInternalServerError)
		return
	}

	volumes, err := providerSvc.ListVolumes()
	if err != nil {
		jsonError(w, fmt.Sprintf("could not list volumes: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &listVolumesResp{Volumes: volumes})
}
