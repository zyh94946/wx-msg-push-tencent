package config

import (
	"encoding/json"
	"errors"
	strip "github.com/grokify/html-strip-tags-go"
	"log"
	"regexp"
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
	Digest	string
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

	if request.Title == "" || request.Content == "" {
		return request, errors.New("Request params fail ")
	}

	regBr, _ := regexp.Compile(`<(?i:br)[\S\s]+?>`)
	request.Digest = regBr.ReplaceAllString(request.Content, "\n")

	request.Digest = strip.StripTags(request.Digest)

	return request, nil
}
