package BaseUnit

func DeleteUnit(id int64) (err error) {
	// 标记删除
	err = unitDB.Delete().AddWhereID(id).NeedSoft(true).ExecNamed(map[string]any{})
	if err != nil {
		return
	}
	//删除缓冲
	deleteUnitCache(id)
	return
}
