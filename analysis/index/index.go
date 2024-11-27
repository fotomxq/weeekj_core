package AnalysisIndex

import (
	"errors"
	"fmt"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetIndexList 获取指标列表参数
type ArgsGetIndexList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否内置
	// 前端应拦截内置指标的删除操作，以免影响系统正常运行，启动重启后将自动恢复，所以删除操作是无法生效的
	NeedIsSystem bool `json:"needIsSystem"`
	IsSystem     bool `db:"is_system" json:"isSystem" index:"true" field_list:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetIndexList 获取指标列表
func GetIndexList(args *ArgsGetIndexList) (dataList []FieldsIndex, dataCount int64, err error) {
	//构建参数
	var conditionFields []BaseSQLTools.ArgsGetListSimpleConditionID
	if args.NeedIsSystem {
		conditionFields = append(conditionFields, BaseSQLTools.ArgsGetListSimpleConditionID{
			Name: "is_system",
			Val:  args.IsSystem,
		})
	}
	//获取数据
	dataCount, err = indexDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages:           args.Pages,
		ConditionFields: conditionFields,
		IsRemove:        args.IsRemove,
		Search:          args.Search,
	}, &dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		var vData FieldsIndex
		err = indexDB.GetInfo().GetInfoByID(v.ID, &vData)
		if err != nil || vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	//反馈
	return
}

// DataGetIndexAll 获取所有指标
type DataGetIndexAll struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//指标名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" index:"true" field_search:"true" field_list:"true"`
	//是否内置
	// 前端应拦截内置指标的删除操作，以免影响系统正常运行，启动重启后将自动恢复，所以删除操作是无法生效的
	IsSystem bool `db:"is_system" json:"isSystem" index:"true" field_list:"true"`
	//指标描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"-1" field_search:"true" field_list:"true" empty:"true"`
	//指标决策建议
	Decision string `db:"decision" json:"decision" check:"des" min:"1" max:"-1" empty:"true" field_search:"true"`
	//指标预警阈值
	// 0-100，归一化后的数据，超出此范围将可触发预警事件记录
	Threshold int64 `db:"threshold" json:"threshold" index:"true"`
	//是否启用
	// 关闭后将不对该指标进行汇总运算
	IsEnable bool `db:"is_enable" json:"isEnable" index:"true" field_list:"true"`
	//子指标
	SubIndex []DataGetIndexAll `json:"subIndex"`
}

// GetIndexAll 获取所有指标
func GetIndexAll() (dataList []DataGetIndexAll, err error) {
	//获取所有指标数据
	var indexRawList []FieldsIndex
	_, err = indexDB.GetList().GetAll(&BaseSQLTools.ArgsGetAll{
		ConditionFields: nil,
		IsRemove:        false,
	}, &indexRawList)
	if err != nil || len(indexRawList) < 1 {
		return
	}
	//获取所有关系数据
	var relationRawList []FieldsIndexRelation
	_, _ = indexRelationDB.GetList().GetAll(&BaseSQLTools.ArgsGetAll{
		ConditionFields: nil,
		IsRemove:        false,
	}, &relationRawList)
	//如果不存在关系，则重组数据后返回
	if len(relationRawList) < 1 {
		for k := 0; k < len(indexRawList); k++ {
			dataList = append(dataList, DataGetIndexAll{
				ID:          indexRawList[k].ID,
				Code:        indexRawList[k].Code,
				Name:        indexRawList[k].Name,
				IsSystem:    indexRawList[k].IsSystem,
				Description: indexRawList[k].Description,
				Decision:    indexRawList[k].Decision,
				Threshold:   indexRawList[k].Threshold,
				IsEnable:    indexRawList[k].IsEnable,
				SubIndex:    []DataGetIndexAll{},
			})
		}
		return
	}
	//如果存在关系
	// 建立第一级指标
	for _, vIndex := range indexRawList {
		isFind := false
		for _, vRel := range relationRawList {
			if vRel.RelationIndexID == vIndex.ID {
				isFind = true
				break
			}
		}
		if !isFind {
			dataList = append(dataList, DataGetIndexAll{
				ID:          vIndex.ID,
				Code:        vIndex.Code,
				Name:        vIndex.Name,
				IsSystem:    vIndex.IsSystem,
				Description: vIndex.Description,
				Decision:    vIndex.Decision,
				Threshold:   vIndex.Threshold,
				IsEnable:    vIndex.IsEnable,
				SubIndex:    []DataGetIndexAll{},
			})
		}
	}
	//检查指标环路
	if b := checkIndexAllRelChild(relationRawList); b {
		err = errors.New(fmt.Sprint("index relation loop"))
		return
	}
	//基于第一级指标，建立下级指标
	for k, vData := range dataList {
		dataList[k].SubIndex = getIndexAllRelChild(indexRawList, relationRawList, vData.ID)
	}
	//反馈
	return
}

