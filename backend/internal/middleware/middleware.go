package middleware

import (
	"backend-lastfm/internal/auth"
	"context"
	"net/http"
)

func Cors(next http.Handler) http.Handler {

	allowed := map[string]bool{
		"http://127.0.0.1:5173": true,
		"http://localhost:5173": true,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateCookie(authService auth.AuthService) func(next http.Handler) http.Handler {
	ctx := context.Background()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookieVal, err := r.Cookie("sid")
			if err != nil {
				http.Error(w, "Unauthorized - Missing or invalid session cookie", http.StatusUnauthorized)
				return
			}
			valid, err := authService.CheckCookieValidity(ctx, cookieVal.Value)
			if !valid || err != nil {
				http.Error(w, "Unauthorized - Invalid session", http.StatusUnauthorized)
				return
			} else {
				next.ServeHTTP(w, r)
				return
			}
		})
	}
}
