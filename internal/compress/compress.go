package compress

import (
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Piktet/azopkov.git/internal/logger"
	"go.uber.org/zap"
)

const (
	compressTypeGzip      = "gzip"
	compressTypeDeflate   = "deflate"
	compressTypeEmpty     = ""
	headerContentEncoding = "Content-Encoding"
	headerContentType     = "Content-Type"
	headerAcceptEncoding  = "Accept-Encoding"
)

type ResponseWriter interface {
	http.ResponseWriter
	Close()
}

type compressResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	writer              io.WriteCloser
}

func (w *compressResponseWriter) Write(b []byte) (int, error) {
	if w.writer != nil {
		return w.writer.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *compressResponseWriter) Close() {
	if w.writer != nil {
		w.writer.Close()
	}
}

func decompress(r *http.Request) *http.Request {

	if r.Body == nil {
		return r
	}
	decompressType := r.Header.Get(headerContentEncoding)

	switch decompressType {
	case compressTypeGzip:
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			logger.Log().Debug("decompress error",
				zap.Error(err),
				zap.String(headerContentEncoding, decompressType))
			return r
		}
		r.Body = gz
	case compressTypeDeflate:
		lz, err := zlib.NewReader(r.Body)
		if err != nil {
			logger.Log().Debug("decompress error",
				zap.Error(err),
				zap.String(headerContentEncoding, decompressType))
			return r
		}
		r.Body = lz
	case compressTypeEmpty:
	default:
		logger.Log().Debug("decompress error",
			zap.Error(errors.New("unsupport decompress type")),
			zap.String(headerContentEncoding, decompressType))
	}

	return r
}

func compress(w http.ResponseWriter, r *http.Request) ResponseWriter {
	cw := &compressResponseWriter{
		ResponseWriter: w,
	}
	x := r.Header.Get(headerContentType)
	if x != "application/json" && x != "text/html" {
		return cw
	}
	for _, contentType := range r.Header.Values(headerAcceptEncoding) {
		var compressType string
		compressLevel := 1
		for _, value := range strings.Split(contentType, ",") {
			value = strings.TrimSpace(value)
			if strings.HasPrefix(value, "q=") {
				fmt.Sscanf(value, "q=%d", compressLevel)
				continue
			}
			if value != "" {
				compressType = value
			}
		}

		switch compressType {
		case compressTypeGzip:

			gz, err := gzip.NewWriterLevel(w, compressLevel)
			if err != nil {
				logger.Log().Debug("compress error",
					zap.Error(err),
					zap.String(headerContentEncoding, compressType))

				continue
			}
			w.Header().Set(headerContentEncoding, compressType)
			cw.writer = gz
			return cw
		case compressTypeDeflate:
			lz, err := zlib.NewWriterLevel(w, compressLevel)
			if err != nil {
				logger.Log().Debug("compress error",
					zap.Error(err),
					zap.String(headerContentEncoding, compressType))

				continue
			}
			w.Header().Set(headerContentEncoding, compressType)
			cw.writer = lz
			return cw
		}
	}

	return cw
}

func WithCompress(h http.HandlerFunc) http.HandlerFunc {

	compressFn := func(w http.ResponseWriter, r *http.Request) {

		cr := decompress(r)
		defer cr.Body.Close()

		cw := compress(w, r)
		defer cw.Close()

		h.ServeHTTP(cw, cr) // внедряем реализацию http.ResponseWriter
	}

	return http.HandlerFunc(compressFn)
}