// getIndexAllRelChild 无限递归构建子指标关系
func getIndexAllRelChild(indexRawList []FieldsIndex, relationRawList []FieldsIndexRelation, indexID int64) (result []DataGetIndexAll) {
	for _, vRel := range relationRawList {
		if vRel.IndexID != indexID {
			continue
		}
		//找到子指标，开始构建数据
		for _, vIndex := range indexRawList {
			if vIndex.ID != vRel.RelationIndexID {
				continue
			}
			result = append(result, DataGetIndexAll{
				ID:          vIndex.ID,
				Code:        vIndex.Code,
				Name:        vIndex.Name,
				IsSystem:    vIndex.IsSystem,
				Description: vIndex.Description,
				Decision:    vIndex.Decision,
				Threshold:   vIndex.Threshold,
				IsEnable:    vIndex.IsEnable,
				SubIndex:    getIndexAllRelChild(indexRawList, relationRawList, vIndex.ID),
			})
			break
		}
	}
	return
}

// checkIndexAllRelChild 指标关系环路检查
// 禁止子指标出现在上级指标中
func checkIndexAllRelChild(relationRawList []FieldsIndexRelation) (b bool) { // 创建一个map，用于存储每个IndexID的访问状态: 0-未访问, 1-正在访问, 2-已访问完
	visitStatus := make(map[int64]int64)
	// 递归函数，用于检查从某个IndexID开始的路径是否形成环路
	var checkLoop func(int64) bool
	checkLoop = func(indexID int64) bool {
		// 如果当前IndexID正在被访问，说明存在环路
		if visitStatus[indexID] == 1 {
			return true
		}
		// 如果当前IndexID已经被访问完，直接返回false
		if visitStatus[indexID] == 2 {
			return false
		}
		// 标记当前IndexID为正在访问
		visitStatus[indexID] = 1
		// 遍历所有关系，查找以当前IndexID为上级指标的关系
		for _, relation := range relationRawList {
			if relation.IndexID == indexID {
				// 递归检查下级指标
				if checkLoop(relation.RelationIndexID) {
					return true
				}
			}
		}
		// 标记当前IndexID为已访问完
		visitStatus[indexID] = 2
		return false
	}
	// 遍历所有关系，从每个未访问的IndexID开始检查
	for _, relation := range relationRawList {
		if visitStatus[relation.IndexID] == 0 && checkLoop(relation.IndexID) {
			return true
		}
	}
	return false
}

// DataGetIndexListByTop 获取指标列表顶部数据
type DataGetIndexListByTop struct {
	//指标ID
	IndexID int64 `db:"index_id" json:"indexID" check:"id" unique:"true"`
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//指标名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" index:"true" field_search:"true" field_list:"true"`
	//指标描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"100" field_search:"true" field_list:"true" empty:"true"`
	//指标决策建议
	Decision string `db:"decision" json:"decision" check:"des" min:"1" max:"-1" empty:"true" field_search:"true"`
}

