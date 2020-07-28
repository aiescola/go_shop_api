package middleware

import (
	"net/http"
	"shopify/util"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

type authMiddleware struct {
	logger      *log.Entry
	cookieStore *sessions.CookieStore
}

func (a authMiddleware) Intercept(next http.HandlerFunc) http.HandlerFunc {
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
			a.logger.Info("User not authenticated")
			util.EncodeError(response, http.StatusNetworkAuthenticationRequired, "Authentication required")
			return
		}

		a.logger.Info("User authenticated: ", username)
		next.ServeHTTP(response, request)
	}
}
