package ClassQueue

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

//通用消息列队处理件
// 外部模块可声明本通用方法，作为列队持久化存储数据
// 完成数据后，将其压出即可

type Queue struct {
	//主表
	TableName string
	//清理多久的数据？
	ClearDay int
}

// Init 初始化对象
func (t *Queue) Init(tagTableName string) {
	t.TableName = tagTableName
}

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//处理状态
	// 如果消息件存在多个状态，可使用，否则应及时删除该消息
	Status int `db:"status" json:"status" check:"than0Int" empty:"true"`
}

// GetList 获取列表
func (t *Queue) GetList(args *ArgsGetList) (dataList []FieldsQueue, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.Status > -1 {
		where = where + "status = :status"
		maps["status"] = args.Status
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		t.TableName,
		"id",
		"SELECT id, create_at, update_at, mod_id, status, params FROM "+t.TableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at"},
	)
	return
}

// ArgsGetByID 获取ID参数
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetByID 获取ID
func (t *Queue) GetByID(args *ArgsGetByID) (data FieldsQueue, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, mod_id, status, params FROM "+t.TableName+" WHERE id = $1", args.ID)
	return
}

// ArgsGetByModID 通过绑定ID获取参数
type ArgsGetByModID struct {
	//Mod ID
	ModID int64 `db:"mod_id" json:"modID" check:"id"`
}

// GetByModID 通过绑定ID获取
func (t *Queue) GetByModID(args *ArgsGetByModID) (data FieldsQueue, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, mod_id, status, params FROM "+t.TableName+" WHERE mod_id = $1", args.ModID)
	return
}

// Pick 提取数据
// 自动按照最早提取原则提取数据
// 提取等同销毁处理
func (t *Queue) Pick() (data FieldsQueue, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, mod_id, status, params FROM "+t.TableName+" ORDER BY id ASC LIMIT 1")
	if err == nil {
		err = t.Delete(&ArgsDelete{
			ID: data.ID,
		})
		return
	}
	return
}

// ArgsAppend 写入数据参数
type ArgsAppend struct {
	//其他模块的ID
	ModID int64 `db:"mod_id" json:"modID"`
	//处理状态
	// 如果消息件存在多个状态，可使用，否则应及时删除该消息
	Status int `db:"status" json:"status"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Append 写入数据
func (t *Queue) Append(args *ArgsAppend) (err error) {
	_, err = t.GetByModID(&ArgsGetByModID{
		ModID: args.ModID,
	})
	if err == nil {
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO "+t.TableName+"(mod_id, status, params) VALUES(:mod_id, :status, :params)", args)
	if err != nil {
		return
	}
	//主动清理数据
	t.clearByDay()
	//反馈
	return
}

// ArgsUpdateStatus 修改status参数
type ArgsUpdateStatus struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//处理状态
	// 如果消息件存在多个状态，可使用，否则应及时删除该消息
	Status int `db:"status" json:"status"`
}

// UpdateStatus 修改status
func (t *Queue) UpdateStatus(args *ArgsUpdateStatus) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE "+t.TableName+" SET update_at = NOW(), status = :status WHERE id = :id", args)
	return
}

// ArgsDelete 删除数据参数
type ArgsDelete struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// Delete 删除数据
func (t *Queue) Delete(args *ArgsDelete) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, t.TableName, "id", args)
	return
}

// ClearByDay 清理超过N天数据
func (t *Queue) clearByDay() {
	if t.ClearDay < 1 {
		t.ClearDay = 365
	}
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, t.TableName, "create_at <= :create_at", map[string]interface{}{
		"create_at": CoreFilter.GetNowTimeCarbon().SubDays(t.ClearDay).Time,
	})
	return
}
