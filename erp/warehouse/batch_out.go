package ERPWarehouse

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

// ArgsBatchOut 批次出库参数
type ArgsBatchOut struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//数量
	// 注意必须是正数，代表要出库的数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//调用模块
	ActionSystem string `db:"action_system" json:"actionSystem" check:"mark" empty:"true"`
	//调用模块ID
	ActionID int64 `db:"action_id" json:"actionID" check:"id" empty:"true"`
}

// BatchOut 批次出库
func BatchOut(args *ArgsBatchOut) (errCode string, err error) {
	//获取原始数据
	batchData := getBatchByID(args.ID)
	if batchData.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, batchData.OrgID) {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//出库数量如果超出批次，则拒绝
	if args.Count > batchData.Count {
		errCode = "err_erp_warehouse_batch_count_out_more"
		err = errors.New("too more batch count")
		return
	}
	//如果少于批次数量，则修改批次数量
	if args.Count <= batchData.Count {
		//执行修改
		err = batchSQL.Update().SetFieldStr("count = count - :count").AddWhereID(batchData.ID).NamedExec(map[string]any{
			"count": args.Count,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
	}
	//如果和批次数量相同，则删除批次
	if args.Count == batchData.Count {
		//执行删除
		err = batchSQL.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(map[string]any{})
		if err != nil {
			errCode = "err_delete"
			return
		}
	}
	//新增批次出库数据
	if args.ActionSystem != "" {
		//err = batchSQL.Update().SetFields([]string{"action_system", "action_id"}).AddWhereID(batchData.ID).NamedExec(map[string]any{
		//	"action_system": args.ActionSystem,
		//	"action_id":     args.ActionID,
		//})
		//if err != nil {
		//	errCode = "err_update"
		//	return
		//}
	}
	//修改产品库存台账信息
	errCode, err = setStore(batchData.OrgID, batchData.WarehouseID, batchData.AreaID, batchData.ProductID, true, 0-args.Count)
	if err != nil {
		err = errors.New(fmt.Sprint("set store code: ", errCode, ", err: ", err))
		return
	}
	//删除缓冲
	deleteBatchCache(args.ID)
	//反馈
	return
}
