package BaseStyle

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//导入和导出设计

// ArgsOutput 导出样式库参数
type ArgsOutput struct {
	//仅选择样式库的mark列
	PickMarks pq.StringArray `json:"pickMarks"`
	//排除样式
	// 相关样式不会导出
	ExcludeMarks []string `json:"excludeMarks"`
}

// DataOutputComponent 导出样式库数据
type DataOutputComponent struct {
	//关联标识码
	// 必填
	// 页面内独特的代码片段，声明后可直接定义该组件的默认参数形式
	Mark string `db:"mark" json:"mark"`
	//组件名称
	Name string `db:"name" json:"name"`
	//组件介绍
	Des string `db:"des" json:"des"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type DataOutput struct {
	//样式库名称
	Name string `db:"name" json:"name"`
	//关联标识码
	// 用于识别代码片段
	Mark string `db:"mark" json:"mark"`
	//样式使用渠道
	// app APP；wxx 小程序等，可以任意定义，模块内不做限制
	SystemMark string `db:"system_mark" json:"systemMark"`
	//分栏样式结构设计
	Components []DataOutputComponent `db:"components" json:"components"`
	//默认标题
	// 标题是展示给用户的，样式库名称和该标题不是一个
	Title string `db:"title" json:"title"`
	//默认描述
	Des string `db:"des" json:"des"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Output 导出样式库
func Output(args *ArgsOutput) (data string, err error) {
	//获取所有样式
	var allStyle []FieldsStyle
	if len(args.PickMarks) > 0 {
		err = Router2SystemConfig.MainDB.Select(&allStyle, "SELECT id, create_at, update_at, delete_at, name, mark, system_mark, components, title, des, cover_file_id, des_files, sort_id, tags, params FROM core_style WHERE delete_at < to_timestamp(1000000) AND mark = ANY($1)", args.PickMarks)
		if err != nil || len(allStyle) < 1 {
			err = errors.New(fmt.Sprint("not find any style, ", err))
			return
		}
	} else {
		err = Router2SystemConfig.MainDB.Select(&allStyle, "SELECT id, create_at, update_at, delete_at, name, mark, system_mark, components, title, des, cover_file_id, des_files, sort_id, tags, params FROM core_style WHERE delete_at < to_timestamp(1000000)")
		if err != nil || len(allStyle) < 1 {
			err = errors.New(fmt.Sprint("not find any style, ", err))
			return
		}
	}
	//初始化集合
	var result []DataOutput
	//遍历样式库并组织数据
	for _, v := range allStyle {
		//跳过不导出的数据
		isFind := false
		for _, v2 := range args.ExcludeMarks {
			if v2 == v.Mark {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		//组织数据集合
		var components []DataOutputComponent
		//组装组件
		var componentList []FieldsComponent
		componentList, err = GetComponentMore(&ArgsGetComponentMore{
			IDs:        v.Components,
			HaveRemove: false,
		})
		if err != nil {
			err = nil
		}
		for _, v2 := range componentList {
			components = append(components, DataOutputComponent{
				Mark:   v2.Mark,
				Name:   v2.Name,
				Des:    v2.Des,
				SortID: v2.SortID,
				Tags:   v2.Tags,
				Params: v2.Params,
			})
		}
		//构建数据
		var vResult DataOutput
		vResult = DataOutput{
			Name:       v.Name,
			Mark:       v.Mark,
			SystemMark: v.SystemMark,
			Components: components,
			Title:      v.Title,
			Des:        v.Des,
			SortID:     v.SortID,
			Tags:       v.Tags,
			Params:     v.Params,
		}
		result = append(result, vResult)
	}
	//整合数据
	var dataByte []byte
	dataByte, err = json.Marshal(result)
	if err != nil {
		return
	}
	data = string(dataByte)
	//反馈数据
	return
}

// ArgsImport 导入数据参数
type ArgsImport struct {
	//数据
	Data string `json:"data"`
	//是否覆盖
	// 如果发现相同数据，是否覆盖
	NeedCover bool `json:"needCover"`
}

// Import 导入数据
func Import(args *ArgsImport) (err error) {
	//解析数据
	var result []DataOutput
	err = json.Unmarshal([]byte(args.Data), &result)
	if err != nil {
		return
	}
	//遍历数据导入
	for _, v := range result {
		//检查系统是否存在数据
		var vData FieldsStyle
		err = Router2SystemConfig.MainDB.Get(&vData, "SELECT id FROM core_style WHERE mark = $1", v.Mark)
		if err == nil && vData.ID > 0 {
			if !args.NeedCover {
				continue
			}
		}
		//写入组件
		var componentIDs []int64
		for _, v2 := range v.Components {
			var v2Data FieldsComponent
			err = Router2SystemConfig.MainDB.Get(&vData, "SELECT id FROM core_style_component WHERE mark = $1", v2.Mark)
			if err == nil && v2Data.ID > 0 {
				if !args.NeedCover {
					continue
				}
			}
			if v2Data.ID > 0 {
				err = UpdateComponent(&ArgsUpdateComponent{
					ID:          v2Data.ID,
					Mark:        v2.Mark,
					Name:        v2.Name,
					Des:         v2.Des,
					CoverFileID: 0,
					DesFiles:    []int64{},
					SortID:      v2.SortID,
					Tags:        v2.Tags,
					Params:      v2.Params,
				})
				if err != nil {
					return
				}
			} else {
				v2Data, err = CreateComponent(&ArgsCreateComponent{
					Mark:        v2.Mark,
					Name:        v2.Name,
					Des:         v2.Des,
					CoverFileID: 0,
					DesFiles:    []int64{},
					SortID:      v2.SortID,
					Tags:        v2.Tags,
					Params:      v2.Params,
				})
				if err != nil {
					return
				}
			}
			componentIDs = append(componentIDs, v2Data.ID)
		}
		//写入数据
		if vData.ID > 0 {
			err = UpdateStyle(&ArgsUpdateStyle{
				ID:          vData.ID,
				Name:        v.Name,
				Mark:        v.Mark,
				SystemMark:  v.SystemMark,
				Components:  componentIDs,
				Title:       v.Title,
				Des:         v.Des,
				CoverFileID: 0,
				DesFiles:    []int64{},
				SortID:      v.SortID,
				Tags:        v.Tags,
				Params:      v.Params,
			})
			if err != nil {
				return
			}
		} else {
			vData, err = CreateStyle(&ArgsCreateStyle{
				Name:        v.Name,
				Mark:        v.Mark,
				SystemMark:  v.SystemMark,
				Components:  componentIDs,
				Title:       v.Title,
				Des:         v.Des,
				CoverFileID: 0,
				DesFiles:    []int64{},
				SortID:      v.SortID,
				Tags:        v.Tags,
				Params:      v.Params,
			})
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
	return
}
