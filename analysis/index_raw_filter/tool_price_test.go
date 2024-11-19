package AnalysisIndexRawFilter

import "testing"

func TestFilterPrice(t *testing.T) {
	str := "123.45万元人民币"
	result := FilterPrice(str)
	t.Log(result)
	str = ""
	result = FilterPrice(str)
	t.Log(result)
}
