package UserAddress

import (
	"database/sql"
	"errors"
	"fmt"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	OrgUserMod "gitee.com/weeekj/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/jmoiron/sqlx"
)

// ArgsUpdate 修改地址参数
type ArgsUpdate struct {
	//地址ID
	ID int64 `db:"id" json:"id" check:"id"`
	//用户ID
	// 用于验证
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//地址昵称
	NiceName string `db:"nice_name" json:"niceName" check:"name" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province" check:"province"`
	//所属城市
	City int `db:"city" json:"city" check:"city"`
	//街道详细信息
	Address string `db:"address" json:"address" check:"address"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType" check:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude" check:"gps"`
	Latitude  float64 `db:"latitude" json:"latitude" check:"gps"`
	//联系人姓名
	Name string `db:"name" json:"name" check:"name"`
	//联系人国家代码
	NationCode string `db:"nation_code" json:"nationCode" check:"nationCode"`
	//联系人手机号
	Phone string `db:"phone" json:"phone" check:"phone"`
	//联系人邮箱
	Email string `db:"email" json:"email" check:"email" empty:"true"`
	//其他联系方式
	Infos CoreSQLConfig.FieldsInfosType `db:"infos" json:"infos"`
}

// Update 修改地址
func Update(args *ArgsUpdate) (err error) {
	//读取原始数据
	// 数据不能被删除，删除的数据无法更新
	var data FieldsAddress
	data, err = GetID(&ArgsGetID{
		ID:       args.ID,
		UserID:   args.UserID,
		IsRemove: false,
	})
	if err != nil {
		return
	}
	//检查数据是否完全相同，则不修改
	if fmt.Sprint(args.UserID, args.NiceName, args.Country, args.Province, args.City, args.Address, args.MapType, args.Longitude, args.Latitude, args.Name, args.Name, args.NationCode, args.Phone, args.Email, args.Infos) == fmt.Sprint(
		data.UserID, data.NiceName, data.Country, data.Province, data.City, data.Address, data.MapType, data.Longitude, data.Latitude, data.Name, data.NationCode, data.Phone, data.Email, data.Infos) {
		return
	}
	//退出捕捉
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	//创建新的数据
	tx := Router2SystemConfig.MainDB.MustBegin()
	var stmt *sqlx.NamedStmt
	stmt, err = tx.PrepareNamed("INSERT INTO user_address (create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos) VALUES (:create_at, NOW(), to_timestamp(0), 0, :user_id, :nice_name, :country, :province, :city, :address, :map_type, :longitude, :latitude, :name, :nation_code, :phone, :email, :infos) RETURNING id;")
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
			return
		}
		err = errors.New("prepare named " + err.Error())
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	var lastID int64
	err = stmt.Get(&lastID, map[string]interface{}{
		"create_at":   data.CreateAt,
		"parent_id":   data.ID,
		"user_id":     data.UserID,
		"nice_name":   args.NiceName,
		"country":     args.Country,
		"province":    args.Province,
		"city":        args.City,
		"address":     args.Address,
		"map_type":    args.MapType,
		"longitude":   args.Longitude,
		"latitude":    args.Latitude,
		"name":        args.Name,
		"nation_code": args.NationCode,
		"phone":       args.Phone,
		"email":       args.Email,
		"infos":       args.Infos,
	})
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
			return
		}
		err = errors.New("get last id " + err.Error())
		return
	}
	var newData FieldsAddress
	err = tx.Get(&newData, "SELECT id, create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos FROM user_address WHERE id = $1", lastID)
	if err != nil {
		return
	}
	//修改旧的数据为该数据的下级，同时标记删除
	var result sql.Result
	result, err = tx.NamedExec("UPDATE user_address SET delete_at = NOW(), parent_id = :parent_id WHERE id = :id", map[string]interface{}{
		"parent_id": newData.ID,
		"id":        data.ID,
	})
	if err != nil {
		return
	} else {
		var rowsAffected int64
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			return
		} else {
			if rowsAffected < 1 {
				err2 := tx.Rollback()
				if err2 != nil {
					err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
					return
				}
				err = errors.New("no update")
				return
			}
		}
	}
	//修改旧的数据的上级为新数据
	result, err = tx.NamedExec("UPDATE user_address SET parent_id = :parent_id WHERE parent_id = :id", map[string]interface{}{
		"parent_id": newData.ID,
		"id":        data.ID,
	})
	if err != nil {
		//return
	}
	//执行事务
	err = tx.Commit()
	//更新组织用户数据包
	OrgUserMod.PushUpdateUserData(0, args.UserID)
	//反馈
	return
}
