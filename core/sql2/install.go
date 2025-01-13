package CoreSQL2

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"reflect"
	"strings"
)

// InstallSQL SQL自动安装工具
/**
dataDefault 初始化需采用的空数据集
1. 如果给予值，则代表有默认值，将按照默认值构建数据表
2. eg: dataDefault = data = ClassSort.FieldsSort{}
3. 如果已经创建过的表，会自动补全数据；但注意，旧的字段需手动变更调整！补增字段，不支持id。

识别规则(tag) :
db: 数据库字段名
index=true: 是否建立索引，用于提高检索效率
unique=true: 是否全局唯一

值类型: 数据库字段类型
max="1/-1": 最大长度，数字或者*代表最大长度，如果给予-1，则例如string按照text处理
index_out="tableName:field_name": 外键索引
default="any/now()": 默认值，sql直接写入；预设值不支持。

值类型转化对应关系:
int: integer
[]int: integer[]
pq.Int32Array: integer[]
int64: bigint
[]int64: bigint[]
pq.Int64Array: bigint[]
bool: boolean
time.Time: timestamp
string: varchar(max), 如果没有给予max，将按照255长度
string: text

预设值规则(tag) :
id: 主键，和unique=true一样
createAt: 创建时间
updateAt: 更新时间
deleteAt: 删除时间
code: 编码
mark: 标识码
comment: 评价
*/
func (t *Client) InstallSQL() (err error) {
	//构建SQL
	sqlData := "CREATE TABLE IF NOT EXISTS " + t.TableName + " ("
	var appendFields []string
	t.installAppendSQLData = []string{}
	t.installNunIndexKeyNum = 0
	//当前表是否已经运行
	// 获取列名称
	var columnNames []string
	err = t.DB.GetPostgresql().Select(&columnNames, "select column_name"+" from information_schema.columns where table_schema = 'public' and table_name = '"+t.TableName+"';")
	if len(columnNames) > 0 {
		sqlData = ""
	} else {
		err = nil
	}
	//计划增加的字段
	var needAddFields []string
	//初始化全量字段
	t.fieldNameList = []clientField{}
	//获取结构体
	paramsType := reflect.TypeOf(t.StructData).Elem()
	step := 0
	for step < paramsType.NumField() {
		//获取当前节点值
		vField := paramsType.Field(step)
		//下一步
		step += 1
		//获取当前节点类型
		vType := vField.Type.String()
		//获取db值
		dbVal := vField.Tag.Get("db")
		needAddFields = append(needAddFields, dbVal)
		//是否索引
		index := vField.Tag.Get("index") == "true"
		//是否唯一索引
		unique := vField.Tag.Get("unique") == "true"
		//最大长度
		maxVal := CoreFilter.GetIntByStringNoErr(vField.Tag.Get("max"))
		//外键关系
		indexOut := vField.Tag.Get("index_out")
		//默认值
		defaultVal := vField.Tag.Get("default")
		//评价值
		commentVal := vField.Tag.Get("comment")
		//写入t.FieldNameAll
		appendClientField := clientField{
			DBName:           dbVal,
			DBType:           "",
			IsList:           true,
			IsKey:            false,
			IsIndex:          index,
			IsUnique:         unique,
			ValueType:        vType,
			MinLen:           CoreFilter.GetIntByStringNoErr(vField.Tag.Get("min")),
			MaxLen:           maxVal,
			DefaultVal:       defaultVal,
			JSONName:         vField.Tag.Get("json"),
			CheckCode:        vField.Tag.Get("check"),
			IsCreateRequired: true,
		}
		//当前字段是否已经存在
		haveField := false
		for _, nowField := range columnNames {
			if nowField == dbVal {
				haveField = true
				break
			}
		}
		//写入数据
		// 识别预设
		switch dbVal {
		case "id":
			if len(columnNames) < 1 {
				appendFields = append(appendFields, "id bigserial constraint "+t.TableName+"_pk primary key")
				t.installAppendUIndex("id")
			}
			appendClientField.IsIndex = true
			appendClientField.DBType = "bigint"
			appendClientField.IsCreateRequired = false
		case "create_at":
			if len(columnNames) > 0 {
				if !haveField {
					appendFields = append(appendFields, "ALTER TABLE"+" "+t.TableName+" ADD COLUMN IF NOT EXISTS create_at timestamp with time zone default CURRENT_TIMESTAMP not null;")
				}
			} else {
				appendFields = append(appendFields, "create_at timestamp with time zone default CURRENT_TIMESTAMP not null")
			}
			appendClientField.DBType = "timestamp with time zone"
			appendClientField.IsCreateRequired = false
		case "update_at":
			if len(columnNames) > 0 {
				if !haveField {
					appendFields = append(appendFields, "ALTER TABLE"+" "+t.TableName+" ADD COLUMN IF NOT EXISTS update_at timestamp with time zone default CURRENT_TIMESTAMP not null;")
				}
			} else {
				appendFields = append(appendFields, "update_at timestamp with time zone default CURRENT_TIMESTAMP not null")
			}
			appendClientField.DBType = "timestamp with time zone"
			appendClientField.IsCreateRequired = false
		case "delete_at":
			if len(columnNames) > 0 {
				if !haveField {
					appendFields = append(appendFields, "ALTER TABLE"+" "+t.TableName+" ADD COLUMN IF NOT EXISTS delete_at timestamp with time zone default to_timestamp((0)::double precision) not null;")
				}
			} else {
				appendFields = append(appendFields, "delete_at timestamp with time zone default to_timestamp((0)::double precision) not null")
			}
			appendClientField.DBType = "timestamp with time zone"
			appendClientField.IsCreateRequired = false
		case "code":
			if maxVal < 1 {
				maxVal = 50
			}
			if len(columnNames) > 0 {
				if !haveField {
					appendFields = append(appendFields, fmt.Sprint("ALTER TABLE"+" "+t.TableName+" ADD COLUMN IF NOT EXISTS code varchar(", maxVal, ") default '"+defaultVal+"' not null;"))
				}
			} else {
				appendFields = append(appendFields, fmt.Sprint("code varchar(", maxVal, ") not null"))
			}
			if unique {
				t.installAppendUIndex("code")
			}
			if index {
				t.installAppendIndex("code")
			}
			appendClientField.DBType = "varchar"
		case "mark":
			if maxVal < 1 {
				maxVal = 50
			}
			if len(columnNames) > 0 {
				if !haveField {
					appendFields = append(appendFields, fmt.Sprint("ALTER TABLE"+" "+t.TableName+" ADD COLUMN IF NOT EXISTS mark varchar(", maxVal, ") default '"+defaultVal+"' not null;"))
				}
			} else {
				appendFields = append(appendFields, fmt.Sprint("mark varchar(", maxVal, ") not null"))
			}
			if unique {
				t.installAppendUIndex("mark")
			}
			if index {
				t.installAppendIndex("mark")
			}
			appendClientField.DBType = "varchar"
		default:
			appendTypeSQL := ""
			appendDefaultSQL := ""
			switch vType {
			case "int":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default 0"
				}
				appendTypeSQL = "integer"
				appendClientField.DBType = "integer"
			case "[]int":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default '{}'::integer[]"
				}
				appendTypeSQL = "integer[]"
				appendClientField.DBType = "integer[]"
			case "pq.Int32Array":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default '{}'::integer[]"
				}
				appendTypeSQL = "integer[]"
			case "int64":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default 0"
				}
				appendTypeSQL = "bigint"
				appendClientField.DBType = "bigint"
			case "[]int64":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default '{}'::bigint[]"
				}
				appendTypeSQL = "bigint[]"
				appendClientField.DBType = "bigint[]"
			case "pq.Int64Array":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default '{}'::bigint[]"
				}
				appendTypeSQL = "bigint[]"
				appendClientField.DBType = "bigint[]"
			case "float64":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default 0"
				}
				appendTypeSQL = "float"
				appendClientField.DBType = "float"
			case "[]float64":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default '{}'::float[]"
				}
				appendTypeSQL = "float[]"
				appendClientField.DBType = "float[]"
			case "bool":
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default false"
				}
				appendTypeSQL = "boolean"
				appendClientField.DBType = "boolean"
			case "time.Time":
				if defaultVal != "" {
					if defaultVal == "now()" {
						appendDefaultSQL = " default CURRENT_TIMESTAMP"
					} else {
						if defaultVal == "0" {
							appendDefaultSQL = " default to_timestamp((0)::double precision)"
						}
					}
				} else {
					appendDefaultSQL = " default to_timestamp((0)::double precision)"
				}
				appendTypeSQL = "timestamp with time zone"
				appendClientField.DBType = "timestamp with time zone"
			case "string":
				if defaultVal != "" {
					appendDefaultSQL = " default '" + defaultVal + "'"
				} else {
					appendDefaultSQL = " default ''"
				}
				if maxVal == -1 {
					appendTypeSQL = "text"
					appendClientField.DBType = "text"
					appendClientField.IsList = false
				} else {
					if maxVal < 1 {
						maxVal = 255
					}
					appendTypeSQL = fmt.Sprint("varchar(", maxVal, ")")
					appendClientField.DBType = "varchar"
				}
			case "[]string":
				if maxVal < 1 {
					maxVal = 255
				}
				if defaultVal != "" {
					appendDefaultSQL = " default '" + defaultVal + "'"
				} else {
					appendDefaultSQL = " default '{}'::varchar[]"
				}
				appendTypeSQL = fmt.Sprint("varchar(", maxVal, ")[]")
				appendClientField.DBType = fmt.Sprint("varchar(", maxVal, ")[]")
			case "pq.StringArray":
				if maxVal < 1 {
					maxVal = 255
				}
				if defaultVal != "" {
					appendDefaultSQL = " default '" + defaultVal + "'"
				} else {
					appendDefaultSQL = " default '{}'::varchar[]"
				}
				appendTypeSQL = fmt.Sprint("varchar(", maxVal, ")[]")
				appendClientField.DBType = fmt.Sprint("varchar(", maxVal, ")[]")
			default:
				//err = errors.New("install sql error: table " + t.TableName + " field " + dbVal + " type " + vType + " not support")
				//return
				//按照jsonb处理，不建议使用
				if defaultVal != "" {
					appendDefaultSQL = " default " + defaultVal
				} else {
					appendDefaultSQL = " default '{}'::jsonb"
				}
				appendTypeSQL = "jsonb"
				appendClientField.DBType = "jsonb"
			}
			//检查是否需继续写入数据
			if len(columnNames) > 0 {
				if !haveField {
					//需要删除字段
					appendFields = append(appendFields, fmt.Sprint("ALTER TABLE"+" "+t.TableName+" ADD COLUMN IF NOT EXISTS ", dbVal, " ", appendTypeSQL, appendDefaultSQL, " not null;"))
				} else {
					//跳过情况，不处理，但此处不能跳出，因为还需后续进行处理优化
				}
			} else {
				appendFields = append(appendFields, fmt.Sprint(dbVal, " ", appendTypeSQL, appendDefaultSQL, " not null"))
			}
			//需要跳过处理，因为字段不需要做任何额外的处理
			if unique {
				appendClientField.IsUnique = true
				t.installAppendUIndex(dbVal)
			}
			if index {
				appendClientField.IsIndex = true
				t.installAppendIndex(dbVal)
			}
		}
		//写入评论值
		if commentVal != "" {
			appendFields = append(appendFields, fmt.Sprintf("COMMENT ON COLUMN %s.%s IS %s;", t.TableName, dbVal, commentVal))
		}
		//TODO: 尚未支持外键
		if indexOut != "" {
			err = errors.New("install sql error: table " + t.TableName + " has index_out, not support")
			return
		}
		t.fieldNameList = append(t.fieldNameList, appendClientField)
	}
	//检查主键数量
	if t.installNunIndexKeyNum > 1 {
		err = errors.New("install sql error: table " + t.TableName + " has more than one primary key")
		return
	}
	//反向核查字段是否被删除
	// alter table if exists table_name drop column if exists field_name;
	for _, nowField := range columnNames {
		isFind := false
		for _, addField := range needAddFields {
			if addField == nowField {
				isFind = true
				break
			}
		}
		if !isFind {
			t.installAppendSQLData = append(t.installAppendSQLData, fmt.Sprintf("alter table if exists %s drop column if exists %s;", t.TableName, nowField))
		}
	}
	//追加sql
	if len(columnNames) > 0 {
		sqlData += strings.Join(appendFields, "")
	} else {
		sqlData += strings.Join(appendFields, ",") + ");"
	}
	sqlData += strings.Join(t.installAppendSQLData, "")
	//执行sql
	_, err = t.DB.GetPostgresql().Exec(sqlData)
	if err != nil {
		err = errors.New(fmt.Sprint("install sql exec error: "+err.Error(), ", ", sqlData))
		return
	}
	//标记已运行
	t.installHaveRun = true
	//反馈
	return
}

// installAppendUIndex 给表格插入uindex序列
func (t *Client) installAppendUIndex(fieldName string) {
	appendSQL := "create unique index if not exists " + t.TableName + "_" + fieldName + "_uindex on " + t.TableName + " (" + fieldName + ");"
	for _, v := range t.installAppendSQLData {
		if v == appendSQL {
			return
		}
	}
	t.installAppendSQLData = append(t.installAppendSQLData, appendSQL)
	t.installNunIndexKeyNum += 1
	return
}

// installAppendIndex 给表格插入index序列
func (t *Client) installAppendIndex(fieldName string) {
	appendSQL := "create index if not exists " + t.TableName + "_" + fieldName + "_index on " + t.TableName + " (" + fieldName + ");"
	for _, v := range t.installAppendSQLData {
		if v == appendSQL {
			return
		}
	}
	t.installAppendSQLData = append(t.installAppendSQLData, appendSQL)
	return
}
