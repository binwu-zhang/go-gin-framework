//Package logic
// @Author binwu.zhang 2022/6/15 8:43 PM
// @Description:
package logic

import (
	"aig-tech-okr/handler/cont"
	"aig-tech-okr/libs"
)

type IndexLogic struct {
	*BaseLogic
}

func (i *IndexLogic) Index() {
	logKey := "AdminLogic-RuleList"

	responseData := struct {
	}{}

	//参数
	params := struct {
	}{}
	if err := i.BindParams(&params); err != nil {
		libs.ErrorLog(logKey, "参数验证失败", i.TraceId, err)
		i.ErrCode = cont.RCParamInvalid
		i.Response(responseData)
		return
	}

	i.Response(responseData)

	return
}