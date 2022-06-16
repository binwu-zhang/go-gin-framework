package middleware

import (
	"aig-tech-okr/handler/cont"
	"aig-tech-okr/handler/entity"
	"aig-tech-okr/handler/util"
	"aig-tech-okr/libs"
	"aig-tech-okr/libs/cache"
	"encoding/json"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"time"
)

func SignVerifyMid() gin.HandlerFunc {
	return func(c *gin.Context) {

		logKey := "SignVerifyMid"

		//url := c.Request.RequestURI

		//请求body暂存
		_ = c.ShouldBindBodyWith(nil, binding.JSON)

		//绑定header
		HeaderInfo := entity.Header{}
		if err := c.ShouldBindHeader(&HeaderInfo); err != nil {
			c.AbortWithStatusJSON(util.ErrResponse(c, cont.LangEn, cont.RCParamInvalid))
			return
		}

		headerInfoByte, _ := json.Marshal(HeaderInfo)
		c.Set("header", string(headerInfoByte))
		c.Set("traceid", HeaderInfo.TraceId)
		c.Set("lang", HeaderInfo.Lang)
		c.Set("platform", HeaderInfo.Platform)

		if HeaderInfo.TraceId == "" {
			c.AbortWithStatusJSON(util.ErrResponse(c, HeaderInfo.Lang, cont.RCParamInvalidTraceIDEmpty))
			return
		}

		if HeaderInfo.Sign == "abc" {
			return
		}

		//时间判断
		nowTime := time.Now().UnixNano() / 1e6
		if nowTime < HeaderInfo.Time || HeaderInfo.Time == 0 || HeaderInfo.Time < nowTime-10000 {
			fmt.Println(nowTime < HeaderInfo.Time)
			fmt.Println(HeaderInfo.Time == 0)
			fmt.Println(HeaderInfo.Time < nowTime-10000)
			fmt.Println(fmt.Sprintf("%+v", HeaderInfo))
			fmt.Println(nowTime)
			//非法时间戳
			c.AbortWithStatusJSON(util.ErrResponse(c, HeaderInfo.Lang, cont.RCParamInvalidTimeInvalid))
			return
		}

		// 签名验证
		if HeaderInfo.Sign == "" {
			c.AbortWithStatusJSON(util.ErrResponse(c, HeaderInfo.Lang, cont.RCSignEmpty))
			return
		}

		// 签名重复使用
		{
			signRepeat, err := cache.FreeCache.Get([]byte(HeaderInfo.Sign))
			if err != nil && err != freecache.ErrNotFound {
				libs.ErrorLog(logKey, "get cache failed", HeaderInfo.TraceId, err)
			} else {
				if len(signRepeat) > 0 && signRepeat[0] == 1 {
					c.AbortWithStatusJSON(util.ErrResponse(c, HeaderInfo.Lang, cont.RCSignInvalid))
					return
				}
				if err := cache.FreeCache.Set([]byte(HeaderInfo.Sign), []byte{1}, 10); err != nil {
					libs.ErrorLog(logKey, "set cache failed", HeaderInfo.TraceId, err)
				}
			}
		}

		paramsStr := libs.Conf.App.SignPrefix + strconv.FormatInt(HeaderInfo.Time, 10) + HeaderInfo.UserToken + HeaderInfo.TraceId + libs.Conf.App.SignSuffix
		
		if HeaderInfo.Sign != util.Md5(paramsStr) {
			c.AbortWithStatusJSON(util.ErrResponse(c, HeaderInfo.Lang, cont.INVALID_REQUEST))
			return
		}

		c.Next()
	}

}
