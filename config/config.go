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
var MsgTypeMarkdown = "markdown"

type ApiRequest struct {
	HttpMethod     string            `json:"httpMethod"`
	QueryString    map[string]string `json:"queryString"`
	Body           string            `json:"body"`
	PathParameters map[string]string `json:"pathParameters"`
}

type SimpleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	MsgType string `json:"type"`
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
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
	request.ToUser = requestPars.QueryString["touser"]
	request.ToParty = requestPars.QueryString["toparty"]
	request.ToTag = requestPars.QueryString["totag"]

	if requestPars.HttpMethod == "POST" {
		if err := json.Unmarshal([]byte(requestPars.Body), &request); err != nil {
			return request, err
		}
	}

	if request.Content == "" {
		return request, errors.New("Request params fail ")
	}

	if request.MsgType == "" {
		request.MsgType = MsgTypeMpNews
	}

	if request.ToUser == "" && request.ToParty == "" && request.ToTag == "" {
		request.ToUser = "@all"
	}

	return request, nil
}
