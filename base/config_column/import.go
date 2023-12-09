package BaseConfigColumn

import (
	"encoding/json"
	"errors"
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//导入和导出设计
// 仅支持系统层的数据处理

type ArgsOutput struct {
	//仅选择样式库的mark列
	PickMarks pq.StringArray `json:"pickMarks"`
	//排除样式
	// 相关样式不会导出
	ExcludeMarks []string `json:"excludeMarks"`
}

type DataOutput struct {
	//标识码
	// 在来源系统内，该数据必须唯一，前端可识别具体是哪个页面的哪个组件
	// 不同层级可声明一个标识码，系统将反馈最大层级的数据
	Mark string `db:"mark" json:"mark"`
	//保存数据集
	// 前后顺序将按照该顺序一致
	Data FieldsChildList `db:"data" json:"data"`
}

func Output(args *ArgsOutput) (data string, err error) {
	//获取所有样式
	var allColumn []FieldsColumn
	if len(args.PickMarks) > 0 {
		err = Router2SystemConfig.MainDB.Select(&allColumn, "SELECT id, create_at, update_at, system, bind_id, mark, data FROM core_config_column WHERE system = 0 AND bind_id = 0 AND mark = ANY($1)", args.PickMarks)
		if err != nil || len(allColumn) < 1 {
			err = errors.New(fmt.Sprint("not find any column, ", err))
			return
		}
	} else {
		err = Router2SystemConfig.MainDB.Select(&allColumn, "SELECT id, create_at, update_at, system, bind_id, mark, data FROM core_config_column WHERE system = 0 AND bind_id = 0")
		if err != nil || len(allColumn) < 1 {
			err = errors.New(fmt.Sprint("not find any column, ", err))
			return
		}
	}
	//数据集合
	var result []DataOutput
	//遍历样式
	for _, v := range allColumn {
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
		//写入数据
		vResult := DataOutput{
			Mark: v.Mark,
			Data: v.Data,
		}
		result = append(result, vResult)
	}
	//反馈数据
	var dataByte []byte
	dataByte, err = json.Marshal(result)
	if err != nil {
		return
	}
	data = string(dataByte)
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
		var vData FieldsColumn
		err = Router2SystemConfig.MainDB.Get(&vData, "SELECT id FROM core_config_column WHERE mark = $1", v.Mark)
		if err == nil && vData.ID > 0 {
			if !args.NeedCover {
				continue
			}
		}
		//写入数据
		err = Set(&ArgsSet{
			System: 0,
			Mark:   v.Mark,
			BindID: 0,
			Data:   v.Data,
		})
		if err != nil {
			return
		}
	}
	return
}
