package CoreMathRFM

import "testing"

func TestCore(t *testing.T) {
	c := Core{
		weightList: []Weight{
			{Number: 0, R: 0.3, F: 0.4, M: 0.5},
		},
		rMin: 10,
		fMin: 15,
		mMin: 20,
		rMax: 100,
		fMax: 200,
		mMax: 300,
	}
	c.SetWeight([]Weight{
		{Number: 0, R: 0.3, F: 0.4, M: 0.5},
	})
	c.SetDataRange(10, 15, 20, 100, 200, 300)
	rfm := c.GetScoreByWeight(15, 16, 17, 0)
	t.Log(rfm)
}
