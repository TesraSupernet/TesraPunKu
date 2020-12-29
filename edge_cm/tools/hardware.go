package tools

import (
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"os/exec"
	"time"
)

func init() {
	if HasNvDriver() {
		// Initialize NVIDIA
		nvml.Init()
	}
}

// check nvidia driver
func HasNvDriver() bool {
	err := exec.Command("cat", "/proc/driver/nvidia/version").Run()
	if err != nil {
		return false
	}
	return true
}

// Initialize NVIDIA
func Init() {
	nvml.Init()
}

// Shutdown NVIDIA
func Shutdown() {
	nvml.Shutdown()
}

// get cpu information
func GetCpuInfo(percpu bool) ([]*CPUStat, error) {
	var cpuStatList []*CPUStat
	var percentage []float64
	var err error
	cores, err := cpu.Counts(false)
	if err != nil {
		Logger.Info("get cores error.")
	}
	if percpu {
		percentage, err = cpu.Percent(0, true)
	} else {
		percentage, err = cpu.Percent(0, false)
	}
	if err != nil {
		return nil, err
	}
	for idx, cpupercent := range percentage {
		var c CPUStat
		c.Index = idx
		c.UsedPercent = cpupercent
		c.Cores = cores
		cpuStatList = append(cpuStatList, &c)
	}
	return cpuStatList, nil
}

// get memory information
func GetMemoryInfo() (*MemoryStat, error) {
	var memory MemoryStat
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	memory.Total = v.Total
	memory.Free = v.Free
	memory.Used = v.Used
	memory.UsedPercent = v.UsedPercent
	return &memory, nil
}

// get disk information
func GetDiskInfo(all bool) ([]*DiskStat, error) {
	var diskStatList []*DiskStat
	var parts []disk.PartitionStat
	var err error
	if all {
		parts, err = disk.Partitions(true)
	} else {
		parts, err = disk.Partitions(false)
	}
	if err != nil {
		return nil, err
	}
	for _, part := range parts {
		var d DiskStat
		diskStat, err := disk.Usage(part.Mountpoint)
		if err != nil {
			return nil, err
		}
		d.Total = diskStat.Total
		d.Free = diskStat.Free
		d.UsedPercent = diskStat.UsedPercent
		d.Path = diskStat.Path
		d.Used = diskStat.Used
		diskStatList = append(diskStatList, &d)
	}
	return diskStatList, nil
}

// get gpu information
func GetGpuInfo() ([]*GPUStat, error) {
	//nvml.Init()
	//defer nvml.Shutdown()
	count, err := nvml.GetDeviceCount()
	if err != nil {
		return nil, err
	}

	var devices []*nvml.Device
	var gpuStatList []*GPUStat
	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	for i, device := range devices {
		var gpuStat GPUStat
		st, err := device.Status()
		if err != nil {
			return nil, err
		}
		gpuStat.Index = uint(i)
		gpuStat.Model = *device.Model
		gpuStat.Power = *st.Power
		gpuStat.Temperature = *st.Temperature
		gpuStat.TotalMemory = *device.Memory
		gpuStat.FreeMemory = *st.Memory.Global.Free
		gpuStat.UsedMemory = *st.Memory.Global.Used
		gpuStat.UtilizationGPU = *st.Utilization.GPU
		gpuStat.UtilizationMemory = *st.Utilization.Memory
		gpuStatList = append(gpuStatList, &gpuStat)
	}
	return gpuStatList, nil
}

// generate hardware information
func HardwareGenerator(d time.Duration) *HardwareGenerate {
	c := make(chan *HardwareStat, 1)
	go func(out chan<- *HardwareStat) {
		for {
			hardwareInfo := GetHardwareInfo()
			out <- hardwareInfo
			time.Sleep(d)
		}
	}(c)
	return &HardwareGenerate{C: c}
}

// get hardware information
func GetHardwareInfo() *HardwareStat {
	var hardwareStat HardwareStat
	var gpuInfo []*GPUStat
	cpuInfo, _ := GetCpuInfo(true)
	if HasNvDriver() {
		gpuInfo, _ = GetGpuInfo()
	}
	diskInfo, _ := GetDiskInfo(false)
	memoryInfo, _ := GetMemoryInfo()

	hardwareStat.CPUStat = cpuInfo
	hardwareStat.GPUStat = gpuInfo
	hardwareStat.DickStat = diskInfo
	hardwareStat.MemoryStat = memoryInfo
	hardwareStat.Mac = GetMacByName(conf.Conf.MMLServer.MAC_ADDR)
	hardwareStat.Timestamp = int64(time.Now().Unix())
	return &hardwareStat
}
