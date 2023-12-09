package CoreDocx

import "github.com/lukasjarosch/go-docx"

// CoreDocxType 数据结构
type CoreDocxType struct {
	////文件结构
	//DocFile *docx.ReplaceDocx
	////编辑文档引用
	//DocData *docx.Docx
	DocData *docx.Document
}

// LoadDocx 读取文件
func (t *CoreDocxType) LoadDocx(src string) (err error) {
	t.DocData, err = docx.Open(src)
	if err != nil {
		return
	}
	return
}

// OpenEdit 打开编辑模式
func (t *CoreDocxType) OpenEdit() {
	//t.DocData = t.DocFile.Editable()
}

// QuickReplaceAll 批量替换数据
func (t *CoreDocxType) QuickReplaceAll(replaceData map[string]interface{}) (err error) {
	err = t.DocData.ReplaceAll(replaceData)
	//for k, v := range replaceData {
	//	err = t.DocData.Replace(k, v, -1)
	//	if err != nil {
	//		err = errors.New(fmt.Sprint("replace ", k, " to ", v, ", err: ", err))
	//		return
	//	}
	//}
	return
}

// SaveDocx 保存文件
func (t *CoreDocxType) SaveDocx(src string) (err error) {
	//err = t.DocData.WriteToFile(src)
	err = t.DocData.WriteToFile(src)
	return
}

// Close 关闭文件操作
func (t *CoreDocxType) Close() (err error) {
	//err = t.DocFile.Close()
	return
}
