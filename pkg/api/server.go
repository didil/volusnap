package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type appController struct {
	authCtrl *authController
}

func buildRouter(app *appController) http.Handler {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/", http.NotFoundHandler())

	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/auth/signup", app.authCtrl.handleSignup).Methods(http.MethodPost)
	s.HandleFunc("/auth/login", app.authCtrl.handleLogin).Methods(http.MethodPost)

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

	err = autoMigrate(db)
	if err != nil {
		return err
	}

	aSvc := newAuthService(db)
	authCtrl := newAuthController(aSvc)

	r := buildRouter(&appController{authCtrl: authCtrl})

	logrus.Infof("Starting server on port %v ...", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

type jsonErr struct {
	Err string `json:"err,omitempty"`
}

func jsonError(w http.ResponseWriter, errStr string, code int) {
	w.WriteHeader(code)
	e := json.NewEncoder(w).Encode(&jsonErr{Err: errStr})
	if e != nil {
		http.Error(w, fmt.Sprintf("json encoding error: %v", e), http.StatusInternalServerError)
	}
}
