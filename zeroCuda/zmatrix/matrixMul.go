package zmatrix

/*
#include <stdbool.h>
todo:
*/
import (
	"C"
)

import (
	"errors"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"time"
)

func MatrixFloat64Mul(a []float64, b []float64, c []float64, grid []int32, block []int32, index int, m, n, k int64) {
	C.matrixDoubleMul((*C.double)(&(a)[0]), (*C.double)(&(b)[0]), (*C.double)(&(c)[0]), (*C.int)(&grid[0]), (*C.int)(&block[0]), C.int(index), C.longlong(m), C.longlong(n), C.longlong(k))
}

func getCudaGridAndBlock(row, col int64) (grid, block []int32) {
	var limitBlock = int32(32)
	grid = make([]int32, 3)
	block = make([]int32, 3)
	grid[0] = 1 + (int32(row)+limitBlock-1)/limitBlock
	grid[1] = 1 + (int32(col)+limitBlock-1)/limitBlock
	grid[2] = 1
	block[0] = limitBlock
	block[1] = limitBlock
	block[2] = 1
	return
}

func GpuMatrixFloat64Mul(a *MatrixFloat64, b *MatrixFloat64, index int) (*MatrixFloat64, error) {
	if a.Col != b.Row {
		return nil, errors.New("the rows and columns are not equal")
	}
	result := NewMatrixFloat64(a.Row, b.Col)

	grid, block := getCudaGridAndBlock(a.Row, a.Col)
	MatrixFloat64Mul(a.Data, b.Data, result.Data, grid, block, index, a.Row, a.Col, b.Col)
	return &result, nil
}

func VerifyResult(a *Matrix, b *Matrix) bool {
	if len(a.Data) != len(b.Data) {
		return false
	}
	for i := range a.Data {
		if a.Data[i] != b.Data[i] {
			return false
		}
	}
	return true
}

func verifyMatrixRow(a *Matrix, b *Matrix, result *Matrix, whereRow int) (bool, error) {
	var num int
	num = whereRow * a.Row
	if result.Row*result.Col < (whereRow+1)*a.Row {
		return false, errors.New("len(result)<(whereRow+1)*a.Row")
	}
	if len(result.Data[num:b.Col+num]) != b.Col {
		return false, errors.New("len(result[num:b.Col+num])!=b.Col")
	}
	for i := whereRow; i < whereRow+1; i++ {
		for j := 0; j < b.Col; j++ {
			var r int32
			for k := 0; k < a.Col; k++ {
				r += a.Data[i*b.Row+k] * b.Data[k*a.Row+j]
			}
			if result.Data[num] != r {
				fmt.Println("cpu result=", r)
				fmt.Println("result[", num, "]=", result.Data[num], "verify=", r)
				return false, nil
			}
			num++
		}
	}
	return true, nil
}

func VerifyRandResult(a *Matrix, b *Matrix, result *Matrix) bool {
	var (
		status = make(chan bool, 50)
	)

	num := 50
	if a.Row < num {
		num = a.Row
	}

	for i := 0; i < num; i++ {
		go func() {
			rand.Seed(time.Now().UnixNano())
			randNumber := rand.Intn(a.Row)
			ok, err := verifyMatrixRow(a, b, result, randNumber)
			if err != nil {
				fmt.Println(err.Error())
				status <- false
			}
			if ok == false {
				status <- false
			} else {
				status <- true
			}
		}()
	}
	for {
		select {
		case s := <-status:
			if s == false {
				return false
			} else {
				num--
				if num <= 0 {
					return true
				}
			}

		}
	}
}

func VerifyGoNumResult(a []float64, b *MatrixFloat64, result []float64, begin int) bool {
	re := make([]float64, b.Row)
	m1 := mat.NewDense(1, int(b.Row), a)
	m2 := mat.NewDense(int(b.Row), int(b.Col), b.Data)
	m3 := mat.NewDense(1, int(b.Col), re)
	m3.Mul(m1, m2)

	for i := 0; i < len(re); i++ {
		if re[i] != result[begin+i] {
			return false
		}
	}
	return true
}

func VerifyFloat64RandResult(a *MatrixFloat64, b *MatrixFloat64, result *MatrixFloat64) bool {
	var (
		status = make(chan bool, 50)
	)
	num := int64(50)
	if a.Row < num {
		num = a.Row
	}
	for i := int64(0); i < num; i++ {
		go func() {
			rand.Seed(time.Now().UnixNano())
			randNumber := rand.Int63n(a.Row)
			ok := VerifyGoNumResult(a.Data[randNumber*b.Row:randNumber*b.Row+a.Col], b, result.Data, int(randNumber*b.Col))
			if ok == false {
				status <- false
			} else {
				status <- true
			}
		}()
	}
	for {
		select {
		case s := <-status:
			if s == false {
				return false
			} else {
				num--
				if num <= 0 {
					return true
				}
			}

		}
	}

}
