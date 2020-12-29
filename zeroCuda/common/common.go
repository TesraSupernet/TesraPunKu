package common

type GoCudaInfo struct {
	FreeMemory          int64
	TotalMemory         int64
	UsedMemory          int64
	DistHostMemory      int64
	DevicesIndex        int
	Power               int
	EngineCount         int
	MaxThreadBlock      int
	MultiProcessorCount int
	MemoryClockRate     int
	ClockRate           int
	Name                string
	MaxGridSize         [3]int
	MaxThreadsDim       [3]int
}

//矩阵的行和列
type Matrix struct {
	Row int64
	Col int64
}

//gpu性能参数
type GpuPowerParameter struct {
	Name      string
	GpuIndex  int
	A         Matrix
	B         Matrix
	UseTimes  int64
	UseGpuMem int64
}
