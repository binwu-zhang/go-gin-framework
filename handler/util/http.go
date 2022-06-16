package util

import (
	"aig-tech-okr/libs"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

type (
	HttpResp struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
)

//
//SendBinaryRequest
//  @date: 2021-11-17 10:44:15
//  @Description: 发送二进制请求
//  @param method string
//  @param url string
//  @param params []byte
//  @param header map[string]string
//  @return resp []byte
//  @return err error
//
func SendBinaryRequest(method, url string, params []byte, header map[string]string) (resp []byte, err error) {

	//请求客户端
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(params))
	if err != nil {
		return
	}
	if header != nil {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}
	response, err := client.Do(req)
	if err != nil {
		return
	}
	return ioutil.ReadAll(response.Body)
}

//
//SendRequest
//  @date: 2021-11-17 10:44:29
//  @Description: 发送普通请求
//  @param method string
//  @param url string
//  @param postParam string
//  @param header map[string]string
//  @return []byte
//  @return error
//
func SendRequest(method, url, postParam string, header map[string]string) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(postParam))
	if err != nil {
		return nil, err
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if header != nil {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	var start, connect, dns, tlsHandshake time.Time
	var dnsT, tLST, connectT, firstT, totalT time.Duration

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			//fmt.Printf("DNS Done: %v\n", time.Since(dns))
			dnsT = time.Since(dns)
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			//fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
			tLST = time.Since(tlsHandshake)
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			//fmt.Printf("Connect time: %v\n", time.Since(connect))
			connectT = time.Since(connect)
		},

		GotFirstResponseByte: func() {
			//fmt.Printf("Time from start to first byte: %v\n", time.Since(start))
			firstT = time.Since(start)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	res, err := client.Do(req)
	totalT = time.Since(start)

	//fmt.Printf("Total time: %v\n", time.Since(start))

	//fmt.Println(fmt.Sprintf("DNS Done: %v\n  TLS Handshake: %v\n  Connect time: %v\n  Time from start to first byte: %v\n  Total time: %v", dnsT, tLST, connectT, firstT, totalT))

	libs.InfoLog("SendRequest", "调用第三方", url, fmt.Sprintf("DNS Done: %v\n  TLS Handshake: %v\n  Connect time: %v\n  Time from start to first byte: %v\n  Total time: %v", dnsT, tLST, connectT, firstT, totalT))

	if err != nil {
		return nil, err
	}

	//fmt.Println(fmt.Sprintf("%v", res), err)

	if res.StatusCode != 200 {
		return nil, errors.New("error reponse code:" + string(res.StatusCode))
	}

	content, err := ioutil.ReadAll(res.Body)

	return content, err
}
