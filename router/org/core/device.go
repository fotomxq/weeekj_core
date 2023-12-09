package RouterOrgCore

import (
	"context"
	"github.com/gin-gonic/gin"
)

//检查设备和组织关系
func CheckOrgDevicePermission(c *gin.Context, ctx context.Context, deviceID, groupID string, permissions []string, action []string) bool {
	//获取组织数据
	/**
	orgData := c.MustGet("OrgData").(OrgCore.FieldsOrg)
	//获取设备对应关系
	err := DeviceCore.CheckOperateByDeviceID(ctx, &DeviceCore.ArgsOperateCheck{
		DeviceID:    deviceID,
		GroupID:     groupID,
		CreateInfo:  GetDataFromByOrgNoName(&orgData),
		Permissions: permissions,
		Action:      action,
	})
	if err != nil {
		CoreLog.Error("user(", orgData.UserID, ") try visit page url: ", c.Request.RequestURI, ", but not have device operate permission, org id: ", orgData.ID, ", device id: ", deviceID, ", group id: ", groupID, ", need permissions: ", permissions, ", action: ", action, ", err: ", err)
		RouterReport.BaseError(c, "no-work-permission", "")
	}
	return err == nil
	*/
	return false
}
