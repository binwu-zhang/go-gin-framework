package entity

//Multilingual
//  @Description: 多语言
type Multilingual struct {
	Zh string `json:"zh"`
	En string `json:"en"`
}

type Header struct {
	Time      int64  `header:"time"`
	UserToken string `header:"usertoken"`
	TraceId   string `header:"traceid"`
	Sign      string `header:"sign"`
	Lang      string `header:"lang"`
	Platform  string `header:"platform"` //darwin代表Mac，windows_nt代表windows android  ios
}

type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"traceid"`
	Time    int64       `json:"time"`
}
