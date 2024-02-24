package BaseLookup

func DeleteDomain(id int64) (err error) {
	// 标记删除
	err = domainDB.Delete().AddWhereID(id).NeedSoft(true).ExecNamed(map[string]any{})
	if err != nil {
		return
	}
	//删除缓冲
	deleteDomainCache(id)
	return
}
