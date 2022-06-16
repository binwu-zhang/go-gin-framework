package middleware

import (
	"aig-tech-okr/handler/cont"
	"aig-tech-okr/handler/util"
	"github.com/gin-gonic/gin"
)

func AuthenticationMid() gin.HandlerFunc {
	return func(c *gin.Context) {

		//platform := c.GetHeader("platform")
		lang := c.GetHeader("lang")

		// 获取token
		token := c.GetHeader("usertoken")
		if token == "" {
			c.AbortWithStatusJSON(util.ErrResponse(c, lang, cont.RCUserTokenEmpty))
			return
		}

		//判断登录状态

		c.Next()
	}
}
