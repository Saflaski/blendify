package auth

import (
	"context"
	"net/http"

	"github.com/golang/glog"
)

func ValidateCookie(h AuthHandler, s AuthService) func(next http.Handler) http.Handler {
	glog.Info("Validating cookie")
	ctx := context.Background()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			glog.Info("Validating cookie")

			cookieVal, err := r.Cookie(h.config.FrontendCookieName)
			if err != nil {
				http.Error(w, "Unauthorized - Missing or invalid session cookie", http.StatusUnauthorized)
				return
			}
			userid, err := s.GetUserByValidSessionID(ctx, cookieVal.Value)
			if userid == "" || err != nil {

				http.Error(w, "Unauthorized - Invalid session", http.StatusUnauthorized)

			} else {
				//Get User ID from SessionID
				// userid, err := s.GetUserByValidSessionID(r.Context(), cookieVal.Value)
				// if err != nil {
				// 	http.Error(w, "Cannot authorize - Internal Server Error", http.StatusInternalServerError)
				// 	glog.Errorf("Could not authorize user due to lack of sid:userid map value in service: %w", err)
				// 	return
				// }
				ctx := context.WithValue(r.Context(), h.UserKey, userid)
				glog.Infof("Validated req from userid: %s", userid)
				next.ServeHTTP(w, r.WithContext(ctx))

			}
		})
	}
}
