package ServiceUserInfo

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateInfo 修改文件信息参数
type ArgsUpdateInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//从属关系
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//从属关系的类型
	// 0 子女; 1 亲属非子女; 2 好友; 3 其他
	BindType int `db:"bind_type" json:"bindType" check:"intThan0" empty:"true"`
	//姓名
	Name string `db:"name" json:"name" check:"name"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender" check:"intThan0" empty:"true"`
	//证件编号
	IDCard string `db:"id_card" json:"idCard" check:"mark" empty:"true"`
	//身份证正反面照片
	IDCardFrontFileID int64 `db:"id_card_front_file_id" json:"idCardFrontFileID" check:"id" empty:"true"`
	IDCardBackFileID  int64 `db:"id_card_back_file_id" json:"idCardBackFileID" check:"id" empty:"true"`
	//联系电话
	Phone string `db:"phone" json:"phone" check:"phone" empty:"true"`
	//个人照片
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//其他照片信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//家庭住址
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//出生年月
	// date结构
	DateOfBirth time.Time `db:"date_of_birth" json:"dateOfBirth" check:"isoTime" empty:"true"`
	//婚姻状态
	MaritalStatus bool `db:"marital_status" json:"maritalStatus" check:"bool"`
	//教育状态
	// 0 无教育; 1 小学; 2 初中; 3 高中; 4 专科技校; 5 本科; 6 研究生; 7 博士; 8 博士后
	EducationStatus int `db:"education_status" json:"educationStatus" check:"intThan0" empty:"true"`
	//职业
	Profession string `db:"profession" json:"profession" check:"name" empty:"true"`
	//级别
	// 看护级别; 其他级别
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//主要紧急联系人
	EmergencyContact      string `db:"emergency_contact" json:"emergencyContact" check:"name" empty:"true"`
	EmergencyContactPhone string `db:"emergency_contact_phone" json:"emergencyContactPhone" check:"phone" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//关联的文件
	DocID int64 `db:"doc_id" json:"docID" check:"id" empty:"true"`
	//描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//负责人1
	Director1 int64 `db:"director_1" json:"director1" check:"id" empty:"true"`
	//负责人2
	Director2 int64 `db:"director_2" json:"director2" check:"id" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateInfo 修改文件信息
