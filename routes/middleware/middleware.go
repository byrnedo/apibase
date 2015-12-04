package middleware
import (
	"time"
	"net/http"
	. "github.com/byrnedo/apibase/logger"
)

func LogTime(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(startTime).Seconds()

		Info.Printf(`[%.6f] "%s"`, duration, r.URL.Path)
	})
}