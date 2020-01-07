package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kou64yama/takanawa"
)

// RequestID returns the middleware.
func RequestID() takanawa.Middleware {
	return takanawa.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(takanawa.HeaderTakanawaRequestID)
			if len(id) == 0 {
				u, _ := uuid.NewRandom()
				id = u.String()
			}

			w.Header().Set(takanawa.HeaderTakanawaRequestID, id)

			ctx := r.Context()
			ctx = context.WithValue(ctx, takanawa.ContextTakanawaRequestID, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}
