package middlewares

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"openlist/internal/model"
	"openlist/server/pkg/rediscli"
	"time"
)

func WebDavSessionControl(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.Next()
		return
	}
	u, ok := user.(*model.User)
	if !ok {
		c.Next()
		return
	}

	clientIP := c.GetHeader("X-Real-IP")
	if clientIP == "" {
		clientIP = c.ClientIP()
	}
	outerIP := c.GetHeader("X-Forwarded-For")
	if outerIP == "" {
		outerIP = clientIP
	}

	// 统计当前账号在线会话
	countKeys, err := rediscli.Client().Keys(c, "webdav:session:"+u.Username+":*").Result()
	if err == nil && u.WebDavMaxOnline > 0 && len(countKeys) >= u.WebDavMaxOnline {
		c.AbortWithStatusJSON(429, gin.H{"error": "WebDAV在线设备已达上限"})
		return
	}

	// 生成session id
	b := make([]byte, 16)
	rand.Read(b)
	sid := hex.EncodeToString(b)

	sess := model.WebDavSession{
		Username:    u.Username,
		ClientIP:    clientIP,
		OuterIP:     outerIP,
		UserAgent:   c.GetHeader("User-Agent"),
		ConnectTime: time.Now(),
		SessionID:   sid,
	}
	sessJson, _ := json.Marshal(sess)
	key := "webdav:session:" + u.Username + ":" + clientIP
	rediscli.Client().Set(c, key, sessJson, 30*time.Minute)

	c.Next()
}
