package poc

/*
#cgo LDFLAGS: -Wl,-rpath="./"
*/
import (
	"C"
)

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
	"zeroCuda/common"
	"zeroCuda/core"
	"zeroCuda/zmatrix"
	"zeroCuda/ztools"
	"zeroCuda/ztools/cuda"
	"zeroCuda/ztools/zeroKnow"
)

func GetGpuMaxRowAndMaxCol(gpu *common.GoCudaInfo, bytesSize int) (row, col int64) {
	var (
		gpuSteps  float64
		hostSteps float64
	)

	gpuSteps = float64(gpu.FreeMemory-200*1024*1024) / float64(bytesSize) / 3 //gpu显存矩阵能最大数
	//获取本机设备信息

	hostSteps = float64(gpu.DistHostMemory-200*1024*1024) / float64(bytesSize) / 3

	gpuMat := int64(math.Sqrt(gpuSteps))
	hostMat := int64(math.Sqrt(hostSteps))
	if gpuMat > hostMat {
		row = hostMat
		col = hostMat
	} else {
		row = gpuMat
		col = gpuMat
	}
	return
}

func GpuPowerParameterToMatrixMul(a *zmatrix.MatrixFloat64, b *zmatrix.MatrixFloat64, index int) (common.GpuPowerParameter, error) {
	var (
		gpuPar common.GpuPowerParameter
	)

	t1 := time.Now()
	c, err := zmatrix.GpuMatrixFloat64Mul(a, b, index)
	if err != nil {
		return gpuPar, err
	}
	//fmt.Println("result=",c.Data)
	t2 := ztools.GetDistrictUnixTime(t1)

	fmt.Println("matrix time=", t2)
	t3 := time.Now()
	ok := zmatrix.VerifyFloat64RandResult(a, b, c)
	if ok == false {
		return gpuPar, errors.New("verify s failure")
	}
	t4 := ztools.GetDistrictUnixTime(t3)
	fmt.Println("verify=", t4)
	gpuPar.UseTimes = t2
	gpuPar.GpuIndex = index
	gpuPar.UseTimes = t2
	gpuPar.A = common.Matrix{Row: a.Row, Col: a.Col}
	gpuPar.B = common.Matrix{Row: b.Row, Col: b.Col}
	gpuPar.UseGpuMem = (a.Row*a.Col + b.Row*b.Col + c.Row*c.Col) * 8

	defer zmatrix.FreeMatrix(*c)

	fmt.Println(gpuPar)
	return gpuPar, nil
}

func MultiGpuCalPower(gIndexs []int) (int64, *zeroKnow.ZkData, error) {
	var (
		newCudaInfo []common.GoCudaInfo
	)

	wg := sync.WaitGroup{}
	cu, err := cuda.InitCuda()
	if err != nil {
		return 0, nil, err
	}
	if len(gIndexs) == 0 {
		newCudaInfo = cu
	} else {
		if len(cu) < len(gIndexs) {
			return 0, nil, errors.New("index out of range")
		}
		for i := 0; i < len(gIndexs); i++ {
			for j := 0; j < len(cu); j++ {
				if gIndexs[i] == cu[j].DevicesIndex {
					newCudaInfo = append(newCudaInfo, cu[j])
					break
				}
			}
		}
	}

	if len(newCudaInfo) == 0 {
		return 0, nil, errors.New("input error")
	}
	wg.Add(len(newCudaInfo))
	powerArr := make([]common.GpuPowerParameter, len(newCudaInfo))
	var mutex sync.Mutex
	fmt.Println(newCudaInfo)
	for k, v := range newCudaInfo {
		go func(cuInfo common.GoCudaInfo, index int, wgc *sync.WaitGroup) {
			x, y := GetGpuMaxRowAndMaxCol(&cuInfo, 8)
			fmt.Println("x=", x, "y=", y)
			mutex.Lock()
			a := zmatrix.NewMatrixFloat64(x, y)
			a.Rand()
			b := zmatrix.NewMatrixFloat64(x, y)
			b.Rand()
			mutex.Unlock()
			gpu, err := GpuPowerParameterToMatrixMul(&a, &b, cuInfo.DevicesIndex)
			if err != nil {
				fmt.Println("name=", cuInfo.Name, "size=", a.Row*b.Col*3*8, "err=", err)
				wgc.Done()
				return
			}
			gpu.Name = cuInfo.Name
			powerArr[index] = gpu
			fmt.Println("wgc.Done()")
			defer zmatrix.FreeMatrix(a)
			defer zmatrix.FreeMatrix(b)
			wgc.Done()
		}(v, k, &wg)

	}

	wg.Wait()
	totalPower, zkd, err := core.GetMultiGpuPower(powerArr)

	fmt.Println("total power=", totalPower)
	return totalPower, zkd, nil
}

func StartPoc() (*zeroKnow.ZkData, error) {
	_, zkd, err := MultiGpuCalPower(nil)
	if err != nil {
		return nil, err
	}

	return zkd, nil
}
