package api

// markdown消息
func NewMarkdown(opts *MsgOpts) *markdownMsg {
	return &markdownMsg{
		msgPublic: msgPublic{
			ToUser:                 opts.ToUser,
			ToParty:                opts.ToParty,
			ToTag:                  opts.ToTag,
			AgentId:                opts.AgentId,
			MsgType:                "markdown",
			EnableDuplicateCheck:   opts.EnableDuplicateCheck,
			DuplicateCheckInterval: opts.DuplicateCheckInterval,
		},
		Markdown: markdownContent{
			Content: opts.Content,
		},
	}
}

type markdownMsg struct { // 是否必须、说明
	msgPublic
	Markdown markdownContent `json:"markdown"` //	是
}

type markdownContent struct {
	Content string `json:"content"` //	是	markdown内容，最长不超过2048个字节，必须是utf8编码
}
