package service

import (
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/DOSNetwork/core/edge_cm/tools"
	"time"
)

//get current cwa value
func CWA() (*tools.CPUStat, *tools.MemoryStat, *tools.DiskStat, *tools.CurrentGpu) {
	var currentGPU tools.CurrentGpu
	cpu, _ := tools.GetCpuInfo(false)
	if tools.HasNvDriver() {
		gpu, _ := tools.GetGpuInfo()
		for _, gs := range gpu {
			currentGPU.Used += gs.UsedMemory
			currentGPU.Total += gs.TotalMemory
			currentGPU.Denominator += uint(100)
			currentGPU.Numerator += gs.UtilizationGPU
		}
	}

	mem, _ := tools.GetMemoryInfo()
	disk, _ := tools.GetDiskInfo(false)
	return cpu[0], mem, disk[0], &currentGPU
}

// store cwa to sqlite, return a dispatcher
func StoreCWA(jobId string) *Dispatcher {
	ch := make(chan struct{})
	go func(c chan struct{}) {
		for {
			select {
			case <-ch:
				close(c)
				return
			default:
				var imageSize int64
				cpuStat, memStat, diskStat, gpuStat := CWA()
				image, err := Client.InspectImage(jobId)
				if err != nil {
					tools.Logger.Error(fmt.Sprintf("%s", err))
				} else {
					imageSize = image.Size
				}
				timestamp := int64(time.Now().Unix())
				cwa := &tools.CWA{
					CPU:           cpuStat.UsedPercent,
					GPU:           gpuStat.Used,
					MEM:           memStat.Used,
					Disk:          diskStat.Used,
					JobId:         jobId,
					CpuCores:      cpuStat.Cores,
					MemTotal:      memStat.Total,
					GpuTotal:      gpuStat.Total,
					ImageSize:     imageSize,
					DiskTotal:     diskStat.Total,
					DetectionTime: timestamp,
					Numerator:     gpuStat.Numerator,
					Denominator:   gpuStat.Denominator,
				}
				db := tools.Insert(cwa)
				if db.Error != nil {
					tools.Logger.Error(fmt.Sprintf("insert cwa error: %s.", db.Error.Error()))
				}
			}
			time.Sleep(time.Second * time.Duration(conf.Conf.Interval.StoreCWAInterval))
		}
	}(ch)
	return &Dispatcher{C: ch}
}

// stop dispatcher
func (d *Dispatcher) Stop() {
	d.C <- struct{}{}
}

// Remove cwa via jobid
func DeleteCWA(jobId string) {
	db := tools.Delete(&tools.CWA{}, "jobid = ? ", jobId)
	if db.Error != nil {
		tools.Logger.Error(fmt.Sprintf("%s", db.Error.Error()))
	}
}

// search peak and valley values were queried by jobid
func Peak(jobId string) (*tools.PeakValue, error) {
	var peakValue tools.PeakValue
	// The structure field needs to be the same as the database field
	cond := "Max(cpu) as cpuMax, Max(gpu) as gpuMax, Max(mem) as memMax, Max(disk) as diskMax, Min(cpu) as cpuMin, " +
		"Min(gpu) as gpuMin, Min(mem) as memMin, Min(disk) as diskMin, cpu_cores as cpuCores, mem_total as memTotal, " +
		"gpu_total as gpuTotal, disk_total as diskTotal, Max(gpu_numerator) as numeratorMax, " +
		"Min(gpu_numerator) as numeratorMin, gpu_denominator as gpuDenominator, image_size as imageSize"
	db := tools.Aggregate(&peakValue, &tools.CWA{}, "jobid = ? ", cond, jobId)

	if db.Error != nil {
		tools.Logger.Error(fmt.Sprintf("%s", db.Error.Error()))
		return nil, db.Error
	}
	return &peakValue, nil
}
