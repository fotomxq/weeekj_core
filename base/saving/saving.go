package BaseSaving

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetMark 获取指定的数据参数
type ArgsGetMark struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
}

// GetMark 获取指定的数据
func GetMark(args *ArgsGetMark) (val string, err error) {
	var data FieldsSaving
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, val FROM core_saving WHERE id = $1 AND mark = $2 AND expire_at >= NOW()", args.ID, args.Mark)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("data not exist")
		return
	}
	val = data.Val
	return
}

// ArgsUpdateMark 更新数据参数
type ArgsUpdateMark struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//数据集合
	Val string `db:"val" json:"val"`
}

// UpdateMark 更新数据
func UpdateMark(args *ArgsUpdateMark) (data FieldsSaving, err error) {
	defer runExpireBlocker.NewEdit()
	if args.ID > 0 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, val FROM core_saving WHERE id = $1 AND mark = $2", args.ID, args.Mark)
	}
	if err != nil || data.ID < 1 {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_saving", "INSERT INTO core_saving (expire_at, mark, val) VALUES (:expire_at,:mark,:val)", map[string]interface{}{
			"expire_at": CoreFilter.GetNowTimeCarbon().AddHour().Time,
			"mark":      args.Mark,
			"val":       args.Val,
		}, &data)
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_saving SET val = :val WHERE id = :id", map[string]interface{}{
		"id":  data.ID,
		"val": args.Val,
	})
	data.Val = args.Val
	return
}

// ArgsDeleteID 删除数据参数
type ArgsDeleteID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteID 删除数据参数
func DeleteID(args *ArgsDeleteID) (err error) {
	defer runExpireBlocker.NewEdit()
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_saving", "id", args)
	return
}
