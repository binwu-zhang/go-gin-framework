package handler

import (
	"aig-tech-okr/handler/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.Engine) {

	//g := router.Group("/group")
	//g.Use(middleware.SignVerifyMid(), middleware.AuthenticationMid())

	router.POST("/", new(controller.IndexController).Index)
}
