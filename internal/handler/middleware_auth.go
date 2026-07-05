package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/google/uuid"
)

// prevent collisions with keys from other packages
type contextKey string

const (
	contextKeyUserID contextKey = "user_id"
	contextKeyRole   contextKey = "role"
)

func AuthenticateMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "missing authorization header", fmt.Errorf(""))
				return
			}

			// Expected format: "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "invalid authorization header format", fmt.Errorf(""))
				return
			}

			claims, err := authService.ValidateToken(parts[1])
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid or expired token", fmt.Errorf(""))
				return
			}

			ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, contextKeyRole, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireDM rejects requests from non-DM users. Must run after AuthenticateMiddleware.
func RequireDM(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := RoleFromContext(r.Context())
		if !ok || role != "dm" {
			respondWithError(w, http.StatusForbidden, "dungeon master access required", fmt.Errorf(""))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(contextKeyUserID).(uuid.UUID)
	return id, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(contextKeyRole).(string)
	return role, ok
}
