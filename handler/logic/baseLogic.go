//Package logic
// @Author binwu.zhang 2022/3/16 3:21 下午
// @Description:
package logic

import (
	"aig-tech-okr/handler/cont"
	"aig-tech-okr/handler/entity"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

type BaseLogic struct {
	C              *gin.Context //
	RecordResponse bool         //是否记录接口返回数据
	HttpCode       int          // http Code  1xx-5xx
	ErrCode        int          //
	Lang           string       //
	Platform       string       //darwin代表Mac，windows_nt代表windows android  ios
	ErrMsg         string       //
	TraceId        string       //
}

func (i *BaseLogic) BindParams(params interface{}) error {

	//默认所有接口记录返回数据
	i.RecordResponse = true

	headerInfo := entity.Header{}
	headerString := i.C.GetString("header")
	_ = json.Unmarshal([]byte(headerString), &headerInfo)

	i.TraceId = headerInfo.TraceId
	i.Lang = headerInfo.Lang
	i.Platform = headerInfo.Platform

	if params == nil {
		return nil
	}

	return i.C.ShouldBindBodyWith(params, binding.JSON)
}

func (i *BaseLogic) Response(responseData interface{}) {
	if i.HttpCode == 0 {
		i.HttpCode = http.StatusOK
	}

	//
	if ok := http.StatusText(i.HttpCode); ok == "" {
		panic("http code error")
	}

	var response entity.Response

	response.Data = responseData
	response.Code = i.ErrCode
	if response.Code == 0 {
		response.Code = cont.RCSuccess
	}

	//
	response.Msg = i.ErrMsg
	if response.Msg == "" {
		response.Msg = cont.GetErrorMessage(response.Code, i.Lang)
	}

	response.TraceId = i.TraceId
	response.Time = time.Now().UnixNano() / 1e6

	i.C.JSON(
		i.HttpCode,
		response,
	)

	responseByte, _ := json.Marshal(response)

	if i.RecordResponse {
		i.C.Set("response", string(responseByte))
	}
}
