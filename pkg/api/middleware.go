package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		logrus.WithFields(logrus.Fields{
			"date":       t.Format(time.RFC3339),
			"method":     r.Method,
			"path":       r.URL.Path,
			"useragent":  r.UserAgent(),
			"remoteaddr": r.RemoteAddr,
		}).Info("API Request")

		next.ServeHTTP(w, r)
	})
}

type ctxKey string

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 || authFields[0] != "Bearer" {
			jsonError(w, "Invalid Authorization Header", http.StatusUnauthorized)
			return
		}

		token := authFields[1]
		userID, err := parseJWT(token)
		if err != nil {
			jsonError(w, fmt.Sprintf("Invalid Token: %v", err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey("userID"), userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
