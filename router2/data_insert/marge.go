package Router2DataInsert

import (
	"encoding/json"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	MapRoom "github.com/fotomxq/weeekj_core/v5/map/room"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	ServiceCompany "github.com/fotomxq/weeekj_core/v5/service/company"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/tidwall/gjson"
)

// AutoMargeRel 映射方法
func AutoMargeRel(rawData interface{}, needMarks []string) (resultData []byte) {
	rawByte, err := json.Marshal(rawData)
	if err != nil {
		return
	}
	return AutoMarge(rawByte, needMarks)
}

//AutoMarge 自动化填充常用数据方法
/**
可识别的ID有：
- userID -> userName / userPhone / userAvatar
- orgID -> orgName / orgLogo
- orgBindID -> orgBindName / orgBindPhone
- deviceID -> deviceName
- roomID -> roomName
- coverFileID -> coverFileURL
- desFiles -> desFileURLs
- companyID -> companyName
needMarks:
"userName", "orgName", "orgBindName", "deviceName", "roomName"
"userPhone", "orgBindPhone"
"userAvatar", "orgLogo"
"coverFileID", "desFiles"
"companyName"
*/
func AutoMarge(rawByte []byte, needMarks []string) (resultData []byte) {
	//预构建数据
	var userID, orgID, orgBindID, deviceID, roomID, coverFileID, companyID int64
	var desFiles []int64
	for _, v := range needMarks {
		switch v {
		case "userName":
			if userID < 1 {
				userID = gjson.GetBytes(rawByte, "userID").Int()
			}
		case "userPhone":
			if userID < 1 {
				userID = gjson.GetBytes(rawByte, "userID").Int()
			}
		case "userAvatar":
			if userID < 1 {
				userID = gjson.GetBytes(rawByte, "userID").Int()
			}
		case "orgName":
			if orgID < 1 {
				orgID = gjson.GetBytes(rawByte, "orgID").Int()
			}
		case "orgLogo":
			if orgID < 1 {
				orgID = gjson.GetBytes(rawByte, "orgID").Int()
			}
		case "orgBindName":
			if orgBindID < 1 {
				orgBindID = gjson.GetBytes(rawByte, "orgBindID").Int()
			}
		case "orgBindPhone":
			if orgBindID < 1 {
				orgBindID = gjson.GetBytes(rawByte, "orgBindID").Int()
			}
		case "deviceName":
			if deviceID < 1 {
				deviceID = gjson.GetBytes(rawByte, "deviceID").Int()
			}
		case "roomName":
			if roomID < 1 {
				deviceID = gjson.GetBytes(rawByte, "roomID").Int()
			}
		case "coverFileID":
			if coverFileID < 1 {
				coverFileID = gjson.GetBytes(rawByte, "coverFileID").Int()
			}
		case "desFiles":
			if len(desFiles) < 1 {
				desFiles2 := gjson.GetBytes(rawByte, "desFiles").Array()
				for _, v2 := range desFiles2 {
					desFiles = append(desFiles, v2.Int())
				}
			}
		case "companyName":
			if companyID < 1 {
				companyID = gjson.GetBytes(rawByte, "companyID").Int()
			}
		}
	}
	//查询数据包
	var userInfo UserCore.FieldsUserType
	if userID > 0 {
		userInfo, _ = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
			ID:    userID,
			OrgID: -1,
		})
	}
	var orgData OrgCore.FieldsOrg
	if orgID > 0 {
		orgData, _ = OrgCore.GetOrg(&OrgCore.ArgsGetOrg{
			ID: orgID,
		})
	}
	var bindData OrgCore.FieldsBind
	if orgBindID > 0 {
		bindData, _ = OrgCore.GetBind(&OrgCore.ArgsGetBind{
			ID: orgBindID,
		})
	}
	var deviceData IOTDevice.FieldsDevice
	if orgBindID > 0 {
		deviceData, _ = IOTDevice.GetDeviceByID(&IOTDevice.ArgsGetDeviceByID{
			ID: deviceID,
		})
	}
	var roomNames map[int64]string
	if roomID > 0 {
		roomNames, _ = MapRoom.GetRoomsName(&MapRoom.ArgsGetRooms{
			IDs:        []int64{roomID},
			HaveRemove: true,
		})
	}
	//解析json结构体
	var jsonData map[string]interface{}
	err := json.Unmarshal(rawByte, &jsonData)
	if err != nil {
		return
	}
	//构建数据
	for _, v := range needMarks {
		switch v {
		case "userName":
			jsonData["userName"] = userInfo.Name
		case "userPhone":
			jsonData["userPhone"] = userInfo.Phone
		case "userAvatar":
			jsonData["userAvatar"] = getURLByFileID(userInfo.Avatar)
		case "orgName":
			jsonData["orgName"] = orgData.Name
		case "orgLogo":
			coverFileID, _ := OrgCore.Config.GetConfigValInt64(&ClassConfig.ArgsGetConfig{
				BindID:    orgID,
				Mark:      "CoverFileID",
				VisitType: "admin",
			})
			jsonData["orgLogo"] = getURLByFileID(coverFileID)
		case "orgBindName":
			jsonData["orgBindName"] = bindData.Name
		case "orgBindPhone":
			jsonData["orgBindPhone"] = bindData.Params.GetValNoBool("phone")
		case "deviceName":
			jsonData["deviceName"] = deviceData.Name
		case "roomName":
			jsonData["roomName"] = CoreFilter.GetMapKey(roomID, roomNames)
		case "coverFileID":
			jsonData["coverFileURL"] = getURLByFileID(coverFileID)
		case "desFiles":
			var urls []string
			for _, v2 := range desFiles {
				urls = append(urls, getURLByFileID(v2))
			}
			jsonData["desFileURLs"] = urls
		case "companyName":
			jsonData["companyName"] = ServiceCompany.GetCompanyName(companyID)
		}
	}
	//制作json包
	resultData, _ = json.Marshal(jsonData)
	//反馈
	return
}
