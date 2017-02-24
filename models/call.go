// call.go
package models

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	XML_HEAD_STR string = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>"
	sid                 = "d9e306a9496363a00295d7184f848100"
	token               = "fd821f20c5d48f72d9c8c062619ecd67"
	appid               = "2b385339bfff4b0183c5551e6c5ef8bc"
)

func Post(url string, body []byte) ([]byte, error) {

	return send("POST", url, body)

}

func send(method string, url string, body []byte) ([]byte, error) {

	reqBody := bytes.NewBuffer(body)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 1))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 1,
		},
	}

	req, err := http.NewRequest(method, url, reqBody)
	if nil != err {
		logs.Error("new request failed.", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		logs.Error("http client failed:", err)
		return nil, errors.New("http request failed!")
	}

	defer resp.Body.Close()

	if 200 != resp.StatusCode {
		logs.Error("post url ", url, "failed.status code:", resp.StatusCode)
		return nil, errors.New("return status code not 200")
	}

	resbody, err := ioutil.ReadAll(resp.Body)
	return resbody, err
}

func makeUperMd5Sig() string {

	/* 账户(sid) + 授权令牌(token) + 时间戳 */
	sign := sid + token + time.Now().Format("20060102150405")
	h := md5.New()
	h.Write([]byte(sign))

	SigParameter := hex.EncodeToString(h.Sum(nil))

	return strings.ToUpper(SigParameter)
}

func CallTel(caller string, callee string) {

	type ivr struct {
		AppId  string `xml:"appId"`
		Caller string `xml:"caller"`
		Called string `xml:"called"`
		Data   string `xml:"data"`
	}

	var p ivr
	p.AppId = appid
	p.Caller = caller
	p.Called = callee
	p.Data = "test"

	//编码 xml格式
	req, err := xml.Marshal(p)
	if nil != err {
		fmt.Println("testGetDtmf xml Marshal failed:", err)
		return
	}

	body := XML_HEAD_STR + string(req)
	resp, err := postRequst("call/outCall", []byte(body))
	if err != nil {
		logs.Error("post failed!", err)
		return
	}

	fmt.Println("get response:", string(resp), err)

}
func postRequst(url string, body []byte) ([]byte, error) {

	uri := "https://api.ucpaas.com/2014-06-30/Accounts/" + sid + `/ipcc/` + url

	uri += "?sig=" + makeUperMd5Sig()

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	fmt.Println(uri)
	//fmt.Println(string(body))

	reqBody := ioutil.NopCloser(strings.NewReader(string(body)))
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", uri, reqBody)

	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8;")
	req.Header.Set("Connection", "close")

	/* Authorization域  使用Base64编码（账户Id + 冒号 + 时间戳）(time.Now().Format("20060102150405"))*/
	auths := sid + ":" + time.Now().Format("20060102150405")

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	b64Auth := b64.EncodeToString([]byte(auths))
	req.Header.Set("Authorization", b64Auth)

	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http client failed:", err)
		return nil, errors.New("http request failed!")
	}

	defer resp.Body.Close()

	resbody, err := ioutil.ReadAll(resp.Body)
	return resbody, err
}
