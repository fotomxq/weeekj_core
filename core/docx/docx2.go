package CoreDocx

import (
	"encoding/json"
	BasePython "github.com/fotomxq/weeekj_core/v5/base/python"
)

// CoreDocx2Type 第二代word处理模块
// 该模块将交给py完成处理操作，并将文件换回临时存储目录
type CoreDocx2Type struct {
}

type ReplaceDataType struct {
	//替换普通文本
	ReplaceText map[string]string `json:"replaceText"`
	//替换图片
	ReplaceImage map[string]ReplaceDataImgType `json:"replaceImage"`
}
type ReplaceDataImgType struct {
	//路径
	Src string `json:"src"`
	//尺寸 毫米
	// 可以通过word生成指定尺寸后设置
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (t *CoreDocx2Type) PutDoc(templateSrc string, replaceData ReplaceDataType) (newFileSrc string, err error) {
	if len(replaceData.ReplaceText) < 1 {
		replaceData.ReplaceText = map[string]string{}
	}
	if len(replaceData.ReplaceImage) < 1 {
		replaceData.ReplaceImage = map[string]ReplaceDataImgType{}
	}
	var paramsJSON []byte
	paramsJSON, err = json.Marshal(replaceData)
	if err != nil {
		return
	}
	newFileSrc, err = BasePython.PushSync("word", 0, "", paramsJSON, templateSrc, 10)
	if err != nil {
		return
	}
	return
}
