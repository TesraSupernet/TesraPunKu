package ztools

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, true)
	for k, v := range percent {
		fmt.Println(k, v)
	}
	return percent[0]
}

func GetCpuNum() (int, error) {
	num, err := cpu.Counts(false)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func GetFreeMemory() uint64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.Available
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}
