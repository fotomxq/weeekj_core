package OrgCoreCore

import (
	"errors"
	"fmt"
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	BaseQiniu "github.com/fotomxq/weeekj_core/v5/base/qiniu"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetBindInfoList 获取列表参数
type ArgsGetBindInfoList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 必填
	OrgID int64 `json:"orgID" check:"id"`
	//分组ID
	GroupID int64 `json:"groupID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataGetBindInfoList struct {
	//创建时间
	CreateAt string `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt string `db:"update_at" json:"updateAt"`
	//入职时间
	InAt string `db:"in_at" json:"inAt"`
	//离职时间
	OutAt string `db:"out_at" json:"outAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender"`
	//出生年月
	// date结构
	DateOfBirth string `db:"date_of_birth" json:"dateOfBirth"`
	//联系电话
	Phone string `db:"phone" json:"phone"`
	//证件编号
	IDCard string `db:"id_card" json:"idCard"`
	//个人照片
	CoverFileID  int64  `db:"cover_file_id" json:"coverFileID"`
	CoverFileURL string `json:"coverFileURL"`
	//其他照片信息
	DesFiles    pq.Int64Array `db:"des_files" json:"desFiles"`
	DesFileURLs []string      `json:"desFileURLs"`
	//教育状态
	// 0 无教育; 1 小学; 2 初中; 3 高中; 4 专科技校; 5 本科; 6 研究生; 7 博士; 8 博士后
	EducationStatus int `db:"education_status" json:"educationStatus"`
	//婚姻状态
	MaritalStatus bool `db:"marital_status" json:"maritalStatus"`
	//家庭住址
	Address string `db:"address" json:"address"`
	//资质
	CertName string `db:"cert_name" json:"certName"`
	//保险
	InsuranceType string `db:"insurance_type" json:"insuranceType"`
	InsuranceAt   string `db:"insurance_at" json:"insuranceAt"`
	//银行信息
	BandName string `db:"band_name" json:"bindName"`
	BandSN   string `db:"band_sn" json:"bandSN"`
	//描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
}

