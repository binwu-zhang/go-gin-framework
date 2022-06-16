package util

import (
	"aig-tech-okr/handler/cont"
	"aig-tech-okr/handler/entity"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"strconv"
)

//UniqueID
//  @date: 2022-03-14 10:43:27
//  @Description: 唯一id
//  @return string
func UniqueID() string {
	return bson.NewObjectId().Hex()
}

//GetGuid
//  @date: 2022-03-14 10:43:10
//  @Description: 唯一ID
//  @return string
func GetGuid() string {
	guid := xid.New()
	return guid.String()
}

//Md5
//  @date: 2022-03-14 10:43:19
//  @Description: md5sum
//  @param str string
//  @return string
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func BuildParams(params map[string]interface{}, sign bool) string {

	vals := url.Values{}
	for k, v := range params {
		if v == "" {
			continue
		}
		val := ""
		switch v.(type) {
		case int:
			val = strconv.Itoa(v.(int))
			break
		case int32:
			val = strconv.Itoa(int(v.(int32)))
			break
		case uint32:
			val = strconv.Itoa(int(v.(uint32)))
			break
		case int64:
			val = strconv.Itoa(int(v.(int64)))
			break
		case uint64:
			val = strconv.Itoa(int(v.(uint64)))
			break
		case float32:
			val = strconv.FormatFloat(float64(v.(float32)), 'f', 8, 64)
			break
		case float64:
			val = strconv.FormatFloat(v.(float64), 'f', 8, 64)
			break
		case []interface{}:
			j, _ := json.Marshal(v)
			val = string(j)
		case string:
			if len((v).(string)) <= 0 {
				continue
			}
			val = (v).(string)
		}
		vals.Add(k, val)
	}

	if sign {
		//sign := Sign(vals.Encode())
		//vals.Add("sign", sign)
	}

	return vals.Encode()
}

//ErrResponse
//  @date: 2022-03-11 18:48:28
//  @Description:
//  @param c *gin.Context
//  @param responseCode int
//  @param args ...interface{}   0:httpCode
//  @return int
//  @return entity.ResponseV2
func ErrResponse(c *gin.Context, lang string, responseCode int, args ...interface{}) (int, entity.Response) {

	responseData := entity.Response{
		Code:    responseCode,
		Msg:     cont.GetErrorMessage(responseCode, lang),
		TraceId: c.GetString("traceid"),
	}

	//http错误吗
	code := http.StatusOK
	if len(args) >= 1 {
		codeTemp, ok := args[0].(int)
		if ok {
			code = codeTemp
		}
	}

	//自定义错误提示
	if len(args) >= 2 {
		msgTemp, ok := args[1].(string)
		if ok {
			responseData.Msg = msgTemp
		}
	}

	responseDataByte, _ := json.Marshal(responseData)
	c.Set("response", string(responseDataByte))

	return code, responseData
}
