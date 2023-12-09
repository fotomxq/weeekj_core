package MapMapDrawing

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "gitee.com/weeekj/weeekj_core/v5/core/sql/gps"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetPointByPic 获取指定图的所有点信息参数
type ArgsGetPointByPic struct {
	//主图ID
	PicID int64 `db:"pic_id" json:"picID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetPointByPic 获取指定图的所有点信息
func GetPointByPic(args *ArgsGetPointByPic) (dataList []FieldsPoint, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, org_id, pic_id, pic_point, gps_point, radius, cover_file_id, cover_icon, cover_rgb, bind_mark, bind_id, params FROM map_map_drawing_point WHERE org_id = $1 AND pic_id = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.PicID)
	return
}

// ArgsCreatePoint 创建新的点参数
type ArgsCreatePoint struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//主图
	PicID int64 `db:"pic_id" json:"picID" check:"id"`
	//图坐标
	PicPoint CoreSQLGPS.FieldsPoint `db:"pic_point" json:"picPoint"`
	//GPS坐标
	GPSPoint CoreSQLGPS.FieldsPoint `db:"gps_point" json:"gpsPoint"`
	//误差半径
	Radius float64 `db:"radius" json:"radius" check:"floatThan0"`
	//显示图标
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//显示位置点系统图标
	CoverIcon string `db:"cover_icon" json:"coverIcon"`
	//图标颜色
	CoverRGB string `db:"cover_rgb" json:"coverRGB"`
	//绑定Mark
	// 绑定的系统，如pic 主图系统; room 房间
	BindMark string `db:"bind_mark" json:"bindMark" check:"mark"`
	//绑定ID
	// 该点绑定的房间，可作为联动处理
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreatePoint 创建新的点
func CreatePoint(args *ArgsCreatePoint) (data FieldsPoint, err error) {
	var picData FieldsPic
	err = Router2SystemConfig.MainDB.Get(&picData, "SELECT id FROM map_map_drawing_pic WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.PicID, args.OrgID)
	if err != nil || picData.ID < 1 {
		err = errors.New(fmt.Sprint("pic not exist, ", err))
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "map_map_drawing_point", "INSERT INTO map_map_drawing_point (org_id, pic_id, pic_point, gps_point, radius, cover_file_id, cover_icon, cover_rgb, bind_mark, bind_id, params) VALUES (:org_id,:pic_id,:pic_point,:gps_point,:radius,:cover_file_id,:cover_icon,:cover_rgb,:bind_mark,:bind_id,:params)", args, &data)
	return
}

// ArgsUpdatePoint 修改点参数
type ArgsUpdatePoint struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//图坐标
	PicPoint CoreSQLGPS.FieldsPoint `db:"pic_point" json:"picPoint"`
	//GPS坐标
	GPSPoint CoreSQLGPS.FieldsPoint `db:"gps_point" json:"gpsPoint"`
	//误差半径
	Radius float64 `db:"radius" json:"radius" check:"floatThan0"`
	//显示图标
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//显示位置点系统图标
	CoverIcon string `db:"cover_icon" json:"coverIcon"`
	//图标颜色
	CoverRGB string `db:"cover_rgb" json:"coverRGB"`
	//绑定Mark
	// 绑定的系统，如pic 主图系统; room 房间
	BindMark string `db:"bind_mark" json:"bindMark" check:"mark"`
	//绑定ID
	// 该点绑定的房间，可作为联动处理
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdatePoint 修改点
func UpdatePoint(args *ArgsUpdatePoint) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_map_drawing_point SET update_at = NOW(), pic_point = :pic_point, gps_point = :gps_point, radius = :radius, cover_file_id = :cover_file_id, cover_icon = :cover_icon, cover_rgb = :cover_rgb, bind_mark = :bind_mark, bind_id = :bind_id, params = :params WHERE id = :id AND org_id = :org_id", args)
	return
}

// ArgsDeletePoint 删除点参数
type ArgsDeletePoint struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeletePoint 删除点
func DeletePoint(args *ArgsDeletePoint) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "map_map_drawing_point", "id = :id AND org_id = :org_id", args)
	return
}