// GetBindInfoList 获取列表
func GetBindInfoList(args *ArgsGetBindInfoList) (dataList []DataGetBindInfoList, dataCount int64, err error) {
	var bindList []FieldsBind
	bindList, dataCount, _ = GetBindList(&ArgsGetBindList{
		Pages:    args.Pages,
		OrgID:    args.OrgID,
		UserID:   -1,
		GroupID:  args.GroupID,
		Manager:  "",
		IsRemove: args.IsRemove,
		Search:   args.Search,
	})
	if len(bindList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, vBind := range bindList {
		vData := getBindInfoByBindID(vBind.ID)
		var desFileURLs []string
		for _, v2 := range vData.DesFiles {
			desFileURLs = append(desFileURLs, BaseQiniu.GetPublicURLStrNoErr(v2))
		}
		dataList = append(dataList, DataGetBindInfoList{
			CreateAt:        CoreFilter.GetTimeToDefaultTime(vData.CreateAt),
			UpdateAt:        CoreFilter.GetTimeToDefaultTime(vData.UpdateAt),
			InAt:            CoreFilter.GetTimeToDefaultTime(vData.InAt),
			OutAt:           CoreFilter.GetTimeToDefaultTime(vData.OutAt),
			OrgID:           vBind.OrgID,
			OrgBindID:       vBind.ID,
			Gender:          vData.Gender,
			DateOfBirth:     CoreFilter.GetTimeToDefaultDate(vData.DateOfBirth),
			Phone:           vData.Phone,
			IDCard:          vData.IDCard,
			CoverFileID:     vData.CoverFileID,
			CoverFileURL:    BaseQiniu.GetPublicURLStrNoErr(vData.CoverFileID),
			DesFiles:        vData.DesFiles,
			DesFileURLs:     desFileURLs,
			EducationStatus: vData.EducationStatus,
			MaritalStatus:   vData.MaritalStatus,
			Address:         vData.Address,
			CertName:        vData.CertName,
			InsuranceType:   vData.InsuranceType,
			InsuranceAt:     CoreFilter.GetTimeToDefaultTime(vData.InsuranceAt),
			BandName:        vData.BandName,
			BandSN:          vData.BandSN,
			Des:             vData.Des,
		})
	}
	return
}

// DataGetBindInfo 获取指定BindID的信息数据
type DataGetBindInfo struct {
	//创建时间
	CreateAt string `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt string `db:"update_at" json:"updateAt"`
	//入职时间
	InAt string `db:"in_at" json:"inAt"`
	//离职时间
	OutAt string `db:"out_at" json:"outAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender"`
	//出生年月
	// date结构
	DateOfBirth string `db:"date_of_birth" json:"dateOfBirth"`
	//联系电话
	Phone string `db:"phone" json:"phone"`
	//证件编号
	IDCard string `db:"id_card" json:"idCard"`
	//个人照片
	CoverFileID  int64  `db:"cover_file_id" json:"coverFileID"`
	CoverFileURL string `json:"coverFileURL"`
	//其他照片信息
	DesFiles    pq.Int64Array `db:"des_files" json:"desFiles"`
	DesFileURLs []string      `json:"desFileURLs"`
	//教育状态
	// 0 无教育; 1 小学; 2 初中; 3 高中; 4 专科技校; 5 本科; 6 研究生; 7 博士; 8 博士后
	EducationStatus int `db:"education_status" json:"educationStatus"`
	//婚姻状态
	MaritalStatus bool `db:"marital_status" json:"maritalStatus"`
	//家庭住址
	Address string `db:"address" json:"address"`
	//资质
	CertName string `db:"cert_name" json:"certName"`
	//保险
	InsuranceType string `db:"insurance_type" json:"insuranceType"`
	InsuranceAt   string `db:"insurance_at" json:"insuranceAt"`
	//银行信息
	BandName string `db:"band_name" json:"bindName"`
	BandSN   string `db:"band_sn" json:"bandSN"`
	//描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
}

// GetBindInfo 获取指定BindID的信息
func GetBindInfo(bindID int64) (data DataGetBindInfo, err error) {
	rawData := getBindInfoByBindID(bindID)
	if rawData.ID < 1 {
		data = DataGetBindInfo{}
		err = errors.New("no data")
		return
	}
	data = DataGetBindInfo{
		CreateAt:        CoreFilter.GetTimeToDefaultTime(rawData.CreateAt),
		UpdateAt:        CoreFilter.GetTimeToDefaultTime(rawData.UpdateAt),
		InAt:            CoreFilter.GetTimeToDefaultDate(rawData.InAt),
		OutAt:           CoreFilter.GetTimeToDefaultDate(rawData.OutAt),
		OrgID:           rawData.OrgID,
		OrgBindID:       rawData.ID,
		Gender:          rawData.Gender,
		DateOfBirth:     CoreFilter.GetTimeToDefaultDate(rawData.DateOfBirth),
		Phone:           rawData.Phone,
		IDCard:          rawData.IDCard,
		CoverFileID:     rawData.CoverFileID,
		CoverFileURL:    BaseFileSys2.GetPublicURLByClaimID(rawData.CoverFileID),
		DesFiles:        rawData.DesFiles,
		DesFileURLs:     BaseFileSys2.GetPublicURLsByClaimIDs(rawData.DesFiles),
		EducationStatus: rawData.EducationStatus,
		MaritalStatus:   rawData.MaritalStatus,
		Address:         rawData.Address,
		CertName:        rawData.CertName,
		InsuranceType:   rawData.InsuranceType,
		InsuranceAt:     CoreFilter.GetTimeToDefaultDate(rawData.InsuranceAt),
		BandName:        rawData.BandName,
		BandSN:          rawData.BandSN,
		Des:             rawData.Des,
	}
	return
}

func getBindInfoByBindID(bindID int64) (data FieldsBindInfo) {
	cacheMark := getBindInfoCacheMark(bindID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, in_at, out_at, org_id, org_bind_id, gender, date_of_birth, phone, id_card, cover_file_id, des_files, education_status, marital_status, address, cert_name, insurance_type, insurance_at, band_name, band_sn, des FROM org_core_bind_info WHERE org_bind_id = $1", bindID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, bindInfoCacheTime)
	return
}

// ArgsSetBindInfo 修改信息参数
type ArgsSetBindInfo struct {
	//入职时间
	InAt string `db:"in_at" json:"inAt"`
	//离职时间
	OutAt string `db:"out_at" json:"outAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender"`
	//出生年月
	// date结构
	DateOfBirth string `db:"date_of_birth" json:"dateOfBirth"`
	//联系电话
	Phone string `db:"phone" json:"phone"`
	//证件编号
	IDCard string `db:"id_card" json:"idCard"`
	//个人照片
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//其他照片信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//教育状态
	// 0 无教育; 1 小学; 2 初中; 3 高中; 4 专科技校; 5 本科; 6 研究生; 7 博士; 8 博士后
	EducationStatus int `db:"education_status" json:"educationStatus"`
	//婚姻状态
	MaritalStatus bool `db:"marital_status" json:"maritalStatus"`
	//家庭住址
	Address string `db:"address" json:"address"`
	//资质
	CertName string `db:"cert_name" json:"certName"`
	//保险
	InsuranceType string `db:"insurance_type" json:"insuranceType"`
	InsuranceAt   string `db:"insurance_at" json:"insuranceAt"`
	//银行信息
	BandName string `db:"band_name" json:"bindName"`
	BandSN   string `db:"band_sn" json:"bandSN"`
	//描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
}

// SetBindInfo 修改信息
func SetBindInfo(args *ArgsSetBindInfo) (err error) {
	var inAt, outAt, dateOfBirth, insuranceAt time.Time
	inAt, err = CoreFilter.GetTimeByDefault(args.InAt)
	if err != nil {
		err = errors.New(fmt.Sprint("in at, ", err))
		return
	}
	if args.OutAt != "" {
		outAt, err = CoreFilter.GetTimeByDefault(args.OutAt)
		if err != nil {
			err = errors.New(fmt.Sprint("out at, ", err))
			return
		}
	}
	dateOfBirth, err = CoreFilter.GetTimeByDefault(args.DateOfBirth)
	if err != nil {
		err = errors.New(fmt.Sprint("date of birth, ", err))
		return
	}
	if args.InsuranceAt != "" {
		insuranceAt, err = CoreFilter.GetTimeByDefault(args.InsuranceAt)
		if err != nil {
			err = errors.New(fmt.Sprint("insurance at, ", err))
			return
		}
	}
	bindInfoData := getBindInfoByBindID(args.OrgBindID)
	if bindInfoData.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind_info SET update_at = NOW(), in_at = :in_at, out_at = :out_at, gender = :gender, date_of_birth = :date_of_birth, phone = :phone, id_card = :id_card, cover_file_id = :cover_file_id, des_files = :des_files, education_status = :education_status, marital_status = :marital_status, address = :address, cert_name = :cert_name, insurance_type = :insurance_type, insurance_at = :insurance_at, band_name = :band_name, band_sn = :band_sn, des = :des WHERE id = :id", map[string]interface{}{
			"id":               bindInfoData.ID,
			"in_at":            inAt,
			"out_at":           outAt,
			"gender":           args.Gender,
			"date_of_birth":    dateOfBirth,
			"phone":            args.Phone,
			"id_card":          args.IDCard,
			"cover_file_id":    args.CoverFileID,
			"des_files":        args.DesFiles,
			"education_status": args.EducationStatus,
			"marital_status":   args.MaritalStatus,
			"address":          args.Address,
			"cert_name":        args.CertName,
			"insurance_type":   args.InsuranceType,
			"insurance_at":     insuranceAt,
			"band_name":        args.BandName,
			"band_sn":          args.BandSN,
			"des":              args.Des,
		})
		if err != nil {
			return
		}
		deleteBindInfoCache(bindInfoData.OrgBindID)
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_core_bind_info (in_at, out_at, org_id, org_bind_id, gender, date_of_birth, phone, id_card, cover_file_id, des_files, education_status, marital_status, address, cert_name, insurance_type, insurance_at, band_name, band_sn, des) VALUES (:in_at,:out_at,:org_id,:org_bind_id,:gender,:date_of_birth,:phone,:id_card,:cover_file_id,:des_files,:education_status,:marital_status,:address,:cert_name,:insurance_type,:insurance_at,:band_name,:band_sn,:des)", map[string]interface{}{
			"org_id":           args.OrgID,
			"org_bind_id":      args.OrgBindID,
			"in_at":            inAt,
			"out_at":           outAt,
			"gender":           args.Gender,
			"date_of_birth":    dateOfBirth,
			"phone":            args.Phone,
			"id_card":          args.IDCard,
			"cover_file_id":    args.CoverFileID,
			"des_files":        args.DesFiles,
			"education_status": args.EducationStatus,
			"marital_status":   args.MaritalStatus,
			"address":          args.Address,
			"cert_name":        args.CertName,
			"insurance_type":   args.InsuranceType,
			"insurance_at":     insuranceAt,
			"band_name":        args.BandName,
			"band_sn":          args.BandSN,
			"des":              args.Des,
		})
		if err != nil {
			return
		}
	}
	return
}

// 仅更新手机号等基本信息
func updateBindInfoBaseInBind(orgBindID int64, phone string) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind_info SET phone = :phone WHERE org_bind_id = :org_bind_id", map[string]interface{}{
		"org_bind_id": orgBindID,
		"phone":       phone,
	})
	if err != nil {
		return
	}
	deleteBindInfoCache(orgBindID)
	return
}

// 缓冲
func getBindInfoCacheMark(bindID int64) string {
	return fmt.Sprint("org:core:bind:info:bind:", bindID)
}

func deleteBindInfoCache(bindID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBindInfoCacheMark(bindID))
}
