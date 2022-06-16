package middleware

import (
	"aig-tech-okr/libs"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		request, _ := c.Get(gin.BodyBytesKey)
		response, _ := c.Get("response")

		libs.ReqLog.WithFields(logrus.Fields{
			"uri":         c.Request.RequestURI,
			"method":      c.Request.Method,
			"client_ip":   c.ClientIP(),
			"status":      c.Writer.Status(),
			"latency":     latency.Seconds(),
			"header":      c.GetString("header"),
			"request":     fmt.Sprintf("%s", request),
			"response":    response,
			"employee_id": c.GetUint("employeeId"),
		}).Info("request_records")
	}
}
