package MapRoom

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsCreateRoom 创建新的房间参数
type ArgsCreateRoom struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//入驻人员列
	Infos pq.Int64Array `db:"infos" json:"infos" check:"ids" empty:"true"`
	//房间编号
	Code string `db:"code" json:"code" check:"mark" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateRoom 创建新的房间
func CreateRoom(args *ArgsCreateRoom) (data FieldsRoom, err error) {
	if args.Code != "" {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND code = $2", args.OrgID, args.Code)
		if err == nil && data.ID > 0 {
			err = errors.New(fmt.Sprint("mark is exist, ", err))
			return
		}
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "map_room", "INSERT INTO map_room (org_id, sort_id, tags, status, service_status, service_bind_id, service_mission_id, infos, code, name, des, cover_file_id, des_files, params) VALUES (:org_id,:sort_id,:tags,3, 0, 0, 0,:infos,:code,:name,:des,:cover_file_id,:des_files,:params)", args, &data)
	if err != nil {
		return
	}
	//推送更新
	pushNatsUpdateStatus(data.ID, "create", "")
	//统计数据
	pushNatsUpdateAnalysis(data.OrgID)
	//反馈
	return
}
