package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func buildRouter(app *appController) http.Handler {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/", http.NotFoundHandler())

	apiR := r.PathPrefix("/api/v1").Subrouter()
	apiR.HandleFunc("/auth/signup", app.authCtrl.handleSignup).Methods(http.MethodPost)
	apiR.HandleFunc("/auth/login", app.authCtrl.handleLogin).Methods(http.MethodPost)

	accountR := apiR.PathPrefix("/account").Subrouter()
	accountR.Use(authMiddleware)
	accountR.HandleFunc("/", app.accountCtrl.handleListAccounts).Methods(http.MethodGet)
	accountR.HandleFunc("/", app.accountCtrl.handleCreateAccount).Methods(http.MethodPost)

	volumeR := apiR.PathPrefix("/account/{accountID:[0-9]+}/volume").Subrouter()
	volumeR.Use(authMiddleware)
	volumeR.HandleFunc("/", app.volumeCtrl.handleListVolumes).Methods(http.MethodGet)

	snapRuleR := apiR.PathPrefix("/account/{accountID:[0-9]+}/snaprule").Subrouter()
	snapRuleR.Use(authMiddleware)
	snapRuleR.HandleFunc("/", app.snapRuleCtrl.handleListSnapRules).Methods(http.MethodGet)
	snapRuleR.HandleFunc("/", app.snapRuleCtrl.handleCreateSnapRule).Methods(http.MethodPost)

	snapshotR := apiR.PathPrefix("/account/{accountID:[0-9]+}/snapshot").Subrouter()
	snapshotR.Use(authMiddleware)
	snapshotR.HandleFunc("/", app.snapshotCtrl.handleListSnapshots).Methods(http.MethodGet)

	return r
}

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

// JSONErr struct
type JSONErr struct {
	Err string `json:"err,omitempty"`
}

func jsonError(w http.ResponseWriter, errStr string, code int) {
	w.Header().Set("Content-Type", "application/JSON")
	w.WriteHeader(code)
	e := json.NewEncoder(w).Encode(&JSONErr{Err: errStr})
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
