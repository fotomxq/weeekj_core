package RouterParams

import (
	"github.com/lib/pq"
	"testing"
)

func TestCheckJSON(t *testing.T) {
	type ArgsCreateOrg struct {
		//上级控制权限限制
		ParentFunc pq.StringArray `db:"parent_func" json:"parentFunc" check:"marks" empty:"true"`
	}

}