package ClassContentRecommend

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math"
)

// 内容推荐系统
// 该模块采用基于内容的推荐算法，根据用户的历史行为，推荐用户可能感兴趣的内容
// 该模块占用资源较少，但是推荐的内容可能不够精准
// 基于内容的推荐模型结构体
/**
使用方法：
1. 建立ContentBased基本模型体
2. 添加商品信息
3. 训练模型
4. 推荐相似的商品
*/

type Item struct {
	ID          string
	OrgID       string
	Conditions  []string
	StringArray []string
}

type SimilarItem struct {
	ItemID     string
	Similarity float64
}

type ContentBased struct {
	items           []*Item
	contentElements map[string]map[string][]float64
}

func NewContentBased() *ContentBased {
	return &ContentBased{
		items:           []*Item{},
		contentElements: map[string]map[string][]float64{},
	}
}

func (c *ContentBased) AddItem(id string, orgID string, conditions []string, stringArray []string) {
	item := &Item{
		ID:          id,
		OrgID:       orgID,
		Conditions:  conditions,
		StringArray: stringArray,
	}

	for _, condition := range conditions {
		if _, ok := c.contentElements[condition]; !ok {
			c.contentElements[condition] = make(map[string][]float64)
		}
		c.contentElements[condition][id] = vectorize(id, condition)
	}

	c.items = append(c.items, item)
}

func (c *ContentBased) Train() error {
	contentConditions := make([]string, 0, len(c.contentElements))
	for condition := range c.contentElements {
		contentConditions = append(contentConditions, condition)
	}

	vectorizer := NewTfidfVectorizer(contentConditions)
	vectorizer.Fit()

	for i, item1 := range c.items {
		for j, item2 := range c.items {
			if i == j {
				continue
			}

			sim := 0.0
			for condition := range c.contentElements {
				if _, ok := c.contentElements[condition]; !ok {
					continue
				}

				vector1 := getAvgVector(c.contentElements[condition], item1.ID)
				vector2 := getAvgVector(c.contentElements[condition], item2.ID)
				sim += vectorizer.IDF(condition) * cosineSimilarity(vector1, vector2)
			}

			if sim > 0 {
				item1.StringArray = append(item1.StringArray, fmt.Sprintf("%s-%.4f", item2.ID, sim))
			}
		}
	}

	return nil
}

func (c *ContentBased) Recommend(id string) ([]string, error) {
	var stringArray []string
	for _, item := range c.items {
		if item.ID == id {
			stringArray = item.StringArray
			break
		}
	}

	if stringArray == nil {
		return nil, errors.New("Item not found")
	}

	return stringArray, nil
}

func vectorize(itemID string, condition string) []float64 {
	return []float64{hash(itemID, condition)}
}

func getAvgVector(items map[string][]float64, itemID string) []float64 {
	avgVector := items[itemID]
	return avgVector
}

func hash(itemID string, condition string) float64 {
	return float64(hashCode(fmt.Sprintf("%s-%s", itemID, condition)))
}

func hashCode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func cosineSimilarity(a, b []float64) float64 {
	dotProduct := 0.0
	magnitudeA := 0.0
	magnitudeB := 0.0

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
}

type TfidfVectorizer struct {
	idf map[string]float64
}

func NewTfidfVectorizer(contentConditions []string) *TfidfVectorizer {
	return &TfidfVectorizer{idf: make(map[string]float64)}
}

func (t *TfidfVectorizer) Fit() {
	// 在这里实现TF-IDF计算逻辑。此处仅为简化示例，您可能需要根据实际需求调整。
	for condition := range t.idf {
		t.idf[condition] = 1.0
	}
}

func (t *TfidfVectorizer) IDF(condition string) float64 {
	return t.idf[condition]
}
