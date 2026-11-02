package middlewares

import (
	"github.com/gin-gonic/gin"
	"openlist/internal/model"
)

func BlockWebLogin(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.Next()
		return
	}
	u, ok := user.(*model.User)
	if ok && u.DisableWebLogin {
		c.AbortWithStatusJSON(403, gin.H{"error": "该账号禁止网页登录"})
		return
	}
	c.Next()
}
