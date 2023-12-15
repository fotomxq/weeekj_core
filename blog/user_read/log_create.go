package BlogUserRead

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BlogCore "github.com/fotomxq/weeekj_core/v5/blog/core"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsCreateLog 添加日志参数
type ArgsCreateLog struct {
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//阅读渠道
	// 访问渠道的特征码
	FromMark string `db:"from_mark" json:"fromMark" check:"mark"`
	FromName string `db:"from_name" json:"fromName"`
	//姓名
	Name string `db:"name" json:"name" check:"name" empty:"true"`
	//IP
	IP string `db:"ip" json:"ip" check:"ip"`
	//文章ID
	ContentID int64 `db:"content_id" json:"contentID" check:"id"`
	//进入时间
	CreateAt string `db:"create_at" json:"createAt" check:"isoTime"`
	//离开时间
	LeaveAt string `db:"leave_at" json:"leaveAt" check:"isoTime" empty:"true"`
}

// CreateLog 添加日志
func CreateLog(args *ArgsCreateLog) (err error) {
	//获取进入时间
	var createAt time.Time
	if args.CreateAt != "" {
		createAt, err = CoreFilter.GetTimeByISO(args.CreateAt)
		if err != nil {
			err = errors.New(fmt.Sprint("get create at, ", err))
			return
		}
	} else {
		createAt = CoreFilter.GetNowTime()
	}
	//获取离开时间
	var leaveAt time.Time
	if args.LeaveAt != "" {
		leaveAt, err = CoreFilter.GetTimeByISO(args.LeaveAt)
		if err != nil {
			err = errors.New(fmt.Sprint("get leave at, ", err))
			return
		}
	}
	//获取文章数据
	var contentData BlogCore.FieldsContent
	contentData, err = BlogCore.GetContentByID(&BlogCore.ArgsGetContentByID{
		ID:        args.ContentID,
		OrgID:     -1,
		IsPublish: true,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get content data, ", err))
		return
	}
	//获取组织信息
	var parentOrgID int64
	var childOrgID int64
	if contentData.OrgID > 0 {
		var orgData OrgCore.FieldsOrg
		orgData, err = OrgCore.GetOrg(&OrgCore.ArgsGetOrg{
			ID: contentData.OrgID,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("get org data, org id: ", contentData.OrgID, ", err: ", err))
			return
		}
		//如果存在上级组织，则修改规则
		if args.ChildOrgID < 1 {
			if orgData.ParentID > 0 {
				parentOrgID = orgData.ParentID
				childOrgID = orgData.ID
			} else {
				parentOrgID = orgData.ID
				childOrgID = 0
			}
		} else {
			parentOrgID = contentData.OrgID
			childOrgID = args.ChildOrgID
		}
	}
	//计算阅读时间
	var readTime int64 = 0
	if args.LeaveAt != "" {
		readTime = leaveAt.Unix() - createAt.Unix()
		if readTime < 1 {
			readTime = 0
		}
	}
	//获取重复记录次数配置
	var blogUserNoReplaceTime int64
	blogUserNoReplaceTime, err = BaseConfig.GetDataInt64("BlogUserNoReplaceTime")
	if err != nil {
		blogUserNoReplaceTime = 3600
	}
	var blogUserReadOnce bool
	blogUserReadOnce, err = BaseConfig.GetDataBool("BlogUserReadOnce")
	if err != nil {
		blogUserReadOnce = true
	}
	//添加日志数据
	leaveAt, readTime, err = createLogLog(args.UserID, args.Name, args.IP, args.FromMark, args.FromName, parentOrgID, childOrgID, args.ContentID, contentData.SortID, blogUserReadOnce, blogUserNoReplaceTime, createAt, leaveAt, readTime)
	if err != nil {
		return
	}
	/** 如果没插入数据，必定拦截，所以不需要执行此段
	else {
		if !blogUserReadOnce {
			err = createLogAnalysis(args.UserID, args.Name, args.IP, args.FromMark, args.FromName, parentOrgID, childOrgID, contentData.SortID, createAt, readTime)
			if err != nil {
				return
			}
		}
	}
	*/
	//反馈
	return
}

// 记录日志记录
func createLogLog(userID int64, name string, ip string, fromMark, fromName string, orgID, childOrgID int64, contentID int64, sortID int64, blogUserReadOnce bool, blogUserNoReplaceTime int64, createAt, leaveAt time.Time, readTime int64) (newLeaveAt time.Time, newReadTime int64, err error) {
	//初始化数据
	newLeaveAt = leaveAt
	newReadTime = readTime
	//检查是否存在用户?
	var data FieldsLog
	if userID > 0 {
		//以用户为主进行记录
		// 尝试获取数据
		data = getLogCache(contentID, userID)
		if data.ID < 1 {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, child_org_id, user_id, from_mark, from_name, name, ip, sort_id, content_id, leave_at, read_time FROM blog_user_read_log WHERE user_id = $1 AND org_id = $2 AND child_org_id = $3 AND content_id = $4 AND sort_id = $5 ORDER BY id DESC LIMIT 1", userID, orgID, childOrgID, contentID, sortID)
			if err == nil && data.ID > 0 {
				setLogCache(data)
			}
		} else {
			if !CoreFilter.EqID2(orgID, data.OrgID) || !CoreFilter.EqID2(childOrgID, data.ChildOrgID) || !CoreFilter.EqID2(sortID, data.SortID) {
				err = errors.New("no data")
			}
		}
	} else {
		//以IP为主进行记录
		var data FieldsLog
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, leave_at FROM blog_user_read_log WHERE (ip = $1 AND from_mark = $2 AND from_name = $3) AND org_id = $4 AND child_org_id = $5 AND content_id = $6 AND sort_id = $7 ORDER BY id DESC LIMIT 1", ip, fromMark, fromName, orgID, childOrgID, contentID, sortID)
	}
	//判断是否存在数据，影响记录方式
	if err == nil && data.ID > 0 {
		//存在数据，则只是更新数据
		if blogUserReadOnce {
			return
		}
		isCreate := false
		//检查是否为离开时刻
		if data.LeaveAt.Unix() > 100000 {
			//如果离开时间少于blogUserNoReplaceTime，则退出
			if CoreFilter.GetNowTime().Unix()-data.LeaveAt.Unix() < blogUserNoReplaceTime {
				return
			}
			//否则按照新的数据来处理
			isCreate = true
		} else {
			//昨天之前创建的数据，也列为离开状态
			if data.CreateAt.Unix() <= CoreFilter.GetNowTimeCarbon().SubDay().Time.Unix() {
				leaveAt = CoreFilter.GetNowTime()
				readTime = leaveAt.Unix() - data.CreateAt.Unix()
				newLeaveAt = leaveAt
				newReadTime = readTime
			}
		}
		//更新数据
		if !isCreate {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE blog_user_read_log SET leave_at = :leave_at, read_time = :read_time WHERE id = :id", map[string]interface{}{
				"id":        data.ID,
				"leave_at":  leaveAt,
				"read_time": readTime,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("update blog user read log, err: ", err))
				return
			}
			return
		}
	}
	//不存在数据，则添加新的记录
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_user_read_log (create_at, org_id, child_org_id, user_id, from_mark, from_name, name, ip, content_id, sort_id, leave_at, read_time) VALUES (:create_at,:org_id,:child_org_id,:user_id,:from_mark,:from_name,:name,:ip,:content_id,:sort_id,:leave_at,:read_time)", map[string]interface{}{
		"create_at":    createAt,
		"org_id":       orgID,
		"child_org_id": childOrgID,
		"user_id":      userID,
		"from_mark":    fromMark,
		"from_name":    fromName,
		"name":         name,
		"ip":           ip,
		"content_id":   contentID,
		"sort_id":      sortID,
		"leave_at":     leaveAt,
		"read_time":    readTime,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("insert blog user read log, err: ", err))
		return
	}
	//添加统计数据
	err = createLogAnalysis(data.UserID, data.Name, data.IP, data.FromMark, data.FromName, data.OrgID, data.ChildOrgID, data.SortID, createAt, readTime)
	if err != nil {
		return
	}
	return
}

