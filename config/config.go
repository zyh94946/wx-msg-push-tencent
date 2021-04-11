package config

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
)
import "os"

var CorpId = os.Getenv("CORP_ID")
var CorpSecret = os.Getenv("CORP_SECRET")
var AgentId, _ = strconv.ParseInt(os.Getenv("AGENT_ID"), 10, 64)
var MediaId = os.Getenv("MEDIA_ID")
var EnableDuplicateCheck = 1     // 是否开启重复消息检查,0表示否,1表示是
var DuplicateCheckInterval = 300 // 重复消息检查的时间间隔,单位秒,最大不超过4小时

var GetTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=ID&corpsecret=SECRET"
var SendMsgUrl = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN"
var MsgTypeMpNews = "mpnews"
var MsgTypeText = "text"

type ApiRequest struct {
	HttpMethod     string            `json:"httpMethod"`
	QueryString    map[string]string `json:"queryString"`
	Body           string            `json:"body"`
	PathParameters map[string]string `json:"pathParameters"`
}

type SimpleRequest struct {
	Title   string
	Content string
	Digest  string
	MsgType string
}

func ParseRequest(event map[string]interface{}) (*SimpleRequest, error) {
	requestPars := ApiRequest{}
	eventJson, _ := json.Marshal(event)
	err := json.Unmarshal(eventJson, &requestPars)
	log.Println("event:", event)
	log.Println("request:", requestPars, "json error:", err)

	request := &SimpleRequest{}
	if requestPars.PathParameters["SECRET"] != CorpSecret {
		return request, errors.New("Request check fail ")
	}

	request.Title = requestPars.QueryString["title"]
	request.Content = requestPars.QueryString["content"]
	request.MsgType = requestPars.QueryString["type"]

	if requestPars.HttpMethod == "POST" {
		type reqBody struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Type    string `json:"type"`
		}
		reqBodyPars := reqBody{}
		if err := json.Unmarshal([]byte(requestPars.Body), &reqBodyPars); err != nil {
			return request, err
		}
		if reqBodyPars.Title != "" {
			request.Title = reqBodyPars.Title
		}
		if reqBodyPars.Content != "" {
			request.Content = reqBodyPars.Content
		}
		if reqBodyPars.Type != "" {
			request.MsgType = reqBodyPars.Type
		}
	}

	if request.Content == "" {
		return request, errors.New("Request params fail ")
	}

	if request.MsgType == "" {
		request.MsgType = MsgTypeMpNews
	}

	return request, nil
}
