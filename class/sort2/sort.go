package ClassSort2

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetList 查询列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//绑定ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//标识码
	Mark string `json:"mark" check:"mark" empty:"true"`
	//上级ID
	ParentID int64 `json:"parentID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 查询列表
func (t *Sort) GetList(args *ArgsGetList) (dataList []FieldsSort, dataCount int64, err error) {
	var rawList []FieldsSort
	dataCount, err = t.sortDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "sort", "name"}).SetIDQuery("bind_id", args.BindID).SetStringQuery("mark", args.Mark).SetIDQuery("parent_id", args.ParentID).SetSearchQuery([]string{"name", "des"}, args.Search).SetPages(args.Pages).ResultAndCount(&rawList)
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

// GetByName 通过名称获取分类
func (t *Sort) GetByName(bindID int64, name string) (data FieldsSort, err error) {
	err = t.sortDB.Get().SetIDQuery("bind_id", bindID).SetStringQuery("name", name).Result(&data)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = t.getByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// GetAll 获取组织下所有分类
func (t *Sort) GetAll(bindID int64, parentID int64) (dataList []FieldsSort, err error) {
	var rawList []FieldsSort
	err = t.sortDB.Select().SetFieldsList([]string{"id"}).SetDeleteQuery("delete_at", false).SetIDQuery("bind_id", bindID).SetIDQuery("parent_id", parentID).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  999,
		Sort: "id",
		Desc: false,
	}).Result(&dataList)
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

// ArgsGetByID 查询指定ID参数
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//绑定ID
	// 用于验证
	BindID int64 `json:"bindID" check:"id" empty:"true"`
}

// GetByID 查询指定ID
func (t *Sort) GetByID(args *ArgsGetByID) (data FieldsSort, err error) {
	data = t.getByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.BindID, data.BindID) {
		err = errors.New("no data")
		return
	}
	return
}

func (t *Sort) GetByIDNoErr(id int64, bindID int64) (data FieldsSort) {
	if id < 1 {
		return
	}
	data = t.getByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(bindID, data.BindID) {
		data = FieldsSort{}
		return
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
func (t *Sort) GetByIDs(args *ArgsGetIDs) (dataList []FieldsSort, err error) {
	var rawList []FieldsSort
	err = t.sortDB.Select().SetFieldsList([]string{"id"}).SetDeleteQuery("delete_at", false).SetIDsQuery("id", args.IDs).SetIDQuery("bind_id", args.BindID).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  int64(args.Limit),
		Sort: "id",
		Desc: false,
	}).Result(&dataList)
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

// ArgsGetID 获取名称参数
type ArgsGetID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//绑定ID
	// 用于验证
	BindID int64 `json:"bindID"`
}

// GetName 获取名称
func (t *Sort) GetName(args *ArgsGetID) (data string, err error) {
	rawData := t.getByID(args.ID)
	if rawData.ID < 1 || !CoreFilter.EqID2(args.BindID, rawData.BindID) {
		err = errors.New("no data")
		return
	}
	data = rawData.Name
	return
}

// GetAllCount 获取有多少分类
func (t *Sort) GetAllCount() (count int64) {
	count = t.sortDB.Analysis().Count("")
	return
}

func (t *Sort) GetNameNoErr(id int64) (data string) {
	if id < 1 {
		return
	}
	rawData := t.getByID(id)
	if rawData.ID < 1 {
		return
	}
	data = rawData.Name
	return
}

func (t *Sort) GetNameMoreNoErr(args *ArgsGetIDs) (data map[int64]string) {
	data = map[int64]string{}
	for _, v := range args.IDs {
		if len(data) > args.Limit {
			break
		}
		data[v] = t.GetNameNoErr(v)
	}
	return
}

// ArgsGetParams 获取扩展参数
type ArgsGetParams struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//来源结构
	// 用于验证
	BindID int64 `json:"bindID" check:"id"`
}

// GetParams 获取扩展
func (t *Sort) GetParams(args *ArgsGetParams) (paramsData CoreSQLConfig.FieldsConfigsType, err error) {
	//data := t.getByID(args.ID)
	//if data.ID < 1 || !CoreFilter.EqID2(args.BindID, data.BindID) {
	//	err = errors.New("no data")
	//	return
	//}
	//paramsData = data.Params
	return
}

// ArgsGetParam 获取指定的扩展参数
type ArgsGetParam struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//来源结构
	// 用于验证
	BindID int64 `json:"bindID" check:"id"`
	//标识码
	Mark string `json:"mark" check:"mark"`
}

// GetParam 获取扩展
func (t *Sort) GetParam(args *ArgsGetParam) (val string, err error) {
	//data := t.getByID(args.ID)
	//if data.ID < 1 || !CoreFilter.EqID2(args.BindID, data.BindID) {
	//	err = errors.New("no data")
	//	return
	//}
	//for _, v := range data.Params {
	//	if v.Mark == args.Mark {
	//		val = v.Val
	//		return
	//	}
	//}
	//err = errors.New("mark not find")
	return
}

// GetParamInt64 获取扩展同时转int64
func (t *Sort) GetParamInt64(args *ArgsGetParam) (val int64, err error) {
	//data := t.getByID(args.ID)
	//if data.ID < 1 || !CoreFilter.EqID2(args.BindID, data.BindID) {
	//	err = errors.New("no data")
	//	return
	//}
	//for _, v := range data.Params {
	//	if v.Mark == args.Mark {
	//		val, err = CoreFilter.GetInt64ByString(v.Val)
	//		return
	//	}
	//}
	//err = errors.New("mark not find")
	return
}

// ArgsCheckBind 检查一组分类是否数据该绑定参数
type ArgsCheckBind struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//绑定D
	BindID int64 `json:"bindID" check:"id"`
}

// CheckBind 检查一组分类是否数据该绑定
func (t *Sort) CheckBind(args *ArgsCheckBind) (err error) {
	type dataType struct {
		//基础
		ID int64 `db:"id" json:"id"`
	}
	var dataList []dataType
	err = t.sortDB.Select().SetFieldsList([]string{"id"}).SetDeleteQuery("delete_at", false).SetIDsQuery("id", args.IDs).SetIDQuery("bind_id", args.BindID).Result(&dataList)
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

// ArgsCreate 创建分组参数
type ArgsCreate struct {
	//来源结构
	BindID int64 `json:"bindID" check:"id"`
	//分组标识码
	Mark string `db:"mark" json:"mark"  check:"mark" empty:"true"`
	//上级ID
	ParentID int64 `json:"parentID" check:"id" empty:"true"`
	//封面图
	CoverFileID int64 `json:"coverFileID" check:"id" empty:"true"`
	//介绍图文
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//名称
	Name string `json:"name" check:"name"`
	//描述
	Des string `json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

