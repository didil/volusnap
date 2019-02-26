package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/didil/volusnap/pkg/models"

	"github.com/gorilla/mux"
)

func newSnapshotController(snapshotSvc snapshotSvcer, accountSvc accountSvcer) *snapshotController {
	return &snapshotController{
		snapshotSvc: snapshotSvc,
		accountSvc:  accountSvc,
	}
}

type snapshotController struct {
	snapshotSvc snapshotSvcer
	accountSvc  accountSvcer
}

type listSnapshotsResp struct {
	Snapshots models.SnapshotSlice `json:"snapshots"`
}

func (ctrl *snapshotController) handleListSnapshots(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get accountID: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// fetch account by user ID to make sure user is authorized to access it
	userID := r.Context().Value(ctxKey("userID")).(int)
	account, err := ctrl.accountSvc.GetForUser(userID, accountID)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get account: %v", err), http.StatusInternalServerError)
		return
	}

	snapshots, err := ctrl.snapshotSvc.List(account.ID)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not list snapshots: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &listSnapshotsResp{Snapshots: snapshots})
}
