package BaseQiniu

import (
	"errors"
	"fmt"
	BaseFileSys "gitee.com/weeekj/weeekj_core/v5/base/filesys"
	CoreHttp "gitee.com/weeekj/weeekj_core/v5/core/http"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"github.com/qiniu/api.v7/v7/storage"
)

// ArgsGetPublicURLs 获取一组文件的Public URL参数
type ArgsGetPublicURLs struct {
	//一组ID
	ClaimIDList []int64 `json:"claimIDList"`
	//用户ID
	// 可选，用于检测
	UserID int64 `json:"userID"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `json:"orgID"`
	//是否仅公开数据
	IsPublic bool `json:"isPublic"`
}

type DataGetPublicURLs struct {
	//文件引用ID
	// 如果直接通过实体文件获取URL，则该数据总是为0
	ClaimID int64 `json:"claimID"`
	//文件ID
	FileID int64 `json:"fileID"`
	//访问地址
	URL string `json:"url"`
}

// GetPublicURLs 获取一组文件的Public URL
func GetPublicURLs(args *ArgsGetPublicURLs) (data []DataGetPublicURLs, err error) {
	//获取文件基本信息
	var claimList []BaseFileSys.FieldsFileClaimType
	var fileList []BaseFileSys.FieldsFileType
	claimList, fileList, _, err = BaseFileSys.GetFileByClaimIDs(&BaseFileSys.ArgsGetFileByClaimIDs{
		ClaimIDList: args.ClaimIDList,
		UserID:      args.UserID,
		OrgID:       args.OrgID,
		IsPublic:    args.IsPublic,
	})
	if err != nil {
		return
	}
	//遍历获取文件URL地址
	data, err = getPublicURLsAppend(claimList, fileList)
	return
}

// ArgsGetPublicURLsByFile 通过文件实体直接获取文件参数
type ArgsGetPublicURLsByFile struct {
	//一组ID
	FileIDs []int64 `json:"fileIDs"`
	//用户ID
	// 可选，用于检测
	UserID int64 `json:"userID"`
	//组织ID
	// 可选，用于检测
	OrgID int64 `json:"orgID"`
	//是否仅公开数据
	IsPublic bool `json:"isPublic"`
}

// GetPublicURLsByFile 通过文件实体直接获取文件
func GetPublicURLsByFile(args *ArgsGetPublicURLsByFile) (data []DataGetPublicURLs, err error) {
	//获取文件基本信息
	var fileList []BaseFileSys.FieldsFileType
	fileList, _, err = BaseFileSys.GetFileByIDsAndClaim(&BaseFileSys.ArgsGetFileByIDsAndClaim{
		IDs:        args.FileIDs,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
		UserID:     args.UserID,
		OrgID:      args.OrgID,
		IsPublic:   args.IsPublic,
	})
	if err != nil {
		return
	}
	//遍历获取文件URL地址
	data, err = getPublicURLsAppend([]BaseFileSys.FieldsFileClaimType{}, fileList)
	return
}

// 获取文件数据结构体
func getPublicURLsAppend(claimList []BaseFileSys.FieldsFileClaimType, fileList []BaseFileSys.FieldsFileType) (data []DataGetPublicURLs, err error) {
	//遍历获取文件URL地址
	if len(claimList) > 0 {
		for _, vClaim := range claimList {
			vFileData := BaseFileSys.FieldsFileType{}
			for _, vFile := range fileList {
				if vClaim.FileID == vFile.ID {
					vFileData = vFile
					break
				}
			}
			if vFileData.ID < 1 {
				continue
			}
			bucketName := ""
			for _, v := range vFileData.Infos {
				if v.Mark == "bucket" {
					bucketName = v.Val
					break
				}
			}
			var qiniuBucketURL string
			qiniuBucketURL, err = getBucketURL(bucketName)
			if err != nil {
				return
			}
			//获取访问URL路径
			fileURL := storage.MakePublicURL(qiniuBucketURL, vFileData.FromInfo.Mark)
			//写入数据组
			data = append(data, DataGetPublicURLs{
				ClaimID: vClaim.ID,
				FileID:  vFileData.ID,
				URL:     fileURL,
			})
		}
	} else {
		for _, vFileData := range fileList {
			bucketName := ""
			for _, v := range vFileData.Infos {
				if v.Mark == "bucket" {
					bucketName = v.Val
					break
				}
			}
			var qiniuBucketURL string
			qiniuBucketURL, err = getBucketURL(bucketName)
			if err != nil {
				return
			}
			//获取访问URL路径
			fileURL := storage.MakePublicURL(qiniuBucketURL, vFileData.FromInfo.Mark)
			//写入数据组
			data = append(data, DataGetPublicURLs{
				ClaimID: 0,
				FileID:  vFileData.ID,
				URL:     fileURL,
			})
		}
	}
	return
}

// Deprecated: 准备废弃
// GetPublicURLsMap 获取文件URL地址
// 注意，通过引用文件获取，map结构int64对应文件引用ID；反之对应文件实体ID
func GetPublicURLsMap(args *ArgsGetPublicURLs) (data map[int64]string, err error) {
	data = map[int64]string{}
	//获取文件基本信息
	var claimList []BaseFileSys.FieldsFileClaimType
	var fileList []BaseFileSys.FieldsFileType
	claimList, fileList, _, err = BaseFileSys.GetFileByClaimIDs(&BaseFileSys.ArgsGetFileByClaimIDs{
		ClaimIDList: args.ClaimIDList,
		UserID:      args.UserID,
		OrgID:       args.OrgID,
		IsPublic:    args.IsPublic,
	})
	if err != nil {
		return
	}
	//遍历获取文件URL地址
	data, err = getPublicURLsMapAppend(claimList, fileList)
	return
}

// Deprecated: 准备废弃
// GetPublicURLStr 获取指定的文件URL
func GetPublicURLStr(fileID int64) (urlStr string, err error) {
	if fileID < 1 {
		err = errors.New("file id less 1")
		return
	}
	var data map[int64]string
	data, err = GetPublicURLsMap(&ArgsGetPublicURLs{
		ClaimIDList: []int64{fileID},
		UserID:      -1,
		OrgID:       -1,
		IsPublic:    false,
	})
	if err != nil {
		return
	}
	for _, v := range data {
		urlStr = v
		return
	}
	err = errors.New("no data")
	return
}

// Deprecated: 准备废弃
// GetPublicURLStrNoErr 无错误获取文件URL
func GetPublicURLStrNoErr(fileID int64) (urlStr string) {
	urlStr, _ = GetPublicURLStr(fileID)
	return
}

// GetPublicURLStrs 获取一组文件URLs
func GetPublicURLStrs(fileIDs []int64) (dataList []string, err error) {
	if len(fileIDs) < 1 {
		err = errors.New("no data")
		return
	}
	for k := 0; k < len(fileIDs); k++ {
		var urlStr string
		urlStr, err = GetPublicURLStr(fileIDs[k])
		if err != nil {
			err = nil
			continue
		}
		dataList = append(dataList, urlStr)
	}
	return
}

// GetPublicURLsMapByFile 通过文件实体获取URL
// 注意，反馈结构int64对应为文件实体ID
func GetPublicURLsMapByFile(args *ArgsGetPublicURLsByFile) (data map[int64]string, err error) {
	data = map[int64]string{}
	//获取文件基本信息
	var fileList []BaseFileSys.FieldsFileType
	fileList, _, err = BaseFileSys.GetFileByIDsAndClaim(&BaseFileSys.ArgsGetFileByIDsAndClaim{
		IDs:        args.FileIDs,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
		UserID:     args.UserID,
		OrgID:      args.OrgID,
		IsPublic:   args.IsPublic,
	})
	if err != nil {
		return
	}
	//遍历获取文件URL地址
	data, err = getPublicURLsMapAppend([]BaseFileSys.FieldsFileClaimType{}, fileList)
	if err != nil {
		return
	}
	return
}

func getPublicURLsMapAppend(claimList []BaseFileSys.FieldsFileClaimType, fileList []BaseFileSys.FieldsFileType) (data map[int64]string, err error) {
	data = map[int64]string{}
	//遍历获取文件URL地址
	if len(claimList) > 0 {
		for k := 0; k < len(claimList); k++ {
			vClaim := claimList[k]
			vFileData := BaseFileSys.FieldsFileType{}
			for _, vFile := range fileList {
				if vFile.ID == vClaim.FileID {
					vFileData = vFile
					break
				}
			}
			if vFileData.ID < 1 {
				continue
			}
			//获取访问URL路径
			var fileURL string
			fileURL, err = GetURLsByFileData(&vFileData)
			if err != nil {
				return
			}
			//写入数据组
			data[vClaim.ID] = fileURL
		}
	} else {
		for k := 0; k < len(fileList); k++ {
			vFileData := fileList[k]
			//获取访问URL路径
			var fileURL string
			fileURL, err = GetURLsByFileData(&vFileData)
			if err != nil {
				return
			}
			//写入数据组
			data[vFileData.ID] = fileURL
		}
	}
	return
}

// Deprecated: 建议采用GetURLsByFileData2
// GetURLsByFileData 直接从文件结构内找到URL地址
func GetURLsByFileData(fileData *BaseFileSys.FieldsFileType) (fileURL string, err error) {
	bucketName, b := fileData.Infos.GetVal("bucket")
	if !b {
		return
	}
	var qiniuBucketURL string
	qiniuBucketURL, err = getBucketURL(bucketName)
	if err != nil {
		return
	}
	//获取访问URL路径
	fileURL = storage.MakePublicURL(qiniuBucketURL, fileData.FromInfo.Mark)
	return
}

// GetURLsByFileData2 直接从文件结构内找到URL地址
func GetURLsByFileData2(bucket string, mark string) (fileURL string, err error) {
	var qiniuBucketURL string
	qiniuBucketURL, err = getBucketURL(bucket)
	if err != nil {
		return
	}
	//获取访问URL路径
	fileURL = storage.MakePublicURL(qiniuBucketURL, mark)
	return
}

// ArgsGetFileData 将文件加载到内存参数
type ArgsGetFileData struct {
	//引用文件ID
	ClaimID int64
}

// GetFileData 将文件加载到内存
func GetFileData(args *ArgsGetFileData) ([]byte, error) {
	//生成访问的URL地址
	fileURLs, err := GetPublicURLs(&ArgsGetPublicURLs{
		ClaimIDList: []int64{args.ClaimID},
		UserID:      0,
		OrgID:       0,
		IsPublic:    true,
	})
	if err != nil {
		return []byte{}, err
	}
	if len(fileURLs) < 1 {
		return []byte{}, errors.New("file urls is empty")
	}
	//爬取文件
	fileData, err := CoreHttp.HttpGet(fileURLs[0].URL)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprint("user try get file, but file src is not exist, file claim id: ", args.ClaimID, ", fileURL: ", fileURLs[0].URL, ", ", err))
	}
	//反馈
	return fileData, nil
}
