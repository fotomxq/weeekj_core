package BaseFileSys

import (
	"testing"
	"time"

	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
)

var (
	testNewData FieldsFileClaimType
	isRun       = false
)

func TestInit(t *testing.T) {
	if isRun {
		return
	}
	isRun = true
	ToolsTest.Init(t)
}

func TestCreate(t *testing.T) {
	var err error
	var errCode string
	//创建文件测试
	testNewData, _, errCode, err = Create(&ArgsCreate{
		CreateIP:   "0.0.0.0",
		CreateInfo: CoreSQLFrom.FieldsFrom{System: "system", ID: 1},
		UserID:     1,
		OrgID:      0,
		FileSize:   1500,
		FileType:   "jpg",
		FileHash:   "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
		FileSrc:    "",
		ExpireAt:   time.Time{},
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "qiniu", Mark: "9c4a8d09ca3762af61e59520943dc26494f22412"},
		Infos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a1",
				Val:  "val1",
			},
		},
		ClaimInfos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a2",
				Val:  "valueA2",
			},
		},
		Des: "des content a1~2.",
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new test data: ", testNewData)
	}
}

func TestCreate2(t *testing.T) {
	newData, _, errCode, err := Create(&ArgsCreate{
		CreateIP:   "0.0.0.0",
		CreateInfo: CoreSQLFrom.FieldsFrom{System: "system", ID: 2},
		UserID:     1,
		OrgID:      0,
		FileSize:   1800,
		FileType:   "jpg",
		FileHash:   "9d969eef6eca33c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
		FileSrc:    "",
		ExpireAt:   time.Time{},
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "qiniu", Mark: "9c4a8d09ca3762af61e59520943dc26494f22412"},
		Infos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a3",
				Val:  "valueA3",
			},
		},
		ClaimInfos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a4",
				Val:  "valueA4",
			},
		},
		Des: "des content a3~4.",
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new data: ", newData)
	}
	newData, _, errCode, err = Create(&ArgsCreate{
		CreateIP:   "0.0.0.0",
		CreateInfo: CoreSQLFrom.FieldsFrom{System: "system", ID: 4},
		UserID:     1,
		OrgID:      0,
		FileSize:   1100,
		FileType:   "jpg",
		FileHash:   "9d969eef6eca33c29a3a629280e333cf0c3f5d5a86aff3ca12020c923adc6c92",
		FileSrc:    "",
		ExpireAt:   time.Time{},
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "qiniu", Mark: "9c4a8d09ca3762af61e59520943dc26494f22412"},
		Infos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a5",
				Val:  "valueA5",
			},
		},
		ClaimInfos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a6",
				Val:  "valueA6",
			},
		},
		Des: "des content a5~6.",
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new data: ", newData)
	}
	newData, _, errCode, err = Create(&ArgsCreate{
		CreateIP:   "0.0.0.0",
		CreateInfo: CoreSQLFrom.FieldsFrom{System: "system", ID: 5},
		UserID:     1,
		OrgID:      0,
		FileSize:   1600,
		FileType:   "jpg",
		FileHash:   "9d969eef6eca33c29a2a629280e333cf0c3f5d5a86a333ca12020c923adc6c92",
		FileSrc:    "",
		ExpireAt:   time.Time{},
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "qiniu", Mark: "9c2a8d09ca3762af61e59520943dc26494f22400"},
		Infos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a7",
				Val:  "valueA7",
			},
		},
		ClaimInfos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a8",
				Val:  "valueA8",
			},
		},
		Des: "des content a7~8.",
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new data: ", newData)
	}
	newData, _, errCode, err = Create(&ArgsCreate{
		CreateIP:   "0.0.0.0",
		CreateInfo: CoreSQLFrom.FieldsFrom{System: "system", ID: 7},
		UserID:     1,
		OrgID:      0,
		FileSize:   1100,
		FileType:   "jpg",
		FileHash:   "9d969eef6eca33c29a3a629280e333cf0c3f5d5a86aff3ca12020c923adc6c92",
		FileSrc:    "",
		ExpireAt:   CoreFilter.GetNowTime().Add(time.Hour * 10),
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "qiniu", Mark: "9c2a8d09ca37j2af55e59520943dc26494f22200"},
		Infos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a7",
				Val:  "valueA7",
			},
		},
		ClaimInfos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "a8",
				Val:  "valueA8",
			},
		},
		Des: "des content a9~10.",
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new data: ", newData)
	}
}

