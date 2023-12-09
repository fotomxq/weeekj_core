package BaseQiniu

//ArgsDeleteWithDay 定时删除某个文件参数
type ArgsDeleteWithDay struct {
	//桶名称
	Bucket string
	//key
	Key string
	//多少天以后删除
	Days int
}

//DeleteWithDay 定时删除某个文件
func DeleteWithDay(args *ArgsDeleteWithDay) error {
	var err error
	bucketManager, err = getManager()
	if err != nil {
		return err
	}
	return bucketManager.DeleteAfterDays(args.Bucket, args.Key, args.Days)
}
