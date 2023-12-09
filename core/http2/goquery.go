package CoreHttp2

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

//goquery是第三方库，实现了类似jQuery的功能，可以用来解析HTML文档

type CoreGoquery struct {
	//核心
	Core *Core
	//文档集合
	Doc *goquery.Document
	//错误信息
	Err error
}

// ResultGoquery 通过goquery获取HTML
// return *goquery.Document , error 文档操作句柄，是否成功
func (t *Core) ResultGoquery() *CoreGoquery {
	//构建集合
	c := &CoreGoquery{
		Core: t,
	}
	//解析数据
	c.Doc, c.Err = goquery.NewDocumentFromReader(t.Response.Body)
	if c.Err != nil {
		return c
	}
	//反馈
	return c
}

// ResultGoqueryFindCount 获取HTML文档字符串长度
func (t *CoreGoquery) ResultGoqueryFindCount(findStr string) int {
	var s string
	s, t.Err = t.Doc.Html()
	if t.Err != nil {
		return -1
	}
	return strings.Count(s, findStr)
}