func TestClaimFile(t *testing.T) {
	newData, errCode, err := ClaimFile(&ArgsClaimFile{
		FileID:   testNewData.FileID,
		UserID:   1,
		OrgID:    0,
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 3),
		ClaimInfos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "b1",
				Val:  "valueB1",
			},
		},
		Des: "des content b1",
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new claim file: ", newData)
	}
}

func TestGetFileByID(t *testing.T) {
	t.Log("test data: ", testNewData)
	getData, err := GetFileByID(&ArgsGetFileByID{
		ID:         testNewData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data: ", getData)
	}
	getData2, err := GetFileByID(&ArgsGetFileByID{
		ID:         testNewData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{System: "system", ID: 2},
	})
	if err == nil && getData2.ID > 0 {
		t.Error("cannot permission file data")
	}
	getData, err = GetFileByID(&ArgsGetFileByID{
		ID:         testNewData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{System: getData.CreateInfo.System, ID: getData.CreateInfo.ID},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data by from: ", getData)
	}
}

func TestGetFileByIDs(t *testing.T) {
	data, dataCount, err := GetFileByIDs(&ArgsGetFileByIDs{
		IDs:        []int64{testNewData.FileID},
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data list: ", data, dataCount)
	}
}

func TestAddVisit(t *testing.T) {
	getData, err := GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: testNewData.ID,
		UserID:  1,
		OrgID:   0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("add visit data visit count: ", getData.VisitCount, ", time: ", getData.VisitLastAt, ", getData: ", getData)
	}
	err = AddVisit(&ArgsAddVisit{
		ClaimID: getData.ID,
	})
	if err != nil {
		t.Error(err)
	}
	getData2, err := GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: testNewData.ID,
		UserID:  0,
		OrgID:   0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("add visit data visit count: ", getData2.VisitCount, ", time: ", getData2.VisitLastAt)
	}
	err = AddVisit(&ArgsAddVisit{
		ClaimID: getData.ID,
	})
	if err != nil {
		t.Error(err)
	}
	getData3, err := GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: testNewData.ID,
		UserID:  0,
		OrgID:   0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("add visit data visit count: ", getData3.VisitCount, ", time: ", getData3.VisitLastAt)
	}
	getData4, err := GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: testNewData.ID,
		UserID:  0,
		OrgID:   0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("add visit data visit count: ", getData4.VisitCount, ", time: ", getData4.VisitLastAt)
	}
	getData5, err := GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: testNewData.ID,
		UserID:  0,
		OrgID:   0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("add visit data visit count: ", getData5.VisitCount, ", time: ", getData5.VisitLastAt)
	}
}

func TestGetVisitList(t *testing.T) {
	dataList, dataCount, err := GetVisitList(&ArgsGetVisitList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:  0,
		ClaimID: 0,
		FileID:  0,
		IP:      "",
		MinTime: time.Time{},
		MaxTime: time.Time{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList, dataCount)
	}
}

func TestGetFileClaimCount(t *testing.T) {
	count := GetFileClaimCount(&ArgsGetFileClaimCount{
		FileID: testNewData.FileID,
	})
	if count < 1 {
		t.Error("cannot get file claim count")
	} else {
		t.Log("get test data claim count: ", count)
	}
}

