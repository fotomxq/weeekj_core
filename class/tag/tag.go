package ClassTag

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//快速标签模块
// 注意，该模块需自行指定数据表，否则无法正常使用
// 所有方法集合采用类结构实现，可任意声明多个使用

// Tag 对象结构
type Tag struct {
	//标签主表名称
	TagTableName string
}

// ArgsGetList 查询列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//绑定ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//搜索标签
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 查询列表
func (t *Tag) GetList(args *ArgsGetList) (dataList []FieldsTag, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.BindID > 0 {
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	var rawList []FieldsTag
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		t.TagTableName,
		"id",
		"SELECT id "+"FROM "+t.TagTableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := t.getByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetIDs 查询一组ID数据参数
type ArgsGetIDs struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//绑定ID
	// 用于验证
	BindID int64 `json:"bindID"`
	//反馈限制
	Limit int `json:"limit"`
}

// GetByIDs 查询一组ID数据
func (t *Tag) GetByIDs(args *ArgsGetIDs) (dataList []FieldsTag, err error) {
	var rawList []FieldsTag
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id "+"FROM "+t.TagTableName+" WHERE id = ANY($1) AND ($2 < 1 OR bind_id = $2) LIMIT $3;", args.IDs, args.BindID, args.Limit)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := t.getByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetByIDNames 获取一组数据map
func (t *Tag) GetByIDNames(ids pq.Int64Array, bindID int64, limit int) (dataList map[int64]string) {
	if len(ids) < 1 {
		return
	}
	var rawData []FieldsTag
	var err error
	err = Router2SystemConfig.MainDB.Select(&rawData, "SELECT id "+"FROM "+t.TagTableName+" WHERE id = ANY($1) AND ($2 < 1 OR bind_id = $2) LIMIT $3;", ids, bindID, limit)
	if err != nil {
		return
	}
	if len(rawData) < 1 {
		err = errors.New("no data")
		return
	}
	dataList = map[int64]string{}
	for k := 0; k < len(rawData); k++ {
		vData := t.getByID(rawData[k].ID)
		if vData.ID < 1 {
			continue
		}
		dataList[vData.ID] = vData.Name
	}
	return
}

// ArgsCheckBind 检查一组标签是否数据该绑定？参数
type ArgsCheckBind struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//绑定D
	BindID int64 `json:"bindID" check:"id"`
}

// CheckBind 检查一组标签是否数据该绑定？
func (t *Tag) CheckBind(args *ArgsCheckBind) (err error) {
	type dataType struct {
		//基础
		ID int64 `db:"id" json:"id"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id "+"FROM "+t.TagTableName+" WHERE id = ANY($1) AND bind_id != :bind_id;", args.IDs, args.BindID)
	if err != nil {
		err = nil
		return
	}
	if len(dataList) > 0 {
		err = errors.New("have not this bind id")
		return
	}
	return
}

// ArgsCreate 创建标签参数
type ArgsCreate struct {
	//绑定ID
	BindID int64 `json:"bindID" check:"id"`
	//名称
	Name string `json:"name" check:"name"`
}

// Create 创建标签
func (t *Tag) Create(args *ArgsCreate) (data FieldsTag, err error) {
	var dataID int64
	dataID, err = t.getByName(args.BindID, args.Name)
	if err == nil && dataID > 0 {
		err = errors.New("name replace")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, t.TagTableName, "INSERT "+"INTO "+t.TagTableName+"(bind_id, name) VALUES(:bind_id, :name)", map[string]interface{}{
		"bind_id": args.BindID,
		"name":    args.Name,
	}, &data)
	return
}

// ArgsUpdateByID 修改标签参数
type ArgsUpdateByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//绑定ID
	// 用于验证
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//名称
	Name string `json:"name" check:"name"`
}

// UpdateByID 修改标签
func (t *Tag) UpdateByID(args *ArgsUpdateByID) (err error) {
	var dataID int64
	dataID, err = t.getByName(args.BindID, args.Name)
	if err == nil {
		if dataID != args.ID {
			err = errors.New("name replace")
			return
		}
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, fmt.Sprint("UPDATE ", t.TagTableName, " SET update_at = NOW(), name = :name WHERE id = :id AND (:bind_id < 1 OR bind_id = :bind_id)"), map[string]interface{}{
		"name":    args.Name,
		"id":      args.ID,
		"bind_id": args.BindID,
	})
	if err != nil {
		return
	}
	t.deleteTagCache(args.ID)
	return
}

// ArgsDeleteByID 删除标签参数
type ArgsDeleteByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//绑定ID
	// 用于验证
	BindID int64 `json:"bindID" check:"id" empty:"true"`
}

// DeleteByID 删除标签
func (t *Tag) DeleteByID(args *ArgsDeleteByID) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, t.TagTableName, "id = :id AND (:bind_id < 1 OR bind_id = :bind_id)", map[string]interface{}{
		"id":      args.ID,
		"bind_id": args.BindID,
	})
	if err != nil {
		return
	}
	t.deleteTagCache(args.ID)
	return
}

// getByID 获取标签
func (t *Tag) getByID(id int64) (data FieldsTag) {
	cacheMark := t.getTagCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, bind_id, name "+"FROM "+t.TagTableName+" WHERE id = $1;", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 259200)
	return
}

func (t *Tag) getByName(bindID int64, name string) (dataID int64, err error) {
	err = Router2SystemConfig.MainDB.Get(&dataID, "SELECT id "+"FROM "+t.TagTableName+" WHERE bind_id = $1 AND name = $2;", bindID, name)
	return
}
