package CoreLanguage

import (
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	"github.com/tidwall/gjson"
	"sync"
)

var (
	//语言包内存缓冲
	languageData []dataLanguage
	//写入语言包数据的锁
	writeLock sync.Mutex
)

// 语言包设定
type dataLanguage struct {
	//语言类型
	Lang string
	//语言数据集
	Data [][]byte
}

// 加载语言包
func loadLanguage(lang string) (result dataLanguage) {
	//检查语言类型
	if !checkLang(lang) {
		return
	}
	//在缓冲中查询语言包数据
	for _, v := range languageData {
		if v.Lang == lang {
			return v
		}
	}
	//锁定机制
	writeLock.Lock()
	defer writeLock.Unlock()
	//语言加载路径
	dirSrc := fmt.Sprint(CoreFile.BaseDir(), CoreFile.Sep, "languages", CoreFile.Sep, lang)
	//加载目录下所有json文件
	fileList, err := CoreFile.GetFileList(dirSrc, []string{"json"}, true)
	if err != nil {
		return
	}
	//构建数据
	result.Lang = lang
	for _, v := range fileList {
		vData, err := CoreFile.LoadFile(v)
		if err != nil {
			continue
		}
		result.Data = append(result.Data, vData)
	}
	//写入集合
	languageData = append(languageData, result)
	//反馈
	return
}

// 加载指定的字段
func loadVal(data *dataLanguage, path string) string {
	if data.Lang == "" {
		return ""
	}
	for _, v := range data.Data {
		val := gjson.GetBytes(v, path).String()
		if val != "" {
			return val
		}
	}
	return ""
}

// 检查支持的语言
func checkLang(lang string) bool {
	switch lang {
	case "zh_cn":
	default:
		return false
	}
	return true
}
