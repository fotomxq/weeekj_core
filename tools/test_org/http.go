package TestOrg

import (
	"encoding/json"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	TestAPI "github.com/fotomxq/weeekj_core/v5/tools/test_api"
	"testing"
)

//组织测试用的整套方案
// 本方案为url设计

// 创建新的组织
func Create(t *testing.T, userID string) {
	type paramsType struct {
		UserID     string `json:"userID" filter:"ID"`
		Name       string `json:"name" filter:"Des" filterMin:"1" filterMax:"600"`
		Des        string `json:"des" filter:"Des" filterMin:"1" filterMax:"1000"`
		ParentID   string `json:"parentID" filter:"ID" empty:"true"`
		WorkTimeID string `json:"workTimeID" filter:"ID" empty:"true"`
	}
	params := paramsType{
		UserID:     userID,
		Name:       "测试组织名称",
		Des:        "测试描述信息",
		ParentID:   "",
		WorkTimeID: "",
	}
	dataByte, err := TestAPI.Put("/v1/work/organizational/manager", params)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data OrgCore.FieldsOrg `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
		return
	}
	t.Log(data.Data, data.Count)
	OrgData = data.Data
}

// 让某个用户成为组织成员
func SetUserBind(t *testing.T, orgID string, userID string, name string, managers []string) {
	type paramsType struct {
		OrganizationalID string                           `json:"organizationalID" filter:"ID"`
		Managers         []string                         `json:"managers"`
		Name             string                           `json:"name" filter:"Name"`
		UserID           string                           `json:"userID" filter:"ID"`
		Params           []CoreSQLConfig.FieldsConfigType `json:"params"`
	}
	params := paramsType{
		OrganizationalID: orgID,
		Managers:         managers,
		Name:             name,
		UserID:           userID,
		Params:           nil,
	}
	dataByte, err := TestAPI.Post("/v1/work/organizational/manager/operate", params)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data interface{} `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
		return
	}
}

// 设置用户当前选择的组织
// url: /v1/work/organizational/user/select
func SetUserSelectOrg(t *testing.T, orgID string) {
	type paramsType struct {
		OrganizationalID string `json:"organizationalID" filter:"ID"`
	}
	params := paramsType{
		OrganizationalID: orgID,
	}
	dataByte, err := TestAPI.Post("/v1/work/organizational/user/select", params)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data interface{} `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
		return
	}
}

// 获取绑定列表
// url: /v1/work/organizational/manager/bind/list
func GetBindList(t *testing.T, orgID string) {
	type paramsType struct {
		Page             int64  `json:"page" filter:"Page"`
		Max              int64  `json:"max" filter:"Max"`
		Sort             string `json:"sort" filter:"Sort"`
		Desc             bool   `json:"desc" filter:"Desc"`
		OrganizationalID string `json:"organizationalID" filter:"ID"`
		Search           string `json:"search" filter:"Search" empty:"true"`
	}
	params := paramsType{
		Page:             1,
		Max:              10,
		Sort:             "_id",
		Desc:             false,
		OrganizationalID: orgID,
		Search:           "",
	}
	dataByte, err := TestAPI.Post("/v1/work/organizational/manager/bind/list", params)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data []OrgCore.FieldsBind `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
		return
	}
	t.Log(data.Data, data.Count)
	BindList = data.Data
}

// 删除组织
func DeleteByID(t *testing.T, id string) {
	dataByte, err := TestAPI.Delete("/v1/work/organizational/manager/id/"+id, nil)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data interface{} `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
		return
	}
}
