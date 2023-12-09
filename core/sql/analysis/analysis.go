package CoreSQLAnalysis

//GetAnalysisQueryField 时间范围抽取，生成请求的字段部分，后续可直接用于分组实现
func GetAnalysisQueryField(fieldName string, timeType string, newFieldName string) (query string) {
	switch timeType{
	case "year":
		query = "TO_CHAR(" + fieldName + ", 'YYYY') AS " + newFieldName
	case "month":
		query = "TO_CHAR(" + fieldName + ", 'YYYY-MM') AS " + newFieldName
	case "day":
		query = "TO_CHAR(" + fieldName + ", 'YYYY-MM-DD') AS " + newFieldName
	case "hour":
		query = "TO_CHAR(" + fieldName + ", 'YYYY-MM-DD HH24') AS " + newFieldName
	case "minute":
		query = "TO_CHAR(" + fieldName + ", 'YYYY-MM-DD HH24:MI') AS " + newFieldName
	case "sec":
		query = "TO_CHAR(" + fieldName + ", 'YYYY-MM-DD HH24:MI:SS') AS " + newFieldName
	default:
		query = "TO_CHAR(" + fieldName + ", 'YYYY-MM-DD HH24:MI:SS') AS " + newFieldName
	}
	return
}