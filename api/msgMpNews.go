package api

import (
	strip "github.com/grokify/html-strip-tags-go"
	"regexp"
)

// 图文消息（mpnews）
func NewMpNews(opts *MsgOpts) *mpNewsMsg {
	// 将内容中的br标签换成换行符，过滤html标签，生成摘要
	regBr, _ := regexp.Compile(`<(?i:br)[\S\s]+?>`)
	digest := regBr.ReplaceAllString(opts.Content, "\n")
	digest = strip.StripTags(digest)

	return &mpNewsMsg{
		msgPublic: msgPublic{
			ToUser:                 opts.ToUser,
			ToParty:                opts.ToParty,
			ToTag:                  opts.ToTag,
			AgentId:                opts.AgentId,
			MsgType:                "mpnews",
			EnableDuplicateCheck:   opts.EnableDuplicateCheck,
			DuplicateCheckInterval: opts.DuplicateCheckInterval,
		},
		MpNews: &mpNewsArticles{Articles: []*mpNewsArticleItem{{
			Title:        opts.Title,
			ThumbMediaId: opts.MediaId,
			Content:      opts.Content,
			Digest:       digest,
		}}},
	}
}

type mpNewsMsg struct { // 是否必须、说明
	msgPublic
	MpNews *mpNewsArticles `json:"mpnews"`         //	是	图文消息，一个图文消息支持1到8条图文
	Safe   int             `json:"safe,omitempty"` //	否	表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，2表示仅限在企业内分享，默认为0；注意仅mpnews类型的消息支持safe值为2，其他消息类型不支持

	isRetry bool
}

type mpNewsArticles struct {
	Articles []*mpNewsArticleItem `json:"articles"`
}

type mpNewsArticleItem struct {
	Title            string `json:"title"`                        //	是	标题，不超过128个字节，超过会自动截断（支持id转译）
	ThumbMediaId     string `json:"thumb_media_id"`               //	是	图文消息缩略图的media_id, 可以通过素材管理接口获得。此处thumb_media_id即上传接口返回的media_id
	Author           string `json:"author,omitempty"`             //	否	图文消息的作者，不超过64个字节
	ContentSourceUrl string `json:"content_source_url,omitempty"` //	否	图文消息点击“阅读原文”之后的页面链接
	Content          string `json:"content"`                      //	是	图文消息的内容，支持html标签，不超过666 K个字节（支持id转译）
	Digest           string `json:"digest,omitempty"`             //	否	图文消息的描述，不超过512个字节，超过会自动截断（支持id转译）
}
