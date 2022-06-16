//Package middleware
// @Author binwu.zhang 2022/3/15 2:57 下午
// @Description:
package middleware

import "github.com/gin-gonic/gin"

// MaxAllowed 限流，同时有多少个请求
func MaxAllowed(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire()       // before request
		defer release() // after request
		c.Next()
	}
}
