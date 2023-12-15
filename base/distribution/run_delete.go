package BaseDistribution

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
)

// 子服务2，检查关联服务的存在状态，如不存在将自动移除服务、子服务、子服务run
func runDelete() error {
	//分批获取负载数据
	//已经检查过的数据列
	var checkServiceList []string
	//遍历数据
	for _, vChild := range cacheChildData {
		//跳过检查过的数据，加速检查
		isFind := false
		for _, v := range checkServiceList {
			if v == vChild.Mark {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		//获取主服务是否存在？
		vService, err := getService(vChild.Mark)
		if err != nil {
			//删除数据
			if err := DeleteChild(&ArgsDeleteChild{
				Mark: vChild.Mark, IP: vChild.ServerIP, Port: vChild.ServerPort,
			}); err != nil {
				CoreLog.Error("delete child, ", err)
				continue
			}
		} else {
			//存在则写入列表
			checkServiceList = append(checkServiceList, vService.Mark)
		}
	}
	return nil
}
