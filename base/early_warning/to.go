package BaseEarlyWarning

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//预警名单

// 查看预警人名单
type ArgsGetToList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//搜索
	Search string
}

func GetToList(args *ArgsGetToList) (dataList []FieldsToType, dataCount int64, err error) {
	where := "(name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%' OR phone_nation_code ILIKE '%' || :search || '%' OR phone ILIKE '%' || :search || '%' OR email ILIKE '%' || :search || '%')"
	maps := map[string]interface{}{
		"search": args.Search,
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_ew_to",
		"id",
		"SELECT id, create_at, update_at, user_id, name, des, phone_nation_code, phone, email, user_id FROM core_ew_to WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at"},
	)
	return
}

// 查看预警人信息
type ArgsGetToByID struct {
	//ID
	ID int64
}

func GetToByID(args *ArgsGetToByID) (data FieldsToType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT  id, create_at, update_at, name, des, phone_nation_code, phone, email, user_id FROM core_ew_to WHERE id=$1", args.ID)
	return
}

// 通过用户ID找到预警人
type ArgsGetToByUserID struct {
	//用户ID
	UserID int64
}

func GetToByUserID(args *ArgsGetToByUserID) (data FieldsToType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, user_id, name, des, phone_nation_code, phone, email, user_id FROM core_ew_to WHERE user_id=$1", args.UserID)
	return
}

// 创建新的预警人
type ArgsCreateTo struct {
	//用户ID
	UserID int64 `db:"user_id"`
	//昵称
	Name string `db:"name"`
	//描述
	Des string `db:"des"`
	//联系电话
	PhoneNationCode string `db:"phone_nation_code"`
	Phone           string `db:"phone"`
	//邮件地址
	Email string `db:"email"`
}

func CreateTo(args *ArgsCreateTo) (data FieldsToType, err error) {
	data, err = GetToByUserID(&ArgsGetToByUserID{
		UserID: args.UserID,
	})
	if err == nil {
		err = errors.New("user id have bind to data")
		return
	}
	var lastID int64
	lastID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_ew_to(user_id, name, des, phone_nation_code, phone, email) VALUES(:user_id, :name, :des, :phone_nation_code, :phone, :email)", args)
	if err == nil {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, name, des, phone_nation_code, phone, email, user_id FROM core_ew_to WHERE id = $1", lastID)
	}
	return
}

// 修改预警人
type ArgsUpdateTo struct {
	//ID
	ID int64 `db:"id"`
	//用户ID
	UserID int64 `db:"user_id"`
	//昵称
	Name string `db:"name"`
	//描述
	Des string `db:"des"`
	//联系电话
	PhoneNationCode string `db:"phone_nation_code"`
	Phone           string `db:"phone"`
	//邮件地址
	Email string `db:"email"`
}

func UpdateTo(args *ArgsUpdateTo) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ew_to SET update_at=NOW(), user_id=:user_id, name=:name, des=:des, phone_nation_code=:phone_nation_code, phone=:phone, email=:email WHERE id=:id;", args)
	return
}

// 删除预警人
type ArgsDeleteToByID struct {
	//ID
	ID int64 `db:"id"`
}

func DeleteToByID(args *ArgsDeleteToByID) (err error) {
	//尝试解绑相关人员
	//无论是否成功，都会继续执行后续
	_ = SetUnBind(&ArgsSetUnBind{
		ToID: args.ID,
	})
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_ew_to", "id", args)
	return
}
