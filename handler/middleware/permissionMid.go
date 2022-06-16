package middleware

import (
	"aig-tech-okr/handler/cont"
	"aig-tech-okr/handler/util"
	"github.com/gin-gonic/gin"
	"strings"
)

func PermissionMid() gin.HandlerFunc {
	return func(c *gin.Context) {

		urlSlice := strings.Split(c.Request.RequestURI, "/")

		if util.InArray(urlSlice[3], []string{"settings"}) {

			//查询权限权限
			var permission bool
			if permission {
				//有权限
				c.Next()
			} else {
				//没有权限
				c.AbortWithStatusJSON(util.ErrResponse(c, "zh", cont.RCSettingsNotAuthority))
				return
			}
		}
	}
}
