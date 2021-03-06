package middleware

import (
	"encoding/json"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/svccommon/msgspec/web"
	"net/http"
	"runtime/debug"
	"time"
)

type statusLoggingResponseWriter struct {
	status int
	http.ResponseWriter
}

func (w *statusLoggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func LogTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapW := statusLoggingResponseWriter{-1, w}
		next.ServeHTTP(&wrapW, r)
		duration := time.Since(startTime)

		var ips string
		if forIps := r.Header.Get("x-forwarded-for"); len(forIps) > 0 {
			ips = forIps
		} else {
			ips = r.RemoteAddr
		}
		Info.Printf("%s -> [%s] %d %q %v\n", ips, r.Method, wrapW.status, r.URL.Path, duration)
	})
}

func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Error.Println("panic:", err, "\n\n", string(debug.Stack()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func AcceptJsonHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// We send a JSON-API error if the Accept header does not have a valid value.
		if r.Header.Get("Accept") != "application/vnd.api+json" {
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(406)
			json.NewEncoder(w).Encode(web.NewErrorResponse().AddError(406, nil, "", "Accept header must be set to 'application/vnd.api+json'."))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//ForceType application/octet-stream
//Header set Content-Disposition attachment
func FileDownloadHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// We send a JSON-API error if the Accept header does not have a valid value.
		shouldDownload := r.URL.Query().Get("dl")
		if shouldDownload == "1" || shouldDownload == "true" {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", "attachment")
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
