package zmatrix

/*
todo:
*/
import (
	"C"
)

import (
	"math/rand"
	"time"
	"unsafe"
)

const maxFloat64 = 15 * 1024 * 1024 * 1024
const maxInt = 20 * 1024 * 1024 * 1024

type Matrix struct {
	Row  int //行
	Col  int //列
	Data *[maxInt]int32
}
type MatrixFloat64 struct {
	Row int64 //行
	Col int64 //列
	//Data *[maxFloat64]float64
	pointer *[maxFloat64]float64
	Data    []float64
}

func NewMatrixFloat64(row, col int64) MatrixFloat64 {
	arr := C.newMatrix(C.longlong(row), C.longlong(col))
	ms := (*[maxFloat64]float64)(unsafe.Pointer(arr))
	da := ms[:row*col]
	return MatrixFloat64{
		row,
		col,
		ms,
		da,
	}
}

func FreeMatrix(mat MatrixFloat64) {
	C.freeMatrix((*C.double)(&mat.pointer[0]))
}

func (m *Matrix) Rand() {
	rand.Seed(time.Now().UnixNano())
	for i := range m.Data {
		m.Data[i] = int32(rand.Intn(10))
	}
}

func (m *MatrixFloat64) Rand() {
	rand.Seed(time.Now().UnixNano())
	var i int64
	for i = 0; i < m.Row*m.Col; i++ {
		m.Data[i] = float64(rand.Intn(10))
	}
}
