package api

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zyh94946/work-wx-msg-push/config"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AccessToken struct {
	CorpId     string
	CorpSecret string
	cacheKey   string
}

var instToken = sync.Map{}

func (at *AccessToken) GetToken(isGetNew bool) (string, error) {

	if isGetNew {
		at.expireToken()
	}

	if tokenVal := at.getTokenCache(); tokenVal != "" {
		return tokenVal, nil
	}

	url := config.GetTokenUrl
	url = strings.ReplaceAll(url, "ID", at.CorpId)
	url = strings.ReplaceAll(url, "SECRET", at.CorpSecret)

	hpClient := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := hpClient.Get(url)
	if err != nil {
		return "", errors.New("Get token return error: " + err.Error())
	}
	if resp.StatusCode != 200 {
		return "", errors.New("Get token return error, http: " + resp.Status)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	type tokenResp struct {
		ErrCode     int64  `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	tokenReturn := tokenResp{}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&tokenReturn); err != nil {
		return "", errors.New("Token json decode error, body: " + string(body) + " error: " + err.Error())
	}

	if tokenReturn.ErrCode != 0 {
		return "", errors.New("Get token return errcode, errcode: " + strconv.FormatInt(tokenReturn.ErrCode, 10) + " errmsg: " + tokenReturn.ErrMsg)
	}

	at.setTokenCache(tokenReturn.AccessToken, tokenReturn.ExpiresIn)

	return tokenReturn.AccessToken, nil

}

func (at *AccessToken) getTokenCache() string {

	if tokenVal, isOk := instToken.Load(at.getTokenKey()); isOk {
		return tokenVal.(string)
	}

	return ""
}

func (at *AccessToken) setTokenCache(tokenVal string, expire int64) {
	instToken.Store(at.getTokenKey(), tokenVal)
	time.AfterFunc(time.Duration(expire)*time.Second, at.expireToken)
}

func (at *AccessToken) expireToken() {
	instToken.Store(at.getTokenKey(), "")
}

func (at *AccessToken) getTokenKey() string {
	if at.cacheKey != "" {
		return at.cacheKey
	}

	h := md5.New()
	_, _ = io.WriteString(h, at.CorpId+"_"+at.CorpSecret)
	at.cacheKey = fmt.Sprintf("%x", h.Sum(nil))
	return at.cacheKey
}
