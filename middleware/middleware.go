package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

type Middleware interface {
	Intercept(next http.HandlerFunc) http.HandlerFunc
}

func NewAuthMiddleware(cookieStore *sessions.CookieStore, logger *log.Logger) *authMiddleware {
	return &authMiddleware{
		logger.WithFields(log.Fields{
			"file": "authMiddleware",
		}),
		cookieStore,
	}
}
