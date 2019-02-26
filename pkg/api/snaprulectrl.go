package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/didil/volusnap/pkg/models"

	"github.com/gorilla/mux"
)

func newSnapRuleController(snapRuleSvc snapRuleSvcer, accountSvc accountSvcer) *snapRuleController {
	return &snapRuleController{
		snapRuleSvc: snapRuleSvc,
		accountSvc:  accountSvc,
	}
}

type snapRuleController struct {
	snapRuleSvc snapRuleSvcer
	accountSvc  accountSvcer
}

type listSnapRulesResp struct {
	SnapRules models.SnapRuleSlice `json:"snaprules"`
}

func (ctrl *snapRuleController) handleListSnapRules(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get accountID: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// fetch account by user ID to make sure user is authorized to access it
	userID := r.Context().Value(ctxKey("userID")).(int)
	account, err := ctrl.accountSvc.Get(userID, accountID)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get account: %v", err), http.StatusInternalServerError)
		return
	}

	snapRules, err := ctrl.snapRuleSvc.List(account.ID)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not list snaprules: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &listSnapRulesResp{SnapRules: snapRules})
}

type createSnapRuleReq struct {
	Frequency    int    `json:"frequency,omitempty"`
	VolumeID     string `json:"volume_id,omitempty"`
	VolumeName   string `json:"volume_name,omitempty"`
	VolumeRegion string `json:"volume_region,omitempty"`
}

type createSnapRuleResp struct {
	ID int `json:"id,omitempty"`
}

func (ctrl *snapRuleController) handleCreateSnapRule(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get accountID: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// fetch account by user ID to make sure user is authorized to access it
	userID := r.Context().Value(ctxKey("userID")).(int)
	account, err := ctrl.accountSvc.Get(userID, accountID)
	if err != nil {
		jsonError(w, fmt.Sprintf("could not get account: %v", err), http.StatusInternalServerError)
		return
	}

	create := &createSnapRuleReq{}

	err = json.NewDecoder(r.Body).Decode(create)
	if err != nil {
		jsonError(w, fmt.Sprintf("JSON err: %v", err), http.StatusInternalServerError)
		return
	}

	snapRuleID, err := ctrl.snapRuleSvc.Create(account.ID, create.Frequency, create.VolumeID, create.VolumeName,create.VolumeRegion)
	if err != nil {
		jsonError(w, fmt.Sprintf("CreateAccount err: %v", err), http.StatusInternalServerError)
		return
	}

	jsonOK(w, &createSnapRuleResp{ID: snapRuleID})
}
