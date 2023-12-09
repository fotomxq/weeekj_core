package BaseFileSys2

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseFileSys "gitee.com/weeekj/weeekj_core/v5/base/filesys"
	BaseQiniu "gitee.com/weeekj/weeekj_core/v5/base/qiniu"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
)

// GetPublicURLByClaimID 获取文件URL地址
func GetPublicURLByClaimID(fileID int64) (url string) {
	if fileID < 1 {
		return
	}
	dataList := GetPublicURLsByClaimIDs([]int64{fileID})
	if len(dataList) < 1 {
		return
	}
	return dataList[0]
}

// GetFileTypeByClaimID 获取文件格式
func GetFileTypeByClaimID(claimID int64) (fileType string) {
	claimData, _ := BaseFileSys.GetFileClaimByID(&BaseFileSys.ArgsGetFileClaimByID{
		ClaimID: claimID,
		UserID:  -1,
		OrgID:   -1,
	})
	if claimData.ID < 1 {
		return
	}
	fileData, _ := BaseFileSys.GetFileByID(&BaseFileSys.ArgsGetFileByID{
		ID:         claimData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if fileData.ID < 1 {
		return
	}
	fileType = fileData.FileType
	return
}

// GetPublicURLsByClaimIDs 批量获取一组文件
func GetPublicURLsByClaimIDs(fileIDs []int64) (fileURLs []string) {
	//如果没有文件则反馈空
	if len(fileIDs) < 1 {
		return []string{}
	}
	//TODO：临时过度方法，获取旧的文件系统文件引用关系列
	claimList, fileList, _, _ := BaseFileSys.GetFileByClaimIDs(&BaseFileSys.ArgsGetFileByClaimIDs{
		ClaimIDList: fileIDs,
		UserID:      -1,
		OrgID:       -1,
		IsPublic:    false,
	})
	//遍历原始文件
	for _, vClaim := range claimList {
		var vFileData BaseFileSys.FieldsFileType
		for _, vFile := range fileList {
			if vClaim.FileID == vFile.ID {
				vFileData = vFile
				break
			}
		}
		if vFileData.ID < 1 {
			continue
		}
		switch vFileData.FromInfo.System {
		case "local":
			appAPI := BaseConfig.GetDataStringNoErr("AppAPI")
			fileURLs = append(fileURLs, fmt.Sprint(appAPI, "/v4/base/file/public/view/", vClaim.ID))
		case "qiniu":
			//获取七牛云URL
			fileURLs = append(fileURLs, BaseQiniu.GetPublicURLStrNoErr(vClaim.ID))
		}
	}
	return fileURLs
}

// GetPublicURLMapsByClaimIDsTo 获取一组文件URL Map结构体
func GetPublicURLMapsByClaimIDsTo(fileIDs []int64) map[int64]string {
	fileURLs, _ := BaseQiniu.GetPublicURLsMapByFile(&BaseQiniu.ArgsGetPublicURLsByFile{
		FileIDs:  fileIDs,
		UserID:   -1,
		OrgID:    -1,
		IsPublic: true,
	})
	return fileURLs
}

// GetPublicURLFirstByList 从一组图片中抽取第一张图反馈
func GetPublicURLFirstByList(fileIDs []int64) string {
	if len(fileIDs) < 1 {
		return ""
	}
	return GetPublicURLByClaimID(fileIDs[0])
}
