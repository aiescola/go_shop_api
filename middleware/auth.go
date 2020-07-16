package middleware

import (
	"net/http"
	"shopify/util"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

type Middleware interface {
	Intercept(next http.HandlerFunc) http.HandlerFunc
}

type authMiddleware struct {
	logger      *log.Entry
	cookieStore *sessions.CookieStore
}

func NewAuthMiddleware(cookieStore *sessions.CookieStore, logger *log.Logger) *authMiddleware {
	return &authMiddleware{
		logger.WithFields(log.Fields{
			"file": "authMiddleware",
		}),
		cookieStore,
	}
}

func (a *authMiddleware) Intercept(next http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		session, err := a.cookieStore.Get(request, "session")
		if err != nil {
			a.logger.Error("Error retrieving session")
			util.EncodeError(response, http.StatusInternalServerError, "Server error")
			return
		}

		a.logger.Info("session: ", *session)

		untyped := session.Values["username"]
		username, ok := untyped.(string)

		if !ok {
			a.logger.Error("User not authenticated")
			util.EncodeError(response, http.StatusNetworkAuthenticationRequired, "Authentication required")
			return
		}

		a.logger.Info("User authenticated: ", username)
		next.ServeHTTP(response, request)
	}
}
