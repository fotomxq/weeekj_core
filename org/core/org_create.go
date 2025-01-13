package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserMessageMod "github.com/fotomxq/weeekj_core/v5/user/message/mod"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateOrg 创建组织参数
type ArgsCreateOrg struct {
	//所属用户
	// 掌管该数据的用户，创建人和根管理员，不可删除只能更换
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//企业唯一标识码
	// 用于特殊识别和登陆识别等操作
	Key string `db:"key" json:"key" check:"mark" empty:"true"`
	//构架名称，或组织名称
	Name string `db:"name" json:"name" check:"name"`
	//组织描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000"`
	//上级ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//上级控制权限限制
	ParentFunc pq.StringArray `db:"parent_func" json:"parentFunc" check:"marks" empty:"true"`
	//开通业务
	// 该内容只有总管理员或订阅能进行控制
	OpenFunc pq.StringArray `db:"open_func" json:"openFunc" check:"marks" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
}

// CreateOrg 创建组织
func CreateOrg(args *ArgsCreateOrg) (orgData FieldsOrg, errCode string, err error) {
	//生成key
	if args.Key == "" {
		args.Key = makeKey(args.Name)
	}
	//获取用户信息
	if args.UserID < 1 {
		errCode = "err_user"
		err = errors.New("user id is empty")
		return
	}
	var userData UserCore.FieldsUserType
	userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    args.UserID,
		OrgID: -1,
	})
	if err != nil {
		errCode = "err_user"
		err = errors.New("get user by id, " + err.Error())
		return
	}
	//检查重名
	oldOrgID := getOrgByName(args.Name)
	if oldOrgID > 0 {
		errCode = "err_name"
		err = errors.New("org name exist")
		return
	}
	//如果开通功能为空，则默认给开通功能
	if len(args.OpenFunc) < 1 {
		args.OpenFunc = Router2SystemConfig.GlobConfig.Org.DefaultOpenFunc
	}
	//检查上一级
	if args.ParentID > 0 {
		//获取上一级
		var parentData FieldsOrg
		parentData, err = GetOrg(&ArgsGetOrg{
			ID: args.ParentID,
		})
		if err != nil {
			errCode = "err_parent"
			err = errors.New("get parent org, " + err.Error())
			return
		}
		//检查上一级功能和本次开通功能
		err = checkOrgInParentFunc(parentData.OpenFunc, args.OpenFunc)
		if err != nil {
			errCode = "err_func_not_in_area"
			err = errors.New("open func not in parent func, " + err.Error())
			return
		}
	}
	//检查key
	if args.Key == "" {
		args.Key, err = CoreFilter.GetRandStr3(10)
		if err != nil {
			errCode = "err_key"
			err = errors.New("get rand key, " + err.Error())
			return
		}
	}
	orgData, err = GetOrgByKey(&ArgsGetOrgByKey{
		Key: args.Key,
	})
	if err == nil && orgData.ID > 0 {
		errCode = "err_key"
		err = errors.New("key is exist")
		return
	}
	//生成数据
	if err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_core", "INSERT INTO org_core (user_id, name, key, des, parent_id, parent_func, open_func, sort_id) VALUES (:user_id, :name, :key, :des, :parent_id, :parent_func, :open_func, :sort_id)", args, &orgData); err != nil {
		errCode = "err_insert"
		err = errors.New("create org, " + err.Error())
		return
	}
	//建立绑定关系
	if _, err = SetBind(&ArgsSetBind{
		UserID:     args.UserID,
		Avatar:     0,
		Name:       userData.Name,
		OrgID:      orgData.ID,
		GroupIDs:   []int64{},
		Manager:    []string{"member", "all"},
		NationCode: userData.NationCode,
		Phone:      userData.Phone,
		Email:      userData.Email,
		SyncSystem: "",
		SyncID:     0,
		SyncHash:   "",
		Params:     CoreSQLConfig.FieldsConfigsType{},
	}); err != nil {
		errCode = "err_org_bind"
		err = errors.New("set org bind, " + err.Error())
		return
	}
	//发送用户消息
	if orgData.UserID > 0 {
		UserMessageMod.CreateSystemToUser(time.Time{}, orgData.UserID, "成功开通商户", fmt.Sprint("您已经成功开通了商户(", orgData.Name, ")。"), nil, nil)
	}
	//反馈
	return
}
