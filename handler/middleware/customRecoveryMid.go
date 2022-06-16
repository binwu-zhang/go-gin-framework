package middleware

import "C"
import (
	"aig-tech-okr/handler/cont"
	"aig-tech-okr/handler/util"
	"aig-tech-okr/libs"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
	"runtime/debug"
	"strconv"
)

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {

			if r := recover(); r != nil {

				buf := debug.Stack()
				libs.ErrorLog("CustomRecovery", "", c.GetString("traceid"), fmt.Sprintln(fmt.Sprintf("%s", buf)))

				// 代码位置
				_, file, line, _ := runtime.Caller(4)

				// recover错误，转string
				var err string
				switch v := r.(type) {
				case error:
					err = v.Error()
				default:
					err = r.(string)
				}
				libs.ErrorLog("CustomRecovery", file+":"+strconv.Itoa(line), c.GetString("traceid"), err)
				c.AbortWithStatusJSON(util.ErrResponse(c, c.GetHeader("lang"), cont.RCServerError))
				return
			}
		}()

		c.Next()
	}
}