func TestGetFileList(t *testing.T) {
	dataList, count, err := GetFileList(&ArgsGetFileList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		CreateInfo:    CoreSQLFrom.FieldsFrom{},
		FromInfo:      CoreSQLFrom.FieldsFrom{},
		FileType:      "",
		FileShaSearch: "",
		Search:        "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find data list, count: ", count, ", dataList: ", dataList)
	}
	dataList, count, err = GetFileList(&ArgsGetFileList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		CreateInfo:    CoreSQLFrom.FieldsFrom{},
		FromInfo:      CoreSQLFrom.FieldsFrom{},
		FileType:      "jpg",
		FileShaSearch: "",
		Search:        "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find data list by where file type, count: ", count, ", dataList: ", dataList)
	}
	dataList, count, err = GetFileList(&ArgsGetFileList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		CreateInfo:    CoreSQLFrom.FieldsFrom{},
		FromInfo:      CoreSQLFrom.FieldsFrom{},
		FileType:      "png",
		FileShaSearch: "",
		Search:        "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find data list by where file type 2, count: ", count, ", dataList: ", dataList)
	}
	dataList, count, err = GetFileList(&ArgsGetFileList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		CreateInfo:    CoreSQLFrom.FieldsFrom{System: "system"},
		FromInfo:      CoreSQLFrom.FieldsFrom{},
		FileType:      "",
		FileShaSearch: "",
		Search:        "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find data list by where create from, count: ", count, ", dataList: ", dataList)
	}
	dataList, count, err = GetFileList(&ArgsGetFileList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true},
		CreateInfo:    CoreSQLFrom.FieldsFrom{},
		FromInfo:      CoreSQLFrom.FieldsFrom{},
		FileType:      "",
		FileShaSearch: "",
		Search:        "0.0",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find data list by where create from, count: ", count, ", dataList: ", dataList)
	}
}

func TestGetFileClaimList(t *testing.T) {
	dataList, count, err := GetFileClaimList(&ArgsGetFileClaimList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		UserID:   0,
		OrgID:    0,
		FileID:   0,
		IsPublic: true,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find file claim list, count: ", count, ", dataList count: ", len(dataList), ", dataList: ", dataList)
	}
	dataList, count, err = GetFileClaimList(&ArgsGetFileClaimList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		UserID:   0,
		OrgID:    0,
		FileID:   testNewData.FileID,
		IsPublic: true,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find file claim list by where file id, count: ", count, ", dataList count: ", len(dataList), ", dataList: ", dataList)
	}
	dataList, count, err = GetFileClaimList(&ArgsGetFileClaimList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		UserID:   0,
		OrgID:    0,
		FileID:   0,
		IsPublic: true,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find file claim list by where user, count: ", count, ", dataList count: ", len(dataList), ", dataList: ", dataList)
	}
	dataList, count, err = GetFileClaimList(&ArgsGetFileClaimList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		UserID:   0,
		OrgID:    0,
		FileID:   0,
		IsPublic: true,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find file claim list by where user and user id, count: ", count, ", dataList count: ", len(dataList), ", dataList: ", dataList)
	}
	dataList, count, err = GetFileClaimList(&ArgsGetFileClaimList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		UserID:   0,
		OrgID:    0,
		FileID:   0,
		IsPublic: true,
		Search:   "a3~4",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find file claim list by where search, count: ", count, ", dataList count: ", len(dataList), ", dataList: ", dataList)
	}
}

