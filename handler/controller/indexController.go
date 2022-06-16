//Package controller
// @Author binwu.zhang 2022/6/15 8:41 PM
// @Description:
package controller

import (
	"aig-tech-okr/handler/logic"
	"github.com/gin-gonic/gin"
)

type IndexController struct {
	BaseController
}

func (i IndexController) Index(c *gin.Context) {
	baseLogic := &logic.BaseLogic{C: c}
	indexLogic := &logic.IndexLogic{BaseLogic: baseLogic}
	indexLogic.Index()
}
