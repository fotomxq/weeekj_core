package ServiceUserInfo

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"strings"
	"time"
)

// GetInfoGender 获取性别
func GetInfoGender(gender int) (genderStr string) {
	switch gender {
	case 0:
		genderStr = "男"
	case 1:
		genderStr = "女"
	default:
		genderStr = "未知"
	}
	return
}

// GetInfoAge 获取年龄
func GetInfoAge(dateOfBirth time.Time) (ageStr string) {
	age := CoreFilter.GetNowTimeCarbon().Year() - CoreFilter.GetCarbonByTime(dateOfBirth).Year()
	if age < 1 {
		age = 0
	}
	ageStr = fmt.Sprint(age)
	return
}

// GetInfoMemberList 拆分数据结构体
// 扩展信息中如果定义分组结构体，需采用本方法进行拆分数据
func GetInfoMemberList(data string) (obj [][]string) {
	objList := strings.Split(data, "&_&")
	for k := 0; k < len(objList); k++ {
		v := objList[k]
		vList := strings.Split(v, "|_|")
		obj = append(obj, vList)
	}
	return
}

// GetInfoEducationStatus 获取学历
func GetInfoEducationStatus(educationStatus int) string {
	switch educationStatus {
	case 0:
		return "无教育"
	case 1:
		return "小学"
	case 2:
		return "初中"
	case 3:
		return "高中"
	case 4:
		return "专科技校"
	case 5:
		return "本科"
	case 6:
		return "研究生"
	case 7:
		return "博士"
	case 8:
		return "博士后"
	default:
		return "未知"
	}
}

// GetInfoMaritalStatus 获取婚姻状态
func GetInfoMaritalStatus(maritalStatus bool) string {
	if maritalStatus {
		return "已婚"
	}
	return "未婚"
}
