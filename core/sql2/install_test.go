package CoreSQL2

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"testing"
	"time"
)

// TestClient_InstallSQL 示例，无法运行
// 因为本模块禁止应用实际数据库链接，如需测试请使用core外其他模块构建
func TestClient_InstallSQL(t *testing.T) {
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
		//Mark
		Mark string `db:"mark" json:"mark" check:"des" min:"1" max:"300"`
		//事件地址
		// nats - 触发的地址
		EventURL string `db:"event_url" json:"eventURL" check:"des" min:"1" max:"600"`
		//发生所属分公司ID
		CompanyID int64 `db:"company_id" json:"companyID" check:"id" index:"true"`
		//发生所属门店ID
		StoreID  int64 `db:"store_id" json:"storeID" check:"id" index:"true"`
		StoreID2 int64 `db:"store_id2" json:"storeID2" check:"id" index:"true"`
		//int: integer
		TestInt int `db:"test_int" json:"testInt" check:"id" index:"true"`
		//[]int: integer[]
		// 注意，实际开发中，不可能遇到该类型，因为postgresql不支持，请使用pq.Int32Array
		TestIntArr []int `db:"test_int_arr" json:"testIntArr" check:"id" index:"true"`
		//pq.Int32Array: integer[]
		TestIntArr2 pq.Int32Array `db:"test_int_arr2" json:"testIntArr2" check:"id"`
		//int64: bigint
		TestInt64 int64 `db:"test_int64" json:"testInt64" check:"id" index:"true"`
		//[]int64: bigint[]
		// 注意，实际开发中，不可能遇到该类型，因为postgresql不支持，请使用pq.Int64Array
		TestInt64Arr []int64 `db:"test_int64_arr" json:"testInt64Arr" check:"id"`
		//pq.Int64Array: bigint[]
		TestInt64Arr2 pq.Int64Array `db:"test_int64_arr2" json:"testInt64Arr2"`
		//bool: boolean
		TestBool bool `db:"test_bool" json:"testBool" check:"id" index:"true"`
		//time.Time: timestamp
		TestTime  time.Time `db:"test_time" json:"testTime" check:"id" default:"now()"`
		TestTime2 time.Time `db:"test_time2" json:"testTime2" check:"id" default:"0"`
		//string: varchar(max)
		TestString string `db:"test_string" json:"testString" check:"des" min:"1" max:"300"`
		//string: text
		TestString2 string `db:"test_string2" json:"testString2" check:"des" min:"1" max:"-1"`
		TestString3 string `db:"test_string3" json:"testString3" check:"des" min:"1" max:"-1" default:"text_default"`
	}
	tableName := "test_db_auto_install"
	//测试
	var testInstallDB Client
	testInstallDB.Init(&Router2SystemConfig.MainSQL, tableName)
	testInstallDB.StructData = &argsType{}
	err := testInstallDB.InstallSQL()
	if err != nil {
		t.Fatal(err)
	}
	_, _ = testInstallDB.DB.GetPostgresql().Exec("DROP TABLE" + " " + tableName)
}
