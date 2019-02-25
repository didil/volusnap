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
	Volumes []Volume `json:"volumes,omitempty"`
}

func (ctrl *volumeController) handleListVolumes(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get accountID: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(ctxKey("userID")).(int)
	account, err := ctrl.accountSvc.Get(userID, accountID)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get account: %v", err), http.StatusInternalServerError)
		return
	}

	factory := pRegistry.getProviderServiceFactory(account.Provider)
	if factory == nil {
		jsonError(w, fmt.Sprintf("could not get provider factory for %v", account.Provider), http.StatusInternalServerError)
		return
	}

	providerSvc := factory.Build(account.Token)
	volumes, err := providerSvc.ListVolumes()
	if err != nil {
		jsonError(w, fmt.Sprintf("could not list volumes: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &listVolumesResp{Volumes: volumes})
}
