// Package compress contains method for Gzip compressing
package compress

import (
	"compress/gzip"
	"net/http"
)

// Middleware compresses gzip format in middleware
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = gz
			defer gz.Close()
		}

		next.ServeHTTP(w, r)
	})
}
