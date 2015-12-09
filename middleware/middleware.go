package middleware
import (
	"time"
	"net/http"
	. "github.com/byrnedo/apibase/logger"
	"encoding/json"
	"github.com/byrnedo/svccommon/msgspec/web"
	"runtime/debug"
)

func LogTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		Info.Printf("[%s] %q %v \n", r.Method, r.URL.Path, duration)
	})
}

func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Error.Println("panic:",err, "\n\n", string(debug.Stack()))
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
