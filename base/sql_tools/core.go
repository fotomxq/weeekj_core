package BaseSQLTools

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

//Quick 快速组装模块
/**
1. 支持快速对定义好的结构体，自动生成通用方法
2. 方法包括：查看列表、查看详情、创建、更新、删除
*/
type Quick struct {
	//sql操作核心
	client CoreSQL2.Client
	//缓冲前缀
	prefixCacheMark string
	//启动delete软删除
	openSoftDelete bool
}

// FieldsTable 结构体示例(不要使用此表)
type eqFieldsTable struct {
	//ID
	// unique:"true" 启动唯一索引支持，不允许多次出现
	ID int64 `db:"id" json:"id" unique:"true"`
	//创建时间
	// default:"now()" 设置默认值
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	// default:"now()" 设置默认值
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	//  default:"0" 设置默认值
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0"`
	//名称
	// check:"des" min:"1" max:"300" 参照字段检查方法
	// field_search 启动列表搜索支持
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" field_search:"true"`
	//上级ID
	// check:"id" 参照字段检查方法
	// index 启动索引支持
	// field_list 支持列表模式检索
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" index:"true" field_list:"true"`
}

// Init 初始化模块
func (c *Quick) Init(tableName string, structData any) (err error) {
	//初始化
	_, err = c.client.Init2(&Router2SystemConfig.MainSQL, tableName, structData)
	if err != nil {
		err = errors.New(fmt.Sprint("init sql client failed: ", err))
		return
	}
	//设置前缀
	c.prefixCacheMark = fmt.Sprint("base:sql.tools.quick:", tableName, ":")
	//获取是否存在delete_at字段，则标记软删除启动
	for _, v := range c.client.GetFields() {
		if v == "delete_at" {
			c.openSoftDelete = true
			break
		}
	}
	//反馈
	return
}

func (c *Quick) GetClient() *CoreSQL2.Client {
	return &c.client
}

// GetInfo 获取Info
func (c *Quick) GetInfo() *QuickInfo {
	return &QuickInfo{
		quickClient: c,
	}
}

// GetList 获取list
func (c *Quick) GetList() *QuickList {
	return &QuickList{
		quickClient: c,
	}
}

// GetInsert 获取insert
func (c *Quick) GetInsert() *QuickInsert {
	return &QuickInsert{
		quickClient: c,
	}
}

// GetUpdate 获取update
func (c *Quick) GetUpdate() *QuickUpdate {
	return &QuickUpdate{
		quickClient: c,
	}
}

// GetDelete 获取delete
func (c *Quick) GetDelete() *QuickDelete {
	return &QuickDelete{
		quickClient: c,
	}
}

// GetAnalysis 获取analysis
func (c *Quick) GetAnalysis() *QuickAnalysis {
	return &QuickAnalysis{
		quickClient: c,
	}
}
