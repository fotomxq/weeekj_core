package CoreSQL2

// InstallSQL SQL自动安装工具
/**
dataDefault 初始化需采用的空数据集
1. 如果给予值，则代表有默认值，将按照默认值构建数据表
2. eg: dataDefault = data = ClassSort.FieldsSort{}

识别规则(tag) :
db: 数据库字段名
index=true: 主键索引
值类型: 数据库字段类型
max: 最大长度
index_out: 外键索引
default: 默认值，sql直接写入

值类型转化对应关系:
int64: bigint
[]int64: bigint[]
pq.Int64Array: bigint[]
int: integer
[]int: integer[]
pq.Int32Array: integer[]
bool: boolean
time.Time: timestamp
string: varchar(max)
string: text

预设值规则(tag) :
createAt: 创建时间
updateAt: 更新时间
deleteAt: 删除时间
code: 编码
*/
func InstallSQL(dataDefault any) {

}
