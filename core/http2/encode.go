package CoreHttp2

import "net/url"

// GetURLEncode URL编码工具
func GetURLEncode(sendURL string) string {
	return url.QueryEscape(sendURL)
}