func TestUpdateFileInfo(t *testing.T) {
	TestInit(t)
	newData, err := GetFileByID(&ArgsGetFileByID{
		ID:         testNewData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("before data: ", newData)
		t.Log("before testNewData: ", testNewData)
	}
	if err := UpdateFileInfo(&ArgsUpdateFileInfo{
		UpdateHash: newData.UpdateHash,
		FileID:     newData.ID,
		Infos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "update1x",
				Val:  "value_update1x",
			},
		},
	}); err != nil {
		t.Error(err)
	}
	newData2, err := GetFileByID(&ArgsGetFileByID{
		ID:         testNewData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if err != nil {
		t.Error(err)
	} else {
		isFind := false
		for _, v := range newData2.Infos {
			if v.Mark == "update1x" && v.Val == "value_update1x" {
				isFind = true
			}
		}
		if !isFind {
			t.Error("cannot find new file info data by update1.")
			t.Log("after data: ", newData2)
		}
	}
}

func TestUpdateClaimInfo(t *testing.T) {
	TestInit(t)
	newData, err := GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: testNewData.ID,
		UserID:  0,
		OrgID:   0,
	})
	if err != nil {
		t.Error(err)
	}
	if err := UpdateClaimInfo(&ArgsUpdateClaimInfo{
		UpdateHash: newData.UpdateHash,
		ClaimID:    newData.ID,
		UserID:     1,
		OrgID:      0,
		ExpireAt:   CoreFilter.GetNowTime().Add(time.Hour * 3),
		ClaimInfos: []CoreSQLConfig.FieldsConfigType{{
			Mark: "update2",
			Val:  "value_update2",
		}},
		IsPublic: true,
		Des:      "update des content",
	}); err != nil {
		t.Error(err)
	}
	newData, err = GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: testNewData.ID,
		UserID:  0,
		OrgID:   0,
	})
	if err != nil {
		t.Error(err)
	} else {
		isFind := false
		for _, v := range newData.Infos {
			if v.Mark == "update2" && v.Val == "value_update2" {
				isFind = true
			}
		}
		if !isFind {
			t.Error("cannot find new file claim info data by update2, data: ", newData)
		}
		isFind = false
		if newData.Des == "update des content" {
			isFind = true
		}
		if !isFind {
			t.Error("cannot find new file claim info data by update2, data: ", newData)
		}
	}
}

func TestDeleteClaim(t *testing.T) {
	//为test data构建新的引用，检查数量，如果少了1个，则说明删除成功
	count1 := GetFileClaimCount(&ArgsGetFileClaimCount{
		FileID: testNewData.FileID,
	})
	var newData, errCode, err = ClaimFile(&ArgsClaimFile{
		FileID:   testNewData.FileID,
		UserID:   1,
		OrgID:    0,
		ExpireAt: time.Time{},
		ClaimInfos: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "d1",
				Val:  "valued1",
			},
		},
		Des: "des content d1",
	})
	if err != nil {
		t.Error(errCode, err)
	}
	count2 := GetFileClaimCount(&ArgsGetFileClaimCount{
		FileID: testNewData.FileID,
	})
	if count2 != count1+1 {
		t.Error("count1 + 1 != count2, create claim is error")
	}
	if err := DeleteClaim(&ArgsDeleteClaim{
		UpdateHash: newData.UpdateHash,
		ClaimID:    newData.ID,
		UserID:     1,
		OrgID:      0,
	}); err != nil {
		t.Error(err)
	}
	count3 := GetFileClaimCount(&ArgsGetFileClaimCount{
		FileID: testNewData.FileID,
	})
	if count1 != count3 {
		t.Error("count1 != count3, delete claim is error")
	}
}

func TestDeleteFile(t *testing.T) {
	if !isRun {
		TestInit(t)
		TestCreate(t)
	}
	newData, err := GetFileClaimByID(&ArgsGetFileClaimByID{
		ClaimID: testNewData.ID,
		UserID:  0,
		OrgID:   0,
	})
	if err != nil {
		t.Error(err)
	}
	fileData, err := GetFileByID(&ArgsGetFileByID{
		ID:         newData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("find file data: ", fileData)
	}
	if err := DeleteFile(&ArgsDeleteFile{
		UpdateHash: fileData.UpdateHash,
		FileID:     fileData.ID,
		CreateInfo: fileData.CreateInfo,
	}); err != nil {
		t.Error(err)
	}
	getData, err := GetFileByID(&ArgsGetFileByID{
		ID:         newData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if err == nil {
		t.Error("get data is success, but data is deleted, data: ", getData)
	}
}

func TestDeleteAll(t *testing.T) {
}
