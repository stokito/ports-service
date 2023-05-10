package api

import (
	"context"
	"encoding/json"
	. "github.com/stokito/ports-service/internal/pkg/db"
	"github.com/stokito/ports-service/internal/pkg/domain"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
)

// ErrorResp is JSON returned on API errors
type ErrorResp struct {
	Code    string
	Message string
}

var apiServer *http.Server

func StartApiServer(listenAddr string, credentials map[string]string) {
	apiServer = createApiServer(listenAddr, credentials)
	err := apiServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("CRIT http error: %s\n", err)
	}
}

func createApiServer(listenAddr string, credentials map[string]string) *http.Server {
	apiServiceMux := http.NewServeMux()
	apiServiceMux.HandleFunc("/api/v1/ports/", handlePortsRequest)
	// profiling
	apiServiceMux.HandleFunc("/debug/pprof/", pprof.Index)
	apiServiceMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	apiServiceMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	apiServiceMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	apiServiceMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	authHandler := NewAuthHandlerWrapper(
		apiServiceMux,
		credentials,
		"Ports Service API",
	)
	recoverHandler := &RecoveryHandlerWrapper{
		Handler: authHandler,
	}
	apiServer := &http.Server{
		Addr:    listenAddr,
		Handler: recoverHandler,
	}
	return apiServer
}

func handlePortsRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetPortRequest(w, r)
	case "POST":
		handlePostPortRequest(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleGetPortRequest(w http.ResponseWriter, r *http.Request) {
	portUnloc := r.URL.Query().Get("unloc")
	port := PortsDbConn.FindPort(context.Background(), portUnloc)
	if port == nil {
		errRes := &ErrorResp{
			Code:    "not_found",
			Message: "The port not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeToJson(w, errRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeToJson(w, port)
}

func handlePostPortRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errRes := &ErrorResp{
			Code:    "io_error",
			Message: "Unable to read a requests",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeToJson(w, errRes)
		return
	}
	port := &domain.Port{}
	err = json.Unmarshal(body, port)
	if err != nil {
		errRes := &ErrorResp{
			Code:    "invalid_json",
			Message: "Error on JSON parsing of the request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeToJson(w, errRes)
		return
	}
	PortsDbConn.UpsertPort(context.Background(), "", port)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	log.Printf("INFO A new port is added: %v\n", port.Unlocs)
}

func writeToJson(w http.ResponseWriter, v any) {
	_ = json.NewEncoder(w).Encode(v)
}

func StopApiServer() {
	log.Printf("INFO Shutdown..,\n")
	_ = apiServer.Shutdown(context.Background())
}
