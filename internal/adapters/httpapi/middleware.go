package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mikail-tommard/task-flow/internal/adapters/token"
)

type ctxKey string

const (
	requestIDKey ctxKey = "request_id"
	ctxUserID ctxKey = "user_id"
)

type middleware func(http.Handler) http.Handler

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func Chain(h http.Handler, m ...middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = newRequestID()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	v := ctx.Value(requestIDKey)
	s, _ := v.(string)
	return s
}

func newRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytes += n
	return n, err
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)

		rid := GetRequestID(r.Context())
		dur := time.Since(start)

		log.Printf("rid=%s method=%s path=%s status=%d bytes=%d dur=%s",
			rid, r.Method, r.URL.Path, sw.status, sw.bytes, dur,
		)
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				rid := GetRequestID(r.Context())

				log.Printf("rid=%s panic=%v", rid, rec)

				writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Auth(jwt *token.Service) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if !strings.HasPrefix(h, prefix) {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return 
			}

			raw := strings.TrimSpace(strings.TrimPrefix(h, prefix))
			claims, err := jwt.ParseToken(raw)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return 
			}

			ctx := context.WithValue(r.Context(), ctxUserID, claims.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}