// GetIndexListByTop 获取指标列表顶部
func GetIndexListByTop() (dataList []DataGetIndexListByTop, dataCount int64, err error) {
	//获取数据
	err = indexDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT i.id as index_id, i.code as code, max(i.name) as name, max(i.description) as description, max(i.decision) as decision FROM analysis_index as i, analysis_index_relation as r WHERE i.delete_at < to_timestamp(100000) and i.id = r.index_id GROUP BY i.id, i.code ORDER BY i.code;")
	if err != nil || len(dataList) < 1 {
		return
	}
	//反馈
	return
}

// GetIndexByCode 通过编码获取指标
func GetIndexByCode(code string) (data FieldsIndex, err error) {
	//获取数据
	err = indexDB.GetInfo().GetInfoByFields(map[string]any{
		"code": code,
	}, true, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// GetIndexByID 通过ID获取指标
func GetIndexByID(id int64) (data FieldsIndex, err error) {
	//获取数据
	err = indexDB.GetInfo().GetInfoByID(id, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// GetIndexNameByID 获取指标名称
func GetIndexNameByID(id int64) (name string) {
	data, _ := GetIndexByID(id)
	return data.Name
}

// GetIndexNameByCode 获取指标名称
func GetIndexNameByCode(code string) (name string) {
	data, _ := GetIndexByCode(code)
	return data.Name
}

// ArgsCreateIndex 创建新的指标参数
type ArgsCreateIndex struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//指标名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" index:"true" field_search:"true" field_list:"true"`
	//是否内置
	// 前端应拦截内置指标的删除操作，以免影响系统正常运行，启动重启后将自动恢复，所以删除操作是无法生效的
	// 通过接口应强制给予false
	IsSystem bool `db:"is_system" json:"isSystem" index:"true"`
	//指标描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"100" field_search:"true" field_list:"true" empty:"true"`
	//指标决策建议
	Decision string `db:"decision" json:"decision" check:"des" min:"1" max:"-1" empty:"true" field_search:"true"`
	//指标预警阈值
	// 0-100，归一化后的数据，超出此范围将可触发预警事件记录
	Threshold int64 `db:"threshold" json:"threshold" index:"true"`
	//是否启用
	// 关闭后将不对该指标进行汇总运算
	// 通过接口应强制给予false
	IsEnable bool `db:"is_enable" json:"isEnable" index:"true"`
}

// CreateIndex 创建新的指标
func CreateIndex(args *ArgsCreateIndex) (err error) {
	//检查code是否重复
	if b, _ := indexDB.GetInfo().CheckInfoByFields(map[string]any{
		"code": args.Code,
	}, true); b {
		err = errors.New("code is exist")
		return
	}
	//修改启动设置
	args.IsEnable = false
	//创建
	_, err = indexDB.GetInsert().InsertRow(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateIndex 更新指标定义参数
type ArgsUpdateIndex struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//指标名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" index:"true" field_search:"true" field_list:"true"`
	//指标描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"100" field_search:"true" field_list:"true" empty:"true"`
	//指标决策建议
	Decision string `db:"decision" json:"decision" check:"des" min:"1" max:"-1" empty:"true" field_search:"true"`
	//指标预警阈值
	// 0-100，归一化后的数据，超出此范围将可触发预警事件记录
	Threshold int64 `db:"threshold" json:"threshold" index:"true"`
	//是否启用
	// 关闭后将不对该指标进行汇总运算
	IsEnable bool `db:"is_enable" json:"isEnable" index:"true"`
}

// UpdateIndex 更新指标定义
func UpdateIndex(args *ArgsUpdateIndex) (err error) {
	//修改数据
	err = indexDB.GetUpdate().UpdateByID(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// DeleteIndex 删除指标
func DeleteIndex(id int64) (err error) {
	//删除数据
	err = indexDB.GetDelete().DeleteByID(id)
	if err != nil {
		return
	}
	//删除指标的参数
	_ = indexParamDB.GetDelete().DeleteByField("index_id", id)
	//删除指标的关系
	_ = indexRelationDB.GetDelete().DeleteByField("index_id", id)
	_ = indexRelationDB.GetDelete().DeleteByField("relation_index_id", id)
	//反馈
	return
}
