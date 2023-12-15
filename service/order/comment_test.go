package ServiceOrder

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitComment(t *testing.T) {
	TestInit(t)
	TestCreate(t)
	TestGetList(t)
	//提交和审核订单
	TestUpdatePost(t)
	TestUpdateAudit(t)
}

func TestUpdateCommentBuyer(t *testing.T) {
	err := UpdateCommentBuyer(&ArgsUpdateCommentBuyer{
		ID:     newOrderData.ID,
		UserID: newOrderData.UserID,
		Des:    "测试comment_buyer",
		GoodFrom: CoreSQLFrom.FieldsFrom{
			System: "mall",
			ID:     123,
			Mark:   "",
			Name:   "",
		},
		CommentBuyerID: 1,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateCommentSeller(t *testing.T) {
	err := UpdateCommentSeller(&ArgsUpdateCommentSeller{
		ID:        newOrderData.ID,
		OrgID:     newOrderData.OrgID,
		OrgBindID: 123,
		Des:       "测试comment_seller",
		GoodFrom: CoreSQLFrom.FieldsFrom{
			System: "mall",
			ID:     123,
			Mark:   "",
			Name:   "",
		},
		CommentSellerID: 2,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearComment(t *testing.T) {
	TestClear(t)
}
