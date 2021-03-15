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

var GetTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=ID&corpsecret=SECRET"
var SendMsgUrl = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN"

type ApiRequest struct {
	HttpMethod     string            `json:"httpMethod"`
	QueryString    map[string]string `json:"queryString"`
	Body           string            `json:"body"`
	PathParameters map[string]string `json:"pathParameters"`
}

type SimpleRequest struct {
	Title   string
	Content string
}

func ParseRequest(event map[string]interface{}) (SimpleRequest, error) {
	requestPars := ApiRequest{}
	eventJson, _ := json.Marshal(event)
	err := json.Unmarshal(eventJson, &requestPars)
	log.Println("event:", event)
	log.Println("request:", requestPars, "json error:", err)

	if requestPars.PathParameters["SECRET"] != CorpSecret {
		return SimpleRequest{}, errors.New("request check fail")
	}

	request := SimpleRequest{
		Title:   requestPars.QueryString["title"],
		Content: requestPars.QueryString["content"],
	}

	if requestPars.HttpMethod == "POST" {
		type reqBody struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		reqBodyPars := reqBody{}
		if err := json.Unmarshal([]byte(requestPars.Body), &reqBodyPars); err == nil {
			request.Title = reqBodyPars.Title
			request.Content = reqBodyPars.Content
		}
	}

	return request, nil
}
