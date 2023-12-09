package BaseFileSys

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsCreate 创建新的文件信息参数
type ArgsCreate struct {
	//IP地址
	CreateIP string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//创建用户
	UserID int64
	//创建组织
	OrgID int64
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//文件大小
	FileSize int64
	//文件类型
	FileType string
	//文件Ｈash
	FileHash string
	//文件路径
	FileSrc string
	//过期时间
	// 留空则永不过期，过期则删除该文件实体
	ExpireAt time.Time
	//存储渠道
	FromInfo CoreSQLFrom.FieldsFrom
	//文件扩展信息
	Infos CoreSQLConfig.FieldsConfigsType
	//引用文件扩展信息
	ClaimInfos CoreSQLConfig.FieldsConfigsType
	//描述
	Des string
}

// Create 创建新的文件信息
// return FieldsFileClaimType 文件领用信息
// return bool 是否创建了新的底层文件，外部可以选择是否保留上传文件
// return error 错误信息
func Create(args *ArgsCreate) (data FieldsFileClaimType, b bool, errCode string, err error) {
	//修正参数
	if args.ClaimInfos == nil {
		args.ClaimInfos = CoreSQLConfig.FieldsConfigsType{}
	}
	//hash不能为空
	if args.FileHash == "" {
		errCode = "hash_empty"
		err = errors.New("hash is empty")
		return
	}
	//检查hash是否存在？
	var fileData FieldsFileType
	fileData, err = getFileByHash(args.FileHash)
	if err == nil && fileData.FromInfo.CheckEg(args.FromInfo) {
		if fileData.ID > 0 {
			data, errCode, err = ClaimFile(&ArgsClaimFile{
				FileID:     fileData.ID,
				UserID:     args.UserID,
				OrgID:      args.OrgID,
				IsPublic:   args.IsPublic,
				ExpireAt:   args.ExpireAt,
				ClaimInfos: args.ClaimInfos,
				Des:        args.Des,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("cannot claim hash file, fileID: ", fileData.ID, ", err: ", err))
				return
			}
			b = false
			return
		}
	}
	var newUpdateHash string
	newUpdateHash, err = CoreFilter.GetRandStr3(10)
	if err != nil {
		errCode = "update_hash"
		err = errors.New("rand hash, " + err.Error())
		return
	}
	//创建数据
	var newData FieldsFileType
	maps := map[string]interface{}{
		"update_hash": newUpdateHash,
		"create_ip":   args.CreateIP,
		"file_size":   args.FileSize,
		"file_type":   args.FileType,
		"file_hash":   args.FileHash,
		"file_src":    args.FileSrc,
		"infos":       args.Infos,
		"create_info": args.CreateInfo,
		"from_info":   args.FromInfo,
	}
	if err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_file", "INSERT INTO core_file (update_hash, create_ip, create_info, file_size, file_type, file_hash, file_src, from_info, infos) VALUES (:update_hash,:create_ip,:create_info,:file_size,:file_type,:file_hash,:file_src,:from_info,:infos)", maps, &newData); err != nil {
		errCode = "insert_file"
		err = errors.New("cannot create new file data, " + err.Error())
		return
	}
	if newData.ID < 1 {
		errCode = "insert_file"
		err = errors.New("cannot create new file data")
	}
	data, errCode, err = ClaimFile(&ArgsClaimFile{
		FileID:     newData.ID,
		UserID:     args.UserID,
		OrgID:      args.OrgID,
		IsPublic:   args.IsPublic,
		ExpireAt:   args.ExpireAt,
		ClaimInfos: args.ClaimInfos,
		Des:        args.Des,
	})
	if err != nil {
		b = true
		err = errors.New("cannot claim new file, " + err.Error())
		return
	}
	b = true
	return
}
