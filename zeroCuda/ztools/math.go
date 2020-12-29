package ztools

import "math"

/*
@Description :截断小数位数
@parameter: prec：从第几位开始截断
*/
func TrunFloat(f float64, prec int) float64 {
	x := math.Pow10(prec)
	return math.Trunc(f*x) / x
}
