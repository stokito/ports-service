package basicauth

import (
	"crypto/subtle"
	"net/http"
)

// AuthHandlerWrapper Wrapper for a http.Handler that adds basic auth
type AuthHandlerWrapper struct {
	Handler     http.Handler
	Credentials map[string]string
	Realm       string
	wwwAuthHdr  string
}

func NewAuthHandlerWrapper(handler http.Handler, credentials map[string]string, realm string) *AuthHandlerWrapper {
	return &AuthHandlerWrapper{
		Handler:     handler,
		Credentials: credentials,
		Realm:       realm,
		wwwAuthHdr:  `Basic realm="` + realm + `", charset="UTF-8"`,
	}
}

func (bah *AuthHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorized := false
	username, password, ok := r.BasicAuth()
	if ok {
		userPassword := bah.Credentials[username]
		if userPassword != "" && subtle.ConstantTimeCompare([]byte(password), []byte(userPassword)) == 1 {
			authorized = true
		}
	}
	if !authorized {
		w.Header().Set("WWW-Authenticate", bah.wwwAuthHdr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	bah.Handler.ServeHTTP(w, r)
}
