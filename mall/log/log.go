package MallLog

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	MallCore "github.com/fotomxq/weeekj_core/v5/mall/core"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

// AppendLogByC 使用顶层路由记录访问记录
func AppendLogByC(c *gin.Context, productID int64, action int) {
	//尝试获取路由用户ID和组织ID信息
	userID, _ := Router2Mid.TryGetUserID(c)
	//获取商品信息
	productData, err := MallCore.GetProduct(&MallCore.ArgsGetProduct{
		ID:    productID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	//记录数据
	AppendLog(userID, c.ClientIP(), productData.OrgID, productID, action)
}

// AppendLog 添加一个记录
func AppendLog(userID int64, ip string, orgID int64, productID int64, action int) {
	//如果存在用户ID，则检查是否存在相同记录了？
	if userID > 0 {
		var data FieldsLog
		err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, product_id, action FROM mall_log WHERE user_id = $1 ORDER BY id DESC LIMIT 1", userID)
		if err == nil && data.ID > 0 {
			if data.ProductID == productID && data.Action == action {
				return
			}
		}
	}
	//添加记录
	_, err := CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO mall_log (org_id, user_id, ip, product_id, action) VALUES (:org_id,:user_id,:ip,:product_id,:action)", map[string]interface{}{
		"org_id":     orgID,
		"user_id":    userID,
		"ip":         ip,
		"product_id": productID,
		"action":     action,
	})
	if err != nil {
		CoreLog.Error("mall log append log, insert, ", err)
		return
	}
}

// DataNearNextProductID 反馈数据排名结构体
type DataNearNextProductID struct {
	//商品ID
	ProductID int64 `json:"productID"`
	//得分
	// 根据action打分:
	// action: 0 浏览行为 +1分; 1 评论行为 +2分; 2 购物车行为 +3分; 3 购买行为 +4分
	Score int `json:"score"`
}

// NearNextProductID 追溯某个商品的下一个商品分布前N名
// 可直接用于简单的推荐算法模型，获取某个商品出现的分布情况做排名，排名靠前的会列入推荐算法列队
// 设计思路：https://blog.csdn.net/qq_39564555/article/details/105881352
// 按照相似度顺序，反馈一组符合条件的商品ID
func NearNextProductID(productID int64, limit int) (scoreList []DataNearNextProductID) {
	//无限循环，直到达到limit或没有数据
	step := 0
	var haveIDs []int64
	for {
		//检查限制
		if len(scoreList) >= limit {
			break
		}
		//获取该商品最近的记录
		var dataList []FieldsLog
		err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM mall_log WHERE product_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", productID, 100, step)
		if err != nil || len(dataList) < 1 {
			break
		}
		//遍历记录，获取所有id的+1数据
		for _, v := range dataList {
			//检查限制
			if len(scoreList) >= limit {
				break
			}
			//获取ID的下一个非本商品的记录
			var vData FieldsLog
			if err := Router2SystemConfig.MainDB.Get(&vData, "SELECT id, product_id, action FROM mall_log WHERE id > $1 AND product_id != $2 ORDER BY id LIMIT 1", v.ID, productID); err != nil {
				continue
			}
			//检查已经筛查过的ID序列，避免重复叠加
			isFindHaveID := false
			for _, v2 := range haveIDs {
				if v2 == vData.ID {
					isFindHaveID = true
					break
				}
			}
			if isFindHaveID {
				continue
			}
			haveIDs = append(haveIDs, vData.ID)
			//计算分数
			score := 0
			switch vData.Action {
			case 0:
				score = 1
			case 1:
				score = 2
			case 2:
				score = 3
			case 3:
				score = 4
			}
			//计入分数或增加数据
			isFind := false
			for k2, v2 := range scoreList {
				if v2.ProductID == vData.ProductID {
					scoreList[k2].Score += score
					isFind = true
					break
				}
			}
			if !isFind {
				scoreList = append(scoreList, DataNearNextProductID{
					ProductID: vData.ProductID,
					Score:     score,
				})
			}
		}
		step += 100
	}
	//反馈
	return
}
