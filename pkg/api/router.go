package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func buildRouter(app *appController) http.Handler {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/", http.NotFoundHandler())

	// root/api router
	apiR := r.PathPrefix("/api/v1").Subrouter()

	// auth router
	authR := apiR.PathPrefix("/auth").Subrouter()
	authR.HandleFunc("/signup", app.authCtrl.handleSignup).Methods(http.MethodPost)
	authR.HandleFunc("/login", app.authCtrl.handleLogin).Methods(http.MethodPost)

	// acounts router
	accountR := apiR.PathPrefix("/account").Subrouter()
	accountR.Use(authMiddleware)
	accountR.HandleFunc("/", app.accountCtrl.handleListAccounts).Methods(http.MethodGet)
	accountR.HandleFunc("/", app.accountCtrl.handleCreateAccount).Methods(http.MethodPost)

	// volumes router
	volumeR := apiR.PathPrefix("/account/{accountID:[0-9]+}/volume").Subrouter()
	volumeR.Use(authMiddleware)
	volumeR.HandleFunc("/", app.volumeCtrl.handleListVolumes).Methods(http.MethodGet)

	// snaprules router
	snapRuleR := apiR.PathPrefix("/account/{accountID:[0-9]+}/snaprule").Subrouter()
	snapRuleR.Use(authMiddleware)
	snapRuleR.HandleFunc("/", app.snapRuleCtrl.handleListSnapRules).Methods(http.MethodGet)
	snapRuleR.HandleFunc("/", app.snapRuleCtrl.handleCreateSnapRule).Methods(http.MethodPost)

	// snapshots routers
	snapshotR := apiR.PathPrefix("/account/{accountID:[0-9]+}/snapshot").Subrouter()
	snapshotR.Use(authMiddleware)
	snapshotR.HandleFunc("/", app.snapshotCtrl.handleListSnapshots).Methods(http.MethodGet)

	return r
}
