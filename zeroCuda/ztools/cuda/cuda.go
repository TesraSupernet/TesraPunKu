package cuda

/*
todo:
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
	"zeroCuda/common"
	"zeroCuda/ztools"
)

//调用c库函数获取gpu参数信息
func getCudaInfo(cuda *C.cudaInfo) bool {
	ok := C.getCudaInfo(cuda)
	return bool(ok)
}

func getDevicesCount() int {
	num := C.getDevicesNum()
	return int(num)
}

func multiDevicesDistActivityMemory(cuda []common.GoCudaInfo) {
	l := len(cuda)
	hostFreeMemory := ztools.GetFreeMemory()
	mem := int64(hostFreeMemory) / int64(l)
	for k, _ := range cuda {
		cuda[k].DistHostMemory = mem
	}
}

func InitCuda() ([]common.GoCudaInfo, error) {
	//获取设备数量
	count := getDevicesCount()
	fmt.Println("count=", count)
	if count == 0 {
		return nil, errors.New("no devices")
	}
	baseCudaInfo := make([]common.GoCudaInfo, count)
	cuda := make([]C.cudaInfo, count)
	ok := getCudaInfo(&cuda[0])
	fmt.Println("ok=", ok)
	if ok == false {
		return nil, errors.New("init error")
	}

	for i := 0; i < count; i++ {
		baseCudaInfo[i].Power = int(cuda[i].power)
		baseCudaInfo[i].DevicesIndex = int(cuda[i].devicesIndex)
		baseCudaInfo[i].Name = C.GoString(&cuda[i].name[0])
		baseCudaInfo[i].FreeMemory = int64(cuda[i].freeMemory)
		baseCudaInfo[i].TotalMemory = int64(cuda[i].totalMemory)
		baseCudaInfo[i].UsedMemory = int64(cuda[i].usedMemory)
		baseCudaInfo[i].ClockRate = int(cuda[i].clockRate)
		baseCudaInfo[i].EngineCount = int(cuda[i].engineCount)
		baseCudaInfo[i].MaxThreadBlock = int(cuda[i].maxThreadBlock)
		baseCudaInfo[i].MultiProcessorCount = int(cuda[i].multiProcessorCount)
		a := (*[3]uint32)(unsafe.Pointer(&cuda[i].maxGridSize))
		//fmt.Println("MaxGridSize=",*a)
		for j := range a {
			baseCudaInfo[i].MaxGridSize[j] = int(a[j])
		}
		b := (*[3]uint32)(unsafe.Pointer(&cuda[i].maxThreadsDim))
		//fmt.Println("MaxGridSize=",*b)
		for k := range b {
			baseCudaInfo[i].MaxThreadsDim[k] = int(b[k])
		}
	}
	multiDevicesDistActivityMemory(baseCudaInfo)

	return baseCudaInfo, nil
}
