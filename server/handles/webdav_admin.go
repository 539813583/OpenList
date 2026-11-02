package handles

import (
	"encoding/json"
	"net/http"
	"openlist/internal/model"
	"openlist/server/pkg/rediscli"
	"strings"
)

// ListWebDavSession 获取所有WebDAV会话
func ListWebDavSession(w http.ResponseWriter, r *http.Request) {
	pattern := "webdav:session:*"
	keys, err := rediscli.Client().Keys(r.Context(), pattern).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var list []model.WebDavSession
	for _, k := range keys {
		val, err := rediscli.Client().Get(r.Context(), k).Result()
		if err != nil {
			continue
		}
		var s model.WebDavSession
		_ = json.Unmarshal([]byte(val), &s)
		list = append(list, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// KickWebDavSession 踢除指定会话
func KickWebDavSession(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")
	if sid == "" {
		http.Error(w, "sid required", http.StatusBadRequest)
		return
	}
	pattern := "webdav:session:*"
	keys, err := rediscli.Client().Keys(r.Context(), pattern).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, k := range keys {
		val, _ := rediscli.Client().Get(r.Context(), k).Result()
		if strings.Contains(val, sid) {
			rediscli.Client().Del(r.Context(), k)
			break
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"msg":"ok"}`))
}
