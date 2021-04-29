package api

// 文本消息
func NewText(opts *MsgOpts) *textMsg {
	return &textMsg{
		msgPublic: msgPublic{
			ToUser:                 opts.ToUser,
			ToParty:                opts.ToParty,
			ToTag:                  opts.ToTag,
			AgentId:                opts.AgentId,
			MsgType:                "text",
			EnableDuplicateCheck:   opts.EnableDuplicateCheck,
			DuplicateCheckInterval: opts.DuplicateCheckInterval,
		},
		Text: textContent{
			Content: opts.Content,
		},
	}
}

type textMsg struct { // 是否必须、说明
	msgPublic
	Text textContent `json:"text"`           //	是
	Safe int         `json:"safe,omitempty"` //	否	表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，2表示仅限在企业内分享，默认为0；注意仅mpnews类型的消息支持safe值为2，其他消息类型不支持
}

type textContent struct {
	Content string `json:"content"` //	是	消息内容，最长不超过2048个字节，超过将截断（支持id转译）
}
