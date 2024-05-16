package Router2Excel

import (
	"fmt"
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	BaseTempFile "github.com/fotomxq/weeekj_core/v5/base/temp_file"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/xuri/excelize/v2"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// ExcelTemplate 快速模板工具模块
type ExcelTemplate struct {
	//模板存储总目录
	rootDir string
	//自定义模板存储子目录
	subDir string
	//临时文件过期时间
	tempFileExpire int
	//输出文件名称
	fileName string
	//文件随机摘要值，用于识别唯一性，构建缓冲机制
	fileHash string
	//Excel句柄
	ExcelData *excelize.File
	//图片后缀
	imgSuffix string
}

func (t *ExcelTemplate) SetRootDir(dir string) {
	t.rootDir = dir
}

func (t *ExcelTemplate) SetSubDir(dir string) {
	t.subDir = dir
}

func (t *ExcelTemplate) SetTempFileExpire(tempFileExpire int) {
	t.tempFileExpire = tempFileExpire
}

func (t *ExcelTemplate) SetFileName(fileName string) {
	t.fileName = fileName
}

func (t *ExcelTemplate) SetFileHash(fileHash string) {
	t.fileHash = fileHash
}

func (t *ExcelTemplate) SetImgSuffix(suffix string) {
	t.imgSuffix = suffix
}

// GetExcelTemplate 获取模版文件
func (t *ExcelTemplate) GetExcelTemplate(c any, logErr string, filename string) (excelData *excelize.File, err error) {
	if t.rootDir == "" {
		t.rootDir = fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep, "data", CoreFile.Sep)
	}
	if t.subDir == "" {
		t.subDir = "default"
	}
	excelData, err = excelize.OpenFile(fmt.Sprint(t.rootDir, t.subDir, CoreFile.Sep, filename))
	if err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", load excel template failed, ", err, "err_excel_template")
		return
	}
	t.ExcelData = excelData
	return
}

// BeforeLoadParamsFile 预先下载文件处理
func (t *ExcelTemplate) BeforeLoadParamsFile(c any, logErr string) (result bool) {
	//获取数据
	newID, hash, b := BaseTempFile.SaveFileBefore(t.fileHash)
	if !b {
		return
	}
	//反馈数据
	Router2Mid.ReportData(c, logErr+", ", nil, "", map[string]interface{}{
		"id":   newID,
		"hash": hash,
	})
	return true
}

// SaveExcelTemplate 第二代保存excel文件
func (t *ExcelTemplate) SaveExcelTemplate(c any, logErr string) error {
	if t.tempFileExpire < 1 {
		t.tempFileExpire = 60
	}
	fileSrc, newID, hash, err := BaseTempFile.SaveFile(t.tempFileExpire, t.fileHash, t.fileName, "", "xlsx")
	if err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", save temp file, ", err, "err_make_file")
		return err
	}
	if err := t.SaveExcelFile(fileSrc); err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", , ", err, "err_make_file")
		return err
	}
	//反馈数据
	Router2Mid.ReportData(c, logErr+", ", nil, "", map[string]interface{}{
		"id":   newID,
		"hash": hash,
	})
	return nil
}

// SaveExcelFile 保存修改结果
func (t *ExcelTemplate) SaveExcelFile(src string) (err error) {
	//保存文件
	err = t.ExcelData.SaveAs(src)
	if err != nil {
		return
	}
	//反馈
	return
}

// SetImgByFileSysClaimID 将图片ID写入对应位置
func (t *ExcelTemplate) SetImgByFileSysClaimID(fileClaimID int64, sheet string, cell string) (err error) {
	fileURL := fmt.Sprint(BaseFileSys2.GetPublicURLByClaimID(fileClaimID), t.imgSuffix)
	if fileURL == "" {
		return
	}
	tempFileDir := fmt.Sprint(CoreFile.BaseDir(), CoreFile.Sep, "temp", CoreFile.Sep, CoreFilter.GetNowTimeCarbon().Format("200601"), CoreFile.Sep, CoreFilter.GetNowTimeCarbon().Format("02"))
	err = CoreFile.CreateFolder(tempFileDir)
	if err != nil {
		return
	}
	tempFileSrc := fmt.Sprint(tempFileDir, CoreFile.Sep, CoreFilter.GetRandStr4(10), ".png")
	err = CoreFile.DownloadByURLToTemp(fileURL, tempFileSrc)
	if err != nil {
		return
	}
	_ = t.ExcelData.AddPicture(sheet, cell, tempFileSrc, &excelize.GraphicOptions{
		AutoFit:     true,
		Positioning: "twoCell",
	})
	return
}
