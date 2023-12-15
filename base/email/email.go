package BaseEmail

import (
	"crypto/tls"
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"net"
	"net/smtp"
	"time"
)

//email服务模块
// 提供标准协议的发送、接收功能
// 警告，本模块不支持并发，请在并发服务前将该模块摘出

// ArgsGetEmailList 获取列表参数
type ArgsGetEmailList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否成功
	IsSuccess bool
	//来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//送达邮件搜索
	ToEmailSearch string
	//搜索
	Search string
}

// GetEmailList 获取列表
func GetEmailList(args *ArgsGetEmailList) (dataList []FieldsEmailType, dataCount int64, err error) {
	where := "(title ILIKE '%' || :search || '%' OR content ILIKE '%' || :search || '%') AND is_success=:is_success AND to_email ILIKE '%' || :to_email || '%' AND delete_at < to_timestamp(1000000) AND (:org_id < 1 OR org_id = :org_id)"
	maps := map[string]interface{}{
		"search":     args.Search,
		"is_success": args.IsSuccess,
		"to_email":   args.ToEmailSearch,
		"org_id":     args.OrgID,
	}
	where, maps, err = args.CreateInfo.GetListAnd("create_info", "create_info", where, maps)
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_email",
		"id",
		fmt.Sprint(
			"SELECT id, create_at, update_at, delete_at, org_id, server_id, create_info, send_at, is_success, is_failed, fail_message, title, content, content_type, to_email FROM core_email WHERE ",
			where,
		),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "send_at"},
	)
	return
}

// ArgsGetEmailByID 获取某个数据参数
type ArgsGetEmailByID struct {
	//ID
	ID int64
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetEmailByID 获取某个数据
func GetEmailByID(args *ArgsGetEmailByID) (data FieldsEmailType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, server_id, create_info, send_at, is_success, is_failed, fail_message, title, content, content_type, to_email FROM core_email WHERE id=$1 AND delete_at < TO_TIMESTAMP(1) AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsSend 申请发送一份邮件参数
type ArgsSend struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//服务配置ID
	ServerID int64
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//计划送达时间
	SendAt time.Time
	//目标邮件地址
	ToEmail string
	//标题
	Title string
	//内容
	Content string
	//是否包含HTML
	IsHtml bool
}

// Send 申请发送一份邮件
func Send(args *ArgsSend) (data FieldsEmailType, err error) {
	contentType := "text"
	if args.IsHtml {
		contentType = "html"
	}
	if args.ServerID < 1 {
		args.ServerID, err = BaseConfig.GetDataInt64("EmailDefaultServerID")
		if err != nil {
			err = errors.New("email server is not exist, " + err.Error())
			return
		}
	}
	if args.SendAt.Unix() < 1 {
		args.SendAt = CoreFilter.GetNowTime()
	}
	if args.ToEmail == "" {
		err = errors.New("email is empty")
		return
	}
	var lastID int64
	lastID, err = CoreSQL.CreateOneAndID(
		Router2SystemConfig.MainDB.DB,
		"INSERT INTO core_email(org_id, server_id, create_info, send_at, title, content, content_type, to_email) VALUES(:org_id, :server_id, :create_info, :send_at, :title, :content, :content_type, :to_email)",
		map[string]interface{}{
			"org_id":       args.OrgID,
			"server_id":    args.ServerID,
			"create_info":  args.CreateInfo,
			"send_at":      args.SendAt,
			"title":        args.Title,
			"content":      args.Content,
			"content_type": contentType,
			"to_email":     args.ToEmail,
		},
	)
	if err == nil {
		data, err = GetEmailByID(&ArgsGetEmailByID{
			ID: lastID,
		})
	}
	//标记阻断器
	runBlocker.NewEdit()
	//反馈
	return
}

// ArgsDeleteEmailByID 删除消息参数
type ArgsDeleteEmailByID struct {
	//ID
	ID int64 `db:"id"`
}

// DeleteEmailByID 删除消息
func DeleteEmailByID(args *ArgsDeleteEmailByID) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "core_email", "id", args)
	return
}

// 标记发送失败并留下日志
func updateFailed(id int64, message string) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_email SET is_failed=true, fail_message=:fail_message, update_at=NOW() WHERE id=:id", map[string]interface{}{
		"fail_message": message,
		"id":           id,
	})
	return
}

// 发送一个普通邮件
func sendMail(serverData FieldsEmailServerType, data FieldsEmailType) error {
	auth := smtp.PlainAuth("", serverData.Email, serverData.Password, serverData.Host)
	to := []string{data.ToEmail}
	msg := []byte(
		"To: " + data.ToEmail + "\r\n" +
			"From: " + serverData.Email + "\r\n" +
			"Subject: " + data.Title + "\r\n" +
			"\r\n" +
			data.Content)
	err := smtp.SendMail(serverData.Host+":"+serverData.Port, auth, serverData.Email, to, msg)
	return err
}

// 发送一个ssl邮件
func sendSSLMail(serverData FieldsEmailServerType, data FieldsEmailType) error {
	header := make(map[string]string)
	header["From"] = serverData.Name + "<" + serverData.Email + ">"
	header["To"] = data.ToEmail
	header["Subject"] = data.Title
	switch data.ContentType {
	case "text":
		header["Content-Type"] = "text/plain; charset=UTF-8"
	case "html":
		header["Content-Type"] = "text/html; charset=UTF-8"
	default:
		return errors.New("content type is error")
	}
	var body string
	for k, v := range header {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body = body + "\r\n\r\n" + data.Content
	auth := smtp.PlainAuth(
		"",
		serverData.Email,
		serverData.Password,
		serverData.Host,
	)
	err := sendMailUsingTLS(
		serverData.Host+":"+serverData.Port,
		auth,
		serverData.Email,
		[]string{data.ToEmail},
		[]byte(body),
	)
	return err
}

// return a smtp client
func sendMailDial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 参考net/smtp的func SendMail()
// 使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
// len(to)>1时,to[1]开始提示是密送
func sendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := sendMailDial(addr)
	if err != nil {
		return err
	}
	defer func() {
		_ = c.Close()
	}()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
