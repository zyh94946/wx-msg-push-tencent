package main

import (
	"context"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
	"github.com/zyh94946/wx-msg-push-tencent/api"
	"github.com/zyh94946/wx-msg-push-tencent/config"
	"log"
)

var resp = events.APIGatewayResponse{
	IsBase64Encoded: false,
	StatusCode:      200,
	Headers:         map[string]string{"Content-Type": "application/json"},
	Body:            `{"errorCode":0,"errorMessage":""}`,
}

func process(ctx context.Context, event map[string]interface{}) (events.APIGatewayResponse, error) {

	request, err := config.ParseRequest(event)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	at := &api.AccessToken{
		CorpId:     config.CorpId,
		CorpSecret: config.CorpSecret,
	}

	var appMsg api.AppMsg

	switch request.MsgType {
	case config.MsgTypeMpNews:
		appMsg = api.NewMpNews(request.Title, request.Content)

	case config.MsgTypeText:
		appMsg = api.NewText(request.Content)
	}

	err = api.Send(appMsg, at)
	if err != nil {
		log.Println(err)
	}

	return resp, err
}

func main() {
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(process)
}