func UpdateInfo(args *ArgsUpdateInfo) (err error) {
	//获取旧的数据
	oldData := getInfoID(args.ID)
	if oldData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//检查该绑定人
	if args.BindID > 0 {
		//获取绑定人
		var bindData FieldsInfo
		err = Router2SystemConfig.MainDB.Get(&bindData, "SELECT id, bind_id FROM service_user_info WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.BindType, args.OrgID)
		if err != nil || bindData.ID < 1 {
			err = errors.New(fmt.Sprint("bind not exist, ", err))
			return
		}
		//绑定人不能存在其他关系
		if bindData.BindID > 0 {
			err = errors.New("bind have other bind")
			return
		}
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info SET update_at = NOW(), create_at = :create_at, user_id = :user_id, bind_id = :bind_id, bind_type = :bind_type, name = :name, country = :country, gender = :gender, id_card = :id_card, id_card_front_file_id = :id_card_front_file_id, id_card_back_file_id = :id_card_back_file_id, phone = :phone, cover_file_id = :cover_file_id, des_files = :des_files, address = :address, date_of_birth = :date_of_birth, marital_status = :marital_status, education_status = :education_status, profession = :profession, level = :level, emergency_contact = :emergency_contact, emergency_contact_phone = :emergency_contact_phone, sort_id = :sort_id, tags = :tags, doc_id = :doc_id, des = :des, director_1 = :director_1, director_2 = :director_2, params = :params WHERE id = :id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	//缓冲
	deleteInfoCache(args.ID)
	//推送nats
	pushNatsInfoStatus("update", args.ID)
	//统计数据
	pushNatsAnalysis(args.OrgID)
	//新的数据
	newData := getInfoID(args.ID)
	if newData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//对比数据，生成日志
	oldDes := ""
	newDes := ""
	if args.CreateAt.Unix() != oldData.CreateAt.Unix() {
		oldDes = fmt.Sprint(oldDes, "旧的时间: ", oldData.CreateAt)
		newDes = fmt.Sprint(newDes, "新的时间: ", args.CreateAt)
	}
	if args.BindType != oldData.BindType {
		oldDes = fmt.Sprint(oldDes, "; 旧的从属关系类型: ", oldData.BindType)
		newDes = fmt.Sprint(newDes, "; 新的从属关系类型: ", args.BindType)
	}
	if args.Name != oldData.Name {
		oldDes = fmt.Sprint(oldDes, "; 旧的姓名: ", oldData.Name)
		newDes = fmt.Sprint(newDes, "; 新的姓名: ", args.Name)
	}
	if args.Country != oldData.Country {
		oldDes = fmt.Sprint(oldDes, "; 旧的国家: ", oldData.Country)
		newDes = fmt.Sprint(newDes, "; 新的国家: ", args.Country)
	}
	if args.Gender != oldData.Gender {
		oldDes = fmt.Sprint(oldDes, "; 旧的性别: ", GetInfoGender(oldData.Gender))
		newDes = fmt.Sprint(newDes, "; 新的性别: ", GetInfoGender(args.Gender))
	}
	if args.IDCard != oldData.IDCard {
		oldDes = fmt.Sprint(oldDes, "; 旧的身份证号: ", oldData.IDCard)
		newDes = fmt.Sprint(newDes, "; 新的身份证号: ", args.IDCard)
	}
	if args.IDCardFrontFileID != oldData.IDCardFrontFileID {
		oldDes = fmt.Sprint(oldDes, "; 旧的身份证正面文件ID: ", oldData.IDCardFrontFileID)
		newDes = fmt.Sprint(newDes, "; 新的身份证正面文件ID: ", args.IDCardFrontFileID)
	}
	if args.IDCardBackFileID != oldData.IDCardBackFileID {
		oldDes = fmt.Sprint(oldDes, "; 旧的身份证背面文件ID: ", oldData.IDCardBackFileID)
		newDes = fmt.Sprint(newDes, "; 新的身份证背面文件ID: ", args.IDCardBackFileID)
	}
	if args.Phone != oldData.Phone {
		oldDes = fmt.Sprint(oldDes, "; 旧的联系电话: ", oldData.Phone)
		newDes = fmt.Sprint(newDes, "; 新的联系电话: ", args.Phone)
	}
	if args.CoverFileID != oldData.CoverFileID {
		oldDes = fmt.Sprint(oldDes, "; 旧的个人照片文件ID: ", oldData.CoverFileID)
		newDes = fmt.Sprint(newDes, "; 新的个人照片文件ID: ", args.CoverFileID)
	}
	if len(args.DesFiles) != len(oldData.DesFiles) {
		oldDes = fmt.Sprint(oldDes, "; 旧的附加文件列: ", oldData.DesFiles)
		newDes = fmt.Sprint(newDes, "; 新的附加文件列: ", args.DesFiles)
	}
	if args.Address != oldData.Address {
		oldDes = fmt.Sprint(oldDes, "; 旧的地址: ", oldData.Address)
		newDes = fmt.Sprint(newDes, "; 新的地址: ", args.Address)
	}
	if args.DateOfBirth != oldData.DateOfBirth {
		oldDes = fmt.Sprint(oldDes, "; 旧的出生年月: ", oldData.DateOfBirth)
		newDes = fmt.Sprint(newDes, "; 新的出生年月: ", args.DateOfBirth)
	}
	if args.MaritalStatus != oldData.MaritalStatus {
		oldDes = fmt.Sprint(oldDes, "; 旧的婚姻状态: ", GetInfoMaritalStatus(oldData.MaritalStatus))
		newDes = fmt.Sprint(newDes, "; 新的婚姻状态: ", GetInfoMaritalStatus(args.MaritalStatus))
	}
	if args.EducationStatus != oldData.EducationStatus {
		oldDes = fmt.Sprint(oldDes, "; 旧的教育水平: ", GetInfoEducationStatus(oldData.EducationStatus))
		newDes = fmt.Sprint(newDes, "; 新的教育水平: ", GetInfoEducationStatus(args.EducationStatus))
	}
	if args.Profession != oldData.Profession {
		oldDes = fmt.Sprint(oldDes, "; 旧的职业: ", oldData.Profession)
		newDes = fmt.Sprint(newDes, "; 新的职业: ", args.Profession)
	}
	if args.Level != oldData.Level {
		oldDes = fmt.Sprint(oldDes, "; 旧的服务级别: ", oldData.Level)
		newDes = fmt.Sprint(newDes, "; 新的服务级别: ", args.Level)
	}
	if args.EmergencyContact != oldData.EmergencyContact {
		oldDes = fmt.Sprint(oldDes, "; 旧的紧急联系人姓名: ", oldData.EmergencyContact)
		newDes = fmt.Sprint(newDes, "; 新的紧急联系人姓名: ", args.EmergencyContact)
	}
	if args.SortID != oldData.SortID {
		oldDes = fmt.Sprint(oldDes, "; 旧的分类: ", oldData.SortID)
		newDes = fmt.Sprint(newDes, "; 新的分类: ", args.SortID)
	}
	if len(args.Tags) != len(oldData.Tags) {
		oldDes = fmt.Sprint(oldDes, "; 旧的标签: ", oldData.Tags)
		newDes = fmt.Sprint(newDes, "; 新的标签: ", args.Tags)
	}
	if args.DocID != oldData.DocID {
		oldDes = fmt.Sprint(oldDes, "; 旧的主要关联文档: ", oldData.DocID)
		newDes = fmt.Sprint(newDes, "; 新的主要关联文档: ", args.DocID)
	}
	if args.Des != oldData.Des {
		oldDes = fmt.Sprint(oldDes, "; 旧的描述 ", oldData.Des)
		newDes = fmt.Sprint(newDes, "; 新的描述: ", args.Des)
	}
	if args.Director1 != oldData.Director1 {
		oldDes = fmt.Sprint(oldDes, "; 旧的主要负责人 ", oldData.Director1)
		newDes = fmt.Sprint(newDes, "; 新的主要负责人: ", args.Director1)
	}
	if args.Director2 != oldData.Director2 {
		oldDes = fmt.Sprint(oldDes, "; 旧的次要负责人 ", oldData.Director2)
		newDes = fmt.Sprint(newDes, "; 新的次要负责人: ", args.Director2)
	}
	for _, v := range oldData.Params {
		isFind := false
		for _, v2 := range newData.Params {
			if v.Mark == v2.Mark {
				if v.Val != v2.Val {
					oldDes = fmt.Sprint(oldDes, "; 旧的扩展信息.[ ", v.Mark, "]: ", v.Val)
					newDes = fmt.Sprint(newDes, "; 新的扩展信息.[ ", v.Mark, "]: ", v2.Val)
				}
				isFind = true
				break
			}
		}
		if !isFind {
			newDes = fmt.Sprint(newDes, "; 移除了扩展信息.[ ", v.Mark, "]: ", v.Val)
		}
	}
	for _, v := range newData.Params {
		isFind := false
		for _, v2 := range oldData.Params {
			if v.Mark == v2.Mark {
				isFind = true
				break
			}
		}
		if !isFind {
			newDes = fmt.Sprint(newDes, "; 新增了扩展信息.[ ", v.Mark, "]: ", v.Val)
		}
	}
	//日志
	appendLog(&argsAppendLog{
		InfoID:     args.ID,
		OrgID:      oldData.OrgID,
		ChangeMark: "update",
		ChangeDes:  "编辑档案",
		OldDes:     oldDes,
		NewDes:     newDes,
	})
	//反馈
	return
}

// ArgsUpdateInfoDie 标记死亡参数
type ArgsUpdateInfoDie struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//时间
	DieAt string `json:"dieAt" check:"isoTime"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateInfoDie 标记死亡
func UpdateInfoDie(args *ArgsUpdateInfoDie) (err error) {
	var data FieldsInfo
	data, err = GetInfoID(&ArgsGetInfoID{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	if CoreSQL.CheckTimeHaveData(data.DieAt) {
		return
	}
	var dieAt time.Time
	dieAt, err = CoreFilter.GetTimeByISO(args.DieAt)
	if err != nil {
		return
	}
	for _, v := range args.Params {
		data.Params = CoreSQLConfig.Set(data.Params, v.Mark, v.Val)
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info SET update_at = NOW(), die_at = :die_at, out_at = :die_at, params = :params WHERE id = :id", map[string]interface{}{
		"id":     args.ID,
		"die_at": dieAt,
		"params": data.Params,
	})
	if err != nil {
		return
	}
	//缓冲
	deleteInfoCache(args.ID)
	//推送nats
	pushNatsInfoStatus("die", args.ID)
	//统计数据
	pushNatsAnalysis(args.OrgID)
	//日志
	appendLog(&argsAppendLog{
		InfoID:     args.ID,
		OrgID:      data.OrgID,
		ChangeMark: "die",
		ChangeDes:  "档案死亡",
		OldDes:     "",
		NewDes:     "",
	})
	//反馈
	return
}

// ArgsUpdateInfoOut 标记出院参数
type ArgsUpdateInfoOut struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//时间
	OutAt string `json:"outAt" check:"isoTime"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateInfoOut 标记出院
func UpdateInfoOut(args *ArgsUpdateInfoOut) (err error) {
	var data FieldsInfo
	data, err = GetInfoID(&ArgsGetInfoID{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	var outAt time.Time
	outAt, err = CoreFilter.GetTimeByISO(args.OutAt)
	if err != nil {
		return
	}
	for _, v := range args.Params {
		data.Params = CoreSQLConfig.Set(data.Params, v.Mark, v.Val)
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info SET update_at = NOW(), out_at = :out_at, params = :params WHERE id = :id", map[string]interface{}{
		"id":     args.ID,
		"out_at": outAt,
		"params": data.Params,
	})
	if err != nil {
		return
	}
	//缓冲
	deleteInfoCache(args.ID)
	//推送nats
	pushNatsInfoStatus("out", args.ID)
	//统计数据
	pushNatsAnalysis(args.OrgID)
	//日志
	appendLog(&argsAppendLog{
		InfoID:     args.ID,
		OrgID:      data.OrgID,
		ChangeMark: "out",
		ChangeDes:  "档案离开",
		OldDes:     "",
		NewDes:     "",
	})
	//反馈
	return
}

// ArgsUpdateInfoDoc 仅修改信息的关联文档参数
type ArgsUpdateInfoDoc struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//关联的文件
	DocID int64 `db:"doc_id" json:"docID" check:"id" empty:"true"`
}

// UpdateInfoDoc 仅修改信息的关联文档
func UpdateInfoDoc(args *ArgsUpdateInfoDoc) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info SET update_at = NOW(), doc_id = :doc_id WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//缓冲
	deleteInfoCache(args.ID)
	//统计数据
	pushNatsAnalysis(args.OrgID)
	//日志
	appendLog(&argsAppendLog{
		InfoID:     args.ID,
		OrgID:      0,
		ChangeMark: "doc",
		ChangeDes:  "创建档案的附属档案",
		OldDes:     "",
		NewDes:     fmt.Sprint("创建档案的附属档案，附属档案ID:", args.DocID),
	})
	//反馈
	return
}
