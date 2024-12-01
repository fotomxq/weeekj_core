package AnalysisIndexFilter

// Append 追加数据
func Append(args *FieldsFilter) (err error) {
	_, err = filterDB.GetInsert().InsertRow(args)
	return
}

// AppendList 追加数据列表
func AppendList(args []FieldsFilter) (err error) {
	for _, v := range args {
		err = Append(&v)
		if err != nil {
			return
		}
	}
	return
}
