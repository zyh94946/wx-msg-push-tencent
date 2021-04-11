package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zyh94946/wx-msg-push-tencent/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AppMsg interface {
	sendMsg(msgJson []byte, token string) (err error, msgErrCode int64)
}

func Send(msg AppMsg, at *AccessToken) error {

	token, err := at.GetToken(false)
	if err != nil {
		return err
	}

	msgJson, err := json.Marshal(msg)
	if err != nil {
		return errors.New("msg to json err:" + err.Error())
	}

	err, errCode := msg.sendMsg(msgJson, token)

	switch errCode {
	case -1, 40014, 42001: // retry
		token, err = at.GetToken(true)
		if err != nil {
			return err
		}
		err, _ = msg.sendMsg(msgJson, token)
		return err

	}

	return err
}

type msgPublic struct {
	ToUser                 string `json:"touser,omitempty"`                   //	否	成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为@all，则向关注该企业应用的全部成员发送
	ToParty                string `json:"toparty,omitempty"`                  //	否	部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	ToTag                  string `json:"totag,omitempty"`                    //	否	标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	MsgType                string `json:"msgtype"`                            //	是	消息类型
	AgentId                int64  `json:"agentid"`                            //	是	企业应用的id，整型。企业内部开发，可在应用的设置页面查看；第三方服务商，可通过接口 获取企业授权信息 获取该参数值
	EnableIdTrans          int    `json:"enable_id_trans,omitempty"`          //	否	表示是否开启id转译，0表示否，1表示是，默认0
	EnableDuplicateCheck   int    `json:"enable_duplicate_check,omitempty"`   //	否	表示是否开启重复消息检查，0表示否，1表示是，默认0
	DuplicateCheckInterval int    `json:"duplicate_check_interval,omitempty"` //	否	表示是否重复消息检查的时间间隔，默认1800s，最大不超过4小时
}

func (mn *msgPublic) sendMsg(msgJson []byte, token string) (error, int64) {

	sendUrl := config.SendMsgUrl
	sendUrl = strings.ReplaceAll(sendUrl, "ACCESS_TOKEN", token)

	hpClient := http.Client{
		Timeout: 2 * time.Second,
	}

	//log.Println(string(msgJson))

	bodyType := "application/json;charset=utf-8"
	resp, err := hpClient.Post(sendUrl, bodyType, bytes.NewBuffer(msgJson))
	if err != nil {
		return errors.New("Send msg return error: " + err.Error()), 0
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return errors.New("Send msg return error, http: " + resp.Status), 0
	}

	type sendMsgResp struct {
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	sendMsgReturn := sendMsgResp{}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&sendMsgReturn); err != nil {
		return errors.New("Send msg json decode error, body: " + string(body) + " error: " + err.Error()), 0
	}

	switch sendMsgReturn.ErrCode {
	case 0:
		return nil, 0
	}

	return errors.New("Send msg return errcode, errcode: " + strconv.FormatInt(sendMsgReturn.ErrCode, 10) + " errmsg: " + sendMsgReturn.ErrMsg), sendMsgReturn.ErrCode

}
