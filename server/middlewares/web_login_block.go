package middlewares

import (
	"net/http"
	"openlist/internal/model"
)

func BlockWebLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*model.User)
		if ok && user.DisableWebLogin {
			http.Error(w, "该账号禁止网页登录", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
