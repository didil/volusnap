package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// StartServer starts the API server
func StartServer(port int) error {
	err := loadConfig("config")
	if err != nil {
		return fmt.Errorf("fatal error config file: %s", err)
	}

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	startSnapRulesChecker(db)

	appCtrl := buildAppController(db)

	r := buildRouter(appCtrl)
	logrus.Infof("Starting server on port %v ...", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func startSnapRulesChecker(db *sql.DB) {
	snapRuleSvc := newSnapRuleService(db)
	snapshotSvc := newSnapshotService(db)
	accountSvc := newAccountService(db)
	shooter := newSnapshotTaker()
	checker := newSnapRulesChecker(snapRuleSvc, snapshotSvc, accountSvc, shooter)
	checker.Start()
}

type jsonErr struct {
	Err string `json:"err,omitempty"`
}

func jsonError(w http.ResponseWriter, errStr string, code int) {
	w.Header().Set("Content-Type", "application/JSON")
	w.WriteHeader(code)
	e := json.NewEncoder(w).Encode(&jsonErr{Err: errStr})
	if e != nil {
		http.Error(w, fmt.Sprintf("json encoding error: %v", e), http.StatusInternalServerError)
	}
}

func jsonOK(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/JSON")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		jsonError(w, fmt.Sprintf("jsonOK err: %v", err), http.StatusInternalServerError)
		return
	}
}
