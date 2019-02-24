package api

import (
	"net/http"
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
