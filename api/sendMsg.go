package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zyh94946/work-wx-msg-push/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 图文消息（mpnews）
type MpNewsMsg struct { // 是否必须、说明
	ToUser                 string          `json:"touser,omitempty"`                   //	否	成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为@all，则向关注该企业应用的全部成员发送
	ToParty                string          `json:"toparty,omitempty"`                  //	否	部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	ToTag                  string          `json:"totag,omitempty"`                    //	否	标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	MsgType                string          `json:"msgtype"`                            //	是	消息类型，此时固定为：mpnews
	AgentId                int64           `json:"agentid"`                            //	是	企业应用的id，整型。企业内部开发，可在应用的设置页面查看；第三方服务商，可通过接口 获取企业授权信息 获取该参数值
	MpNews                 *MpNewsArticles `json:"mpnews"`                             //	是	图文消息，一个图文消息支持1到8条图文
	Safe                   int             `json:"safe,omitempty"`                     //	否	表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，2表示仅限在企业内分享，默认为0；注意仅mpnews类型的消息支持safe值为2，其他消息类型不支持
	EnableIdTrans          int             `json:"enable_id_trans,omitempty"`          //	否	表示是否开启id转译，0表示否，1表示是，默认0
	EnableDuplicateCheck   int             `json:"enable_duplicate_check,omitempty"`   //	否	表示是否开启重复消息检查，0表示否，1表示是，默认0
	DuplicateCheckInterval int             `json:"duplicate_check_interval,omitempty"` //	否	表示是否重复消息检查的时间间隔，默认1800s，最大不超过4小时

	isRetry bool
}

type MpNewsArticles struct {
	Articles []*MpNewsArticleItem `json:"articles"`
}

type MpNewsArticleItem struct {
	Title            string `json:"title"`                        //	是	标题，不超过128个字节，超过会自动截断（支持id转译）
	ThumbMediaId     string `json:"thumb_media_id"`               //	是	图文消息缩略图的media_id, 可以通过素材管理接口获得。此处thumb_media_id即上传接口返回的media_id
	Author           string `json:"author,omitempty"`             //	否	图文消息的作者，不超过64个字节
	ContentSourceUrl string `json:"content_source_url,omitempty"` //	否	图文消息点击“阅读原文”之后的页面链接
	Content          string `json:"content"`                      //	是	图文消息的内容，支持html标签，不超过666 K个字节（支持id转译）
	Digest           string `json:"digest,omitempty"`             //	否	图文消息的描述，不超过512个字节，超过会自动截断（支持id转译）
}

func (mn *MpNewsMsg) Send(at *AccessToken) error {

	token, err := at.GetToken(false)
	if err != nil {
		return err
	}

	sendUrl := config.SendMsgUrl
	sendUrl = strings.ReplaceAll(sendUrl, "ACCESS_TOKEN", token)

	mn.MsgType = "mpnews"

	msgJson, err := json.Marshal(mn)
	if err != nil {
		return errors.New("mpNews to json err:" + err.Error())
	}

	hpClient := http.Client{
		Timeout: 2 * time.Second,
	}

	//log.Println(string(msgJson))

	bodyType := "application/json;charset=utf-8"
	resp, err := hpClient.Post(sendUrl, bodyType, bytes.NewBuffer(msgJson))
	if err != nil {
		return errors.New("Send msg return error: " + err.Error())
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return errors.New("Send msg return error, http: " + resp.Status)
	}

	type sendMsgResp struct {
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	sendMsgReturn := sendMsgResp{}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&sendMsgReturn); err != nil {
		return errors.New("Send msg json decode error, body: " + string(body) + " error: " + err.Error())
	}

	switch sendMsgReturn.ErrCode {
	case 0:
		return nil

	case -1, 40014, 42001: // retry
		if mn.isRetry == false {
			_, _ = at.GetToken(true)
			mn.isRetry = true
			err = mn.Send(at)
			return err
		}

	}

	return errors.New("Send msg return errcode, errcode: " + strconv.FormatInt(sendMsgReturn.ErrCode, 10) + " errmsg: " + sendMsgReturn.ErrMsg)

}
