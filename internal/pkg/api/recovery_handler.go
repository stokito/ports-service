package api

import (
	"bytes"
	"log"
	"net/http"
	"runtime/debug"
)

// RecoveryHandlerWrapper Middleware to catch and log panics
type RecoveryHandlerWrapper struct {
	Handler http.Handler
}

func (h *RecoveryHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// catch panic
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			stackTrace := debug.Stack()
			// linearize stacktrace
			stackTrace = bytes.Replace(stackTrace, []byte("\n"), []byte("|"), -1)
			log.Printf("ERR panic: %s: %s\n", panicErr, stackTrace)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	h.Handler.ServeHTTP(w, r)
}