// 记录统计记录
func createLogAnalysis(userID int64, name string, ip string, fromMark, fromName string, orgID, childOrgID int64, sortID int64, createAt time.Time, readTime int64) (err error) {
	//检查是否存在用户?
	var data FieldsAnalysis
	if userID > 0 {
		//以用户为主进行记录
		// 尝试获取数据
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at FROM blog_user_read_analysis WHERE user_id = $1 AND org_id = $2 AND child_org_id = $3 AND sort_id = $4", userID, orgID, childOrgID, sortID)
	} else {
		//以IP为主进行记录
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at FROM blog_user_read_analysis WHERE ip = $1 AND from_mark = $2 AND from_name = $3 AND org_id = $4 AND child_org_id = $5 AND sort_id = $6", ip, fromMark, fromName, orgID, childOrgID, sortID)
	}
	//检查是否数据
	if err == nil && data.ID > 0 {
		//更新数据
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE blog_user_read_analysis SET read_time = read_time + :read_time, read_count = read_count + 1 WHERE id = :id", map[string]interface{}{
			"id":        data.ID,
			"read_time": readTime,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update blog user read analysis, analysis id: ", data.ID, ", err: ", err))
			return
		}
		return
	}
	//添加新的数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_user_read_analysis (create_at, org_id, child_org_id, user_id, from_mark, from_name, name, ip, sort_id, read_time, read_count) VALUES (:create_at,:org_id,:child_org_id,:user_id,:from_mark,:from_name,:name,:ip,:sort_id,:read_time,:read_count)", map[string]interface{}{
		"create_at":    createAt,
		"org_id":       orgID,
		"child_org_id": childOrgID,
		"user_id":      userID,
		"from_mark":    fromMark,
		"from_name":    fromName,
		"name":         name,
		"ip":           ip,
		"sort_id":      sortID,
		"read_time":    readTime,
		"read_count":   1,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("insert blog user read analysis, err: ", err))
		return
	}
	//反馈成功
	return
}
