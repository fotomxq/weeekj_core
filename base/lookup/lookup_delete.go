package BaseLookup

func DeleteLookup(id int64) (err error) {
	// 标记删除
	err = lookupDB.Delete().AddWhereID(id).NeedSoft(true).ExecNamed(map[string]any{})
	if err != nil {
		return
	}
	//删除缓冲
	deleteLookupCache(id)
	return
}
