package ServiceUserInfo

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateInfo 创建新的文件参数
type ArgsCreateInfo struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" check:"isoTime"`
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

// CreateInfo 创建新的文件
func CreateInfo(args *ArgsCreateInfo) (data FieldsInfo, err error) {
	if args.BindID > 0 {
		//获取绑定人
		var bindData FieldsInfo
		err = Router2SystemConfig.MainDB.Get(&bindData, "SELECT id, bind_id FROM service_user_info WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.BindID, args.OrgID)
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
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_user_info", "INSERT INTO service_user_info (create_at, org_id, user_id, bind_id, bind_type, name, country, gender, id_card, id_card_front_file_id, id_card_back_file_id, phone, cover_file_id, des_files, address, date_of_birth, marital_status, education_status, profession, level, emergency_contact, emergency_contact_phone, sort_id, tags, doc_id, des, director_1, director_2, params) VALUES (:create_at,:org_id,:user_id,:bind_id,:bind_type,:name,:country,:gender,:id_card,:id_card_front_file_id,:id_card_back_file_id,:phone,:cover_file_id,:des_files,:address,:date_of_birth,:marital_status,:education_status,:profession,:level,:emergency_contact,:emergency_contact_phone,:sort_id,:tags,:doc_id,:des,:director_1,:director_2,:params)", args, &data)
	if err != nil {
		err = errors.New(fmt.Sprint("create new info, ", err))
		return
	}
	//推送nats
	pushNatsInfoStatus("create", data.ID)
	//统计数据
	pushNatsAnalysis(data.OrgID)
	//日志
	appendLog(&argsAppendLog{
		InfoID:     data.ID,
		OrgID:      data.OrgID,
		ChangeMark: "create",
		ChangeDes:  "创建档案",
		OldDes:     "",
		NewDes:     "",
	})
	//反馈
	return
}
