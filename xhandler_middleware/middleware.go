package xhandler_middleware

import (
	"github.com/rs/xhandler"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"time"
	"foundation-go/logging"
)

func LoggingMiddleware(handler xhandler.HandlerC, ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// before request
	t := time.Now()

	handler.ServeHTTPC(ctx, w, r)

	log.WithFields(log.Fields{
		"duration":  time.Since(t).Seconds()*1000,
		"path":      r.URL.Path,
		"method":    r.Method,
		"remote_ip": r.RemoteAddr,
		"app":       logging.LoggingApp,
		"project":   logging.LoggingProject,
		"version":   logging.LoggingAppVersion,
	}).Info("API call")
}
