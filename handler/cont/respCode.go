package cont

import (
	"fmt"
	"strings"
)

var respCodeMap map[int]string

const (
	RCSuccess                 = 200
	RCInvalidUserNamePassword = 1000 // 账号或者密码错误

	//usertoken 登录验证
	RCUserTokenEmpty   = 1001 // userToken 为空
	RCUserTokenExpired = 1002 // userToken 过期了
	RCUserTokenInvalid = 1003 // userToken 无效

	//sign 签名验证
	RCSignEmpty   = 1030 //签名为空
	RCSignInvalid = 1031 //非法签名，重复使用

	//params 参数验证
	RCParamInvalid             = 1060 // 参数不合法
	RCParamInvalidTraceIDEmpty = 1061 // traceid为空
	RCParamInvalidTimeEmpty    = 1062 // time为空
	RCParamInvalidTimeInvalid  = 1063 // time非法 大于当前时间，或小于十秒前

	RCUserNotExist = 1080 //用户不存在

	//服务相关错误
	RCServerError         = 1100 //
	RCServerErrorDatabase = 1101 //数据库错误

	RCRequestDataNotExist = 1201 //请求数据不存在

	ERROR           = 1
	INVALID_REQUEST = 12 // 不合法的请求

	RCSettingsNotAuthority = 1402 //没有权限
)

func init() {
	respCodeMap = map[int]string{
		RCSuccess:                 "ok:ok",
		ERROR:                     "error:error",
		RCSignInvalid:             "非法签名:error",
		RCParamInvalid:            "参数不合法:error",
		RCInvalidUserNamePassword: "账号或者密码错误:error",
		RCUserTokenEmpty:          "userToken 为空:error",
		RCUserTokenExpired:        "userToken 过期:error",
		RCUserTokenInvalid:        "userToken 无效:error",
		INVALID_REQUEST:           "不合法的请求:error",
		RCServerErrorDatabase:     "数据库错误:error",
		RCServerError:             "服务错误:server error",
		RCRequestDataNotExist:     "请求数据不存在:error",
		RCUserNotExist:            "用户不存在:User doesn't exist",

		//sign 签名验证
		RCSignEmpty: "签名为空:error",
		//params 参数验证
		RCParamInvalidTraceIDEmpty: "traceid为空:error",
		RCParamInvalidTimeEmpty:    " time为空:error",
		RCParamInvalidTimeInvalid:  " time非法 大于当前时间，或小于十秒前:error",

		RCSettingsNotAuthority: "没有权限:error",
	}
}

func GetErrorMessage(errCode int, lang string) string {

	if msg, ok := respCodeMap[errCode]; ok {
		msgSlice := strings.Split(msg, ":")
		if len(msgSlice) == 0 {
			return fmt.Sprintf("error: %d", errCode)
		}

		if len(msgSlice) == 1 {
			return msgSlice[0]
		}

		if lang == LangZh {
			return msgSlice[0]
		} else {
			if msgSlice[1] == "error" {
				return fmt.Sprintf("error: %d", errCode)
			}
			return msgSlice[1]
		}
	}

	return fmt.Sprintf("error: %d", errCode)
}
