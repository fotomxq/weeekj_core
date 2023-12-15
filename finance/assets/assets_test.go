package FinanceAssets

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	"testing"
)

var (
	newAssetsData FieldsAssets
)

func TestInit5(t *testing.T) {
	TestInit(t)
	TestCreateProduct(t)
}

func TestSetAssets(t *testing.T) {
	var err error
	newAssetsData, err = SetAssets(&ArgsSetAssets{
		OrgID:     orgID,
		BindID:    orgBindID,
		UserID:    userID,
		ProductID: newProductData.ID,
		Count:     3,
		Des:       "测试",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(newAssetsData)
	}
	newAssetsData, err = SetAssets(&ArgsSetAssets{
		OrgID:     orgID,
		BindID:    orgBindID,
		UserID:    userID,
		ProductID: newProductData.ID,
		Count:     3,
		Des:       "测试",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(newAssetsData)
		if newAssetsData.Count != 6 {
			t.Error("not 6")
		}
	}
}

func TestGetAssetsList(t *testing.T) {
	data, count, err := GetAssetsList(&ArgsGetAssetsList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:     orgID,
		UserID:    0,
		ProductID: 0,
		IsRemove:  false,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
	}
}

func TestGetAssetsByID(t *testing.T) {
	var err error
	newAssetsData, err = GetAssetsByID(&ArgsGetAssetsByID{
		ID:    newAssetsData.ID,
		OrgID: orgID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(newAssetsData)
	}
}

func TestClearAssetsByProductID(t *testing.T) {
	var err error
	err = ClearAssetsByProductID(&ArgsClearAssetsByProductID{
		OrgID:     orgID,
		ProductID: newProductData.ID,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClearAssetsByOrgID(t *testing.T) {
	var err error
	err = ClearAssetsByOrgID(&ArgsClearAssetsByOrgID{
		OrgID: orgID,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClear5(t *testing.T) {
	TestDeleteProduct(t)
	TestClear(t)
}
