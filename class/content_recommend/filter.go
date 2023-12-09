package ClassContentRecommend

// 在 ContentBased 结构体后添加以下方法：

// FindItems 根据条件过滤数据
/**
# 查找OrgID为"org1"且包含condition1和condition2条件的项目
items := cb.FindItems("org1", []string{"condition1", "condition2"})
fmt.Println("Filtered items:", items)

# 如果需要忽略，给与-1或空切片即可
// 忽略 OrgID，仅根据 Conditions 检索
items = cb.FindItems("-1", []string{"condition1", "condition2"})
// 忽略 Conditions，仅根据 OrgID 检索
items = cb.FindItems("org1", []string{})
// 忽略 OrgID 和 Conditions，检索所有项目
items = cb.FindItems("-1", []string{})
*/
func (c *ContentBased) FindItems(orgID string, conditions []string) []*Item {
	var filteredItems []*Item

	for _, item := range c.items {
		if orgID != "-1" && item.OrgID != orgID {
			continue
		}

		if len(conditions) > 0 {
			matched := true
			for _, condition := range conditions {
				if !contains(item.Conditions, condition) {
					matched = false
					break
				}
			}

			if !matched {
				continue
			}
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
}

// contains 添加一个辅助函数，用于检查某个元素是否在一个字符串切片中：
func contains(arr []string, item string) bool {
	for _, elem := range arr {
		if elem == item {
			return true
		}
	}
	return false
}
