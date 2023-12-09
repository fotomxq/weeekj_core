package CoreMathLinearRegression

import (
	"gonum.org/v1/gonum/mat"
)

// LinearRegression 给与一组数字，通过线性回归预测未来N的值
func LinearRegression(data []float64, n int) []float64 {
	//如果不存在数据，则跳出无法预测
	if len(data) < 1 {
		return []float64{0}
	}
	// 将数据转换为矩阵形式
	x := mat.NewDense(len(data), 2, nil)
	y := mat.NewVecDense(len(data), data)
	for i := 0; i < len(data); i++ {
		x.Set(i, 0, float64(i))
		x.Set(i, 1, 1.0)
	}
	// 计算最小二乘法解析解
	xt := new(mat.Dense)
	xt.Mul(x.T(), x)
	if det := mat.Det(xt); det == 0 {
		return make([]float64, n)
	}
	a := new(mat.Dense)
	err := a.Inverse(xt)
	if err != nil {
		return nil
	}
	b := new(mat.VecDense)
	b.MulVec(x.T(), y)
	result := new(mat.VecDense)
	result.MulVec(a, b)
	res := make([]float64, n)
	for i := len(data); i < len(data)+n; i++ {
		res[i-len(data)] = result.At(0, 0)*float64(i) + result.At(1, 0)
	}
	return res
}
