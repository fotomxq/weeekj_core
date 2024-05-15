package CoreSQL2

import (
	"fmt"
)

//TODO: 需要继续开发，当前只有sort排序方案可以参考，但可能还有其他通用形式
// 除此之外，还需确保针对分组日期的查询方案。不过日期差分方式，和排序查询逻辑上不兼容，需考虑特性和应用场景继续开发

// ClientAnalysisCtx 排名处理器
type ClientAnalysisCtx struct {
	//对象
	clientCtx *ClientCtx
	//需要聚合的字段
	countFields []ClientAnalysisCountFieldCtx
	//聚合字段
	groupFields []string
	//条件
	where string
	//限制
	limit int
	//时间范围
	betweenAt ArgsTimeBetween
	//排序方式
	sort   []string
	isDesc bool
	//是否判断删除操作
	sqlNeedNoDelete bool
}

type ClientAnalysisCountFieldCtx struct {
	//字段名
	Name string
	//聚合方法
	// sum 加; avg 平均
	Type string
}

func (t *ClientAnalysisCtx) getMode(mode string) (query string) {
	switch mode {
	case "day":
		query = "to_char(create_at::DATE, 'YYYY-MM-DD') as day_time"
	case "week":
		query = "to_char(create_at::DATE-(extract(dow from create_at::TIMESTAMP)-1||'day')::interval, 'YYYY-mm-dd') day_time"
	case "year":
		query = "to_char(create_at::DATE, 'YYYY') as day_time"
	default:
		query = "to_char(create_at::DATE, 'YYYY-MM') as day_time"
	}
	return
}

// DataNoDelete SQL结构体不含删除标记
func (t *ClientAnalysisCtx) DataNoDelete() *ClientAnalysisCtx {
	t.sqlNeedNoDelete = true
	return t
}

// Count 计算总数
// eg: where = "config = $1", args = "[config_id]"...
func (t *ClientAnalysisCtx) Count(where string, args ...any) (val int64) {
	if t.sqlNeedNoDelete {
		_ = t.clientCtx.Get(&val, t.clientCtx.DataNoDelete().getSQLWhere(fmt.Sprint("SELECT COUNT("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	} else {
		_ = t.clientCtx.Get(&val, t.clientCtx.getSQLWhere(fmt.Sprint("SELECT COUNT("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	}
	return
}

func (t *ClientAnalysisCtx) Max(where string, args ...any) (val int64) {
	if t.sqlNeedNoDelete {
		_ = t.clientCtx.Get(&val, t.clientCtx.DataNoDelete().getSQLWhere(fmt.Sprint("SELECT MAX("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	} else {
		_ = t.clientCtx.Get(&val, t.clientCtx.getSQLWhere(fmt.Sprint("SELECT MAX("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	}
	return
}

func (t *ClientAnalysisCtx) Min(where string, args ...any) (val int64) {
	if t.sqlNeedNoDelete {
		_ = t.clientCtx.Get(&val, t.clientCtx.DataNoDelete().getSQLWhere(fmt.Sprint("SELECT MIN("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	} else {
		_ = t.clientCtx.Get(&val, t.clientCtx.getSQLWhere(fmt.Sprint("SELECT MIN("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	}
	return
}

func (t *ClientAnalysisCtx) AVG(where string, args ...any) (val float64) {
	if t.sqlNeedNoDelete {
		_ = t.clientCtx.Get(&val, t.clientCtx.DataNoDelete().getSQLWhere(fmt.Sprint("SELECT AVG("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	} else {
		_ = t.clientCtx.Get(&val, t.clientCtx.getSQLWhere(fmt.Sprint("SELECT AVG("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	}
	return
}

func (t *ClientAnalysisCtx) Sum2(where string, args ...any) (val int64) {
	if t.sqlNeedNoDelete {
		_ = t.clientCtx.Get(&val, t.clientCtx.DataNoDelete().getSQLWhere(fmt.Sprint("SELECT SUM("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	} else {
		_ = t.clientCtx.Get(&val, t.clientCtx.getSQLWhere(fmt.Sprint("SELECT SUM("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	}
	return
}

func (t *ClientAnalysisCtx) Sum2ByFloat64(where string, args ...any) (val float64) {
	if t.sqlNeedNoDelete {
		_ = t.clientCtx.Get(&val, t.clientCtx.DataNoDelete().getSQLWhere(fmt.Sprint("SELECT SUM("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	} else {
		_ = t.clientCtx.Get(&val, t.clientCtx.getSQLWhere(fmt.Sprint("SELECT SUM("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where), args...)
	}
	return
}

func (t *ClientAnalysisCtx) Sum(data any, fields []string, where string, args ...any) (err error) {
	sumQuery := ""
	if len(fields) == 1 {
		sumQuery = fmt.Sprint("SUM(", fields[0], ") as ", fields[0])
	} else {
		for k := 0; k < len(fields); k++ {
			v := fields[k]
			if k >= len(fields)-1 {
				sumQuery = fmt.Sprint(sumQuery, ", SUM(", v, ") as ", v)
			} else {
				sumQuery = fmt.Sprint("SUM(", v, ") as ", v)
			}
		}
	}
	if t.sqlNeedNoDelete {
		err = t.clientCtx.Get(data, t.clientCtx.DataNoDelete().getSQLWhere(fmt.Sprint("SELECT ", sumQuery, " FROM ", t.clientCtx.client.TableName), where), args...)
	} else {
		err = t.clientCtx.Get(data, t.clientCtx.getSQLWhere(fmt.Sprint("SELECT ", sumQuery, " FROM ", t.clientCtx.client.TableName), where), args...)
	}
	return
}
