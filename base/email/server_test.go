package BaseEmail

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"testing"
)

var (
	serverGlobData FieldsEmailServerType
)

func TestInitServer(t *testing.T) {
	TestInit(t)
}

func TestCreateServer(t *testing.T) {
	newData, err := CreateServer(&ArgsCreateServer{
		OrgID:    0,
		Name:     "ffa-name",
		Host:     "tt.qq.com",
		Port:     "",
		IsSSL:    true,
		Email:    "abc@qq.com",
		Password: "abc",
		Params:   nil,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("create new data server, ", newData)
		t.Log("create new data server id: ", newData.ID)
		serverGlobData = newData
		serverID = serverGlobData.ID
	}
}

func TestGetServerByID(t *testing.T) {
	getData, err := GetServerByID(&ArgsGetServerByID{
		ID:    serverGlobData.ID,
		OrgID: -1,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data success, ", getData)
	}
}

func TestGetServerList(t *testing.T) {
	data, count, err := GetServerList(&ArgsGetServerList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:  -1,
		Search: "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get list, ", count, ", data:", data)
	}
}

func TestUpdateServer(t *testing.T) {
	err := UpdateServer(&ArgsUpdateServer{
		ID:       serverGlobData.ID,
		Host:     "tt3.qq.com",
		Port:     serverGlobData.Port,
		Email:    "abcd@qq.com",
		Password: "abc345",
		IsSSL:    serverGlobData.IsSSL,
		Name:     "f33",
	})
	if err != nil {
		t.Error(err)
	} else {
		getData, err := GetServerByID(&ArgsGetServerByID{
			ID: serverGlobData.ID,
		})
		if err != nil {
			t.Error(err)
		} else {
			t.Log("get data success and update success, ", getData)
		}
	}
}

func TestDeleteServerByID(t *testing.T) {
	if err := DeleteServerByID(&ArgsDeleteServerByID{
		ID:    serverGlobData.ID,
		OrgID: -1,
	}); err != nil {
		t.Error(err)
	} else {
		if _, err = CreateServer(&ArgsCreateServer{
			Host: "mail.qq.com", Email: "abc@qq.com", Password: "abc", IsSSL: true, Name: "fotomxq-name",
		}); err != nil {
			t.Error(err)
		}
	}
}

func TestClearServer(t *testing.T) {
	//清理
	if _, err := CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_email_server", "", nil); err != nil {
		t.Error(err)
	}
}