// Create 创建分组
func (t *Sort) Create(args *ArgsCreate) (data FieldsSort, err error) {
	if args.DesFiles == nil {
		args.DesFiles = pq.Int64Array{}
	}
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	if args.ParentID > 0 {
		data = t.getByID(args.ParentID)
		if data.ID < 1 || !CoreFilter.EqID2(args.BindID, data.BindID) {
			err = errors.New(fmt.Sprint("parent id not exist, ", err, ", parent id: ", args.ParentID, ", bind id: ", args.BindID))
			return
		}
	}
	var sort int
	var dataList []FieldsSort
	dataList, _, err = t.GetList(&ArgsGetList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  1,
			Sort: "sort",
			Desc: true,
		},
		BindID:   args.BindID,
		Mark:     "",
		ParentID: args.ParentID,
		Search:   "",
	})
	if err == nil && len(dataList) > 0 {
		sort = dataList[0].Sort
	} else {
		sort = 1
	}
	var newID int64
	newID, err = t.sortDB.Insert().SetFields([]string{"bind_id", "mark", "parent_id", "sort", "cover_file_id", "des_files", "name", "des"}).Add(map[string]any{
		"bind_id":       args.BindID,
		"mark":          args.Mark,
		"parent_id":     args.ParentID,
		"sort":          sort + 1,
		"cover_file_id": args.CoverFileID,
		"des_files":     args.DesFiles,
		"name":          args.Name,
		"des":           args.Des,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	data = t.getByID(newID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	err = t.UpdateParamsAdd(&ArgsUpdateParams{
		ID:     data.ID,
		BindID: data.BindID,
		Params: args.Params,
	})
	if err != nil {
		return
	}
	return
}

// ArgsUpdateByID 修改分组参数
type ArgsUpdateByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//来源结构
	// 用于验证
	BindID int64 `json:"bindID" check:"id"`
	//分组标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//上级ID
	ParentID int64 `json:"parentID" check:"id" empty:"true"`
	//排序
	Sort int `json:"sort"`
	//封面图
	CoverFileID int64 `json:"coverFileID" check:"id" empty:"true"`
	//介绍图文
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//名称
	Name string `json:"name" check:"name"`
	//描述
	Des string `json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

// UpdateByID 修改分组
func (t *Sort) UpdateByID(args *ArgsUpdateByID) (err error) {
	data := t.getByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.BindID, data.BindID) {
		return errors.New("no data")
	}
	if data.ParentID != args.ParentID {
		err = t.checkParent(data.BindID, args.ParentID, args.ID, []int64{})
		if err != nil {
			err = errors.New("parent is cycle or not exist, " + err.Error())
			return
		}
	}
	err = t.sortDB.Update().NeedSoft(true).AddWhereID(args.ID).SetFields([]string{"mark", "parent_id", "sort", "cover_file_id", "des_files", "name", "des"}).NamedExec(map[string]any{
		"mark":          args.Mark,
		"parent_id":     args.ParentID,
		"sort":          args.Sort,
		"cover_file_id": args.CoverFileID,
		"des_files":     args.DesFiles,
		"name":          args.Name,
		"des":           args.Des,
	})
	if err != nil {
		return
	}
	t.deleteSortCache(args.ID)
	return t.UpdateParams(&ArgsUpdateParams{
		ID:     data.ID,
		BindID: data.BindID,
		Params: args.Params,
	})
}

// ArgsUpdateParams 修改扩展参数参数
type ArgsUpdateParams struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//来源结构
	// 用于验证
	BindID int64 `json:"bindID" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

// UpdateParams 修改扩展参数参数
func (t *Sort) UpdateParams(args *ArgsUpdateParams) (err error) {
	//_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE "+t.SortTableName+" SET params = :params WHERE id = :id AND (:bind_id < 1 OR bind_id = :bind_id)", map[string]interface{}{
	//	"id":      args.ID,
	//	"bind_id": args.BindID,
	//	"params":  args.Params,
	//})
	//if err != nil {
	//	return
	//}
	//t.deleteSortCache(args.ID)
	return
}

// UpdateParamsAdd 增量修改扩展参数
func (t *Sort) UpdateParamsAdd(args *ArgsUpdateParams) (err error) {
	//data := t.getByID(args.ID)
	//if data.ID < 1 || !CoreFilter.EqID2(args.BindID, data.BindID) {
	//	return errors.New("no data")
	//}
	//for _, v := range args.Params {
	//	isFind := false
	//	for k2, v2 := range data.Params {
	//		if v.Mark == v2.Mark {
	//			isFind = true
	//			data.Params[k2].Val = v.Val
	//			break
	//		}
	//	}
	//	if !isFind {
	//		data.Params = append(data.Params, v)
	//	}
	//}
	//_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE "+t.SortTableName+" SET params = :params WHERE id = :id AND (:bind_id < 1 OR bind_id = :bind_id)", map[string]interface{}{
	//	"id":      args.ID,
	//	"bind_id": args.BindID,
	//	"params":  data.Params,
	//})
	//if err != nil {
	//	return
	//}
	//t.deleteSortCache(args.ID)
	return
}

// ArgsDeleteByID 删除分组参数
type ArgsDeleteByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//来源结构
	// 用于验证
	BindID int64 `json:"bindID" check:"id" empty:"true"`
}

// DeleteByID 删除分组参数
func (t *Sort) DeleteByID(args *ArgsDeleteByID) (err error) {
	err = t.sortDB.Delete().NeedSoft(true).AddWhereID(args.ID).SetWhereOrThan("bind_id", args.BindID).ExecNamed(nil)
	if err != nil {
		return
	}
	t.deleteSortCache(args.ID)
	return
}

// checkParent 检查上级是否形成循环
// 将检查自己的上级，直到到达头部
// params fromID string 来源ID，每个上级必须一致，否则异常
// params parentID string 上级ID
// params id string 要监测的ID，该值提交后恒定
func (t *Sort) checkParent(bindID int64, parentID int64, id int64, checkList []int64) error {
	//如果上级为空，则正常跳出
	if parentID < 1 {
		return nil
	}
	//如果ID和上级ID相同，则异常
	if id == parentID {
		return errors.New("parent cycle")
	}
	//检查该数据是否已经存在于列队？
	for _, v := range checkList {
		if v == parentID {
			return errors.New("parent cycle")
		}
	}
	//获取上级信息结构
	data := t.getByID(parentID)
	if data.ID < 1 || !CoreFilter.EqID2(bindID, data.BindID) {
		return errors.New("no data")
	}
	//写入检查列队
	checkList = append(checkList, data.ID)
	//继续监测
	return t.checkParent(bindID, data.ParentID, id, checkList)
}

// getByID 获取分组
func (t *Sort) getByID(id int64) (data FieldsSort) {
	cacheMark := t.getSortCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = t.sortDB.Get().SetDefaultFields().GetByID(id).Result(&data)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 259200)
	return
}
