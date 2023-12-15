package MapMapDrawing

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetPicList 获取主图列表参数
type ArgsGetPicList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//主图ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetPicList 获取主图列表
func GetPicList(args *ArgsGetPicList) (dataList []FieldsPic, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "map_map_drawing_pic"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, parent_id, name, des, file_id, fix_height, fix_width, button_name, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetPicID 查看主图参数
type ArgsGetPicID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetPicID 查看主图
func GetPicID(args *ArgsGetPicID) (data FieldsPic, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, parent_id, name, des, file_id, fix_height, fix_width, button_name, bind_area_id, params FROM map_map_drawing_pic WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	return
}

// ArgsGetPicMore 获取一组ID参数
type ArgsGetPicMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetPicMore 获取一组ID
func GetPicMore(args *ArgsGetPicMore) (dataList []FieldsPic, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "map_map_drawing_pic", "id, create_at, update_at, delete_at, org_id, parent_id, name, des, file_id, fix_height, fix_width, button_name, bind_area_id, params", args.IDs, args.HaveRemove)
	return
}

// ArgsGetPicChild 获取指定主图的所有幅图参数
type ArgsGetPicChild struct {
	//主图ID
	ParentID int64 `db:"parent_id" json:"parentID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetPicChild 获取指定主图的所有幅图
func GetPicChild(args *ArgsGetPicChild) (dataList []FieldsPic, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, org_id, parent_id, name, des, file_id, fix_height, fix_width, button_name, bind_area_id, params FROM map_map_drawing_pic WHERE parent_id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.ParentID, args.OrgID)
	return
}

// ArgsCreatePic 创建新的主图参数
type ArgsCreatePic struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//上级ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//核心地图
	FileID int64 `db:"file_id" json:"fileID" check:"id"`
	//修正图片高度和宽度
	FixHeight int `db:"fix_height" json:"fixHeight" check:"intThan0"`
	FixWidth  int `db:"fix_width" json:"fixWidth" check:"intThan0"`
	//按钮文字
	ButtonName string `db:"button_name" json:"buttonName" check:"name"`
	//绑定电子围栏
	BindAreaID int64 `db:"bind_area_id" json:"bindAreaID" check:"id" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreatePic 创建新的主图
func CreatePic(args *ArgsCreatePic) (data FieldsPic, err error) {
	//不能存在多个级别
	if args.ParentID > 0 {
		var parentData FieldsPic
		err = Router2SystemConfig.MainDB.Get(&parentData, "SELECT id FROM map_map_drawing_pic WHERE id = $1 AND org_id = $2 AND parent_id < 1 AND delete_at < to_timestamp(1000000)", args.ParentID, args.OrgID)
		if err != nil || parentData.ID < 1 {
			err = errors.New(fmt.Sprint("parent not exist, ", err))
			return
		}
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "map_map_drawing_pic", "INSERT INTO map_map_drawing_pic (org_id, parent_id, name, des, file_id, fix_height, fix_width, button_name, bind_area_id, params) VALUES (:org_id,:parent_id,:name,:des,:file_id,:fix_height,:fix_width,:button_name,:bind_area_id,:params)", args, &data)
	return
}

// ArgsUpdatePic 修改主图参数
type ArgsUpdatePic struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//核心地图
	FileID int64 `db:"file_id" json:"fileID" check:"id"`
	//修正图片高度和宽度
	FixHeight int `db:"fix_height" json:"fixHeight" check:"intThan0"`
	FixWidth  int `db:"fix_width" json:"fixWidth" check:"intThan0"`
	//按钮文字
	ButtonName string `db:"button_name" json:"buttonName" check:"name"`
	//绑定电子围栏
	BindAreaID int64 `db:"bind_area_id" json:"bindAreaID" check:"id" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdatePic 修改主图
func UpdatePic(args *ArgsUpdatePic) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_map_drawing_pic SET update_at = NOW(), name = :name, des = :des, file_id = :file_id, fix_height = :fix_height, fix_width = :fix_width, button_name = :button_name, bind_area_id = :bind_area_id, params = :params WHERE id = :id AND org_id = :org_id", args)
	return
}

// ArgsDeletePic 删除主图参数
type ArgsDeletePic struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeletePic 删除主图
func DeletePic(args *ArgsDeletePic) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "map_map_drawing_pic", "id = :id AND org_id = :org_id", args)
	return
}
