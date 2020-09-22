package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netspace/common"
	"netspace/util"
)

// 拦截器
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		if len(username) < 3 || !IsTokenValid(token) {
			c.Abort()
			resp := util.NewRespMsg(common.StatusInvalidToken, "token 无效", nil)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}
}
