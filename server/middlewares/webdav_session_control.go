package middlewares

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"openlist/internal/model"
	"openlist/server/pkg/rediscli"
	"time"
)

type ctxKey string
const davCtxKey ctxKey = "dav_ctrl"

type DavCtrl struct {
	rc *http.ResponseController
}

func WebDavSessionControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*model.User)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		rc := http.NewResponseController(w)

		clientIP := r.Header.Get("X-Real-IP")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}
		outerIP := r.Header.Get("X-Forwarded-For")
		if outerIP == "" {
			outerIP = clientIP
		}

		countKeys, err := rediscli.Client().Keys(r.Context(), "webdav:session:"+user.Username+":*").Result()
		if err == nil && user.WebDavMaxOnline > 0 && len(countKeys) >= user.WebDavMaxOnline {
			http.Error(w, "WebDAV在线设备已达上限", http.StatusTooManyRequests)
			return
		}

		b := make([]byte, 16)
		rand.Read(b)
		sid := hex.EncodeToString(b)

		sess := model.WebDavSession{
			Username:    user.Username,
			ClientIP:    clientIP,
			OuterIP:     outerIP,
			UserAgent:   r.UserAgent(),
			ConnectTime: time.Now(),
			SessionID:   sid,
		}
		sessJson, _ := json.Marshal(sess)
		key := "webdav:session:" + user.Username + ":" + clientIP
		rediscli.Client().Set(r.Context(), key, sessJson, 30*time.Minute)

		ctx := context.WithValue(r.Context(), davCtxKey, &DavCtrl{rc: rc})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
