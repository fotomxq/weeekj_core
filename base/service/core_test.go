package BaseService

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
	Init()
}

// 测试数据库自动安装
// 请单独运行该单元测试!
func TestAutoInstallSQL(t *testing.T) {
	//初始化
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
	Init()
	//构建结构体
	type argsType struct {
		//ID
		ID int64 `db:"id" json:"id" check:"id" unique:"true"`
		//创建时间
		CreateAt time.Time `db:"create_at" json:"createAt"`
		//更新时间
		UpdateAt time.Time `db:"update_at" json:"updateAt"`
		//删除时间
		DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
		//事件编码
		Code string `db:"code" json:"code" check:"des" min:"1" max:"300" index:"true"`
		//事件地址
		// nats - 触发的地址
		EventURL string `db:"event_url" json:"eventURL" check:"des" min:"1" max:"600"`
		//发生所属分公司ID
		CompanyID int64 `db:"company_id" json:"companyID" check:"id" index:"true"`
		//发生所属门店ID
		StoreID  int64 `db:"store_id" json:"storeID" check:"id" index:"true"`
		StoreID2 int64 `db:"store_id2" json:"storeID2" check:"id" index:"true"`
	}
	tableName := "test_db_auto_install"
	//测试
	var testInstallDB CoreSQL2.Client
	testInstallDB.Init(&Router2SystemConfig.MainSQL, tableName)
	testInstallDB.StructData = &argsType{}
	err := testInstallDB.InstallSQL()
	if err != nil {
		t.Fatal(err)
	}
	_, _ = testInstallDB.DB.GetPostgresql().Exec("DROP TABLE" + " " + tableName)
}

func TestClear(t *testing.T) {
}
