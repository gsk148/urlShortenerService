package comperess

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var compressingTypes = []string{
	`application/javascript`,
	`application/json`,
	`text/css`,
	`text/html`,
	`text/plain`,
	`text/xml`,
}

var supportedTypesMap = make(map[string]bool)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func CompressGzip(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Accept-Encoding")
		if !strings.Contains(contentType, "gzip") && isNeedCompress(contentType) {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
}

func isNeedCompress(contentType string) bool {
	for _, data := range compressingTypes {
		supportedTypesMap[data] = true
	}

	return supportedTypesMap[contentType]
}
