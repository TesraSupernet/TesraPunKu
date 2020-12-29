package tools

import "fmt"

// GPU stat information struct
type GPUStat struct {
	Index             uint   `json:"index"`
	Model             string `json:"model"`
	Power             uint   `json:"power"`
	Temperature       uint   `json:"temperature"`
	TotalMemory       uint64 `json:"totalMemory"`
	UsedMemory        uint64 `json:"usedMemory"`
	FreeMemory        uint64 `json:"freeMemory"`
	UtilizationGPU    uint   `json:"gpuUtilization"`
	UtilizationMemory uint   `json:"memoryUtilization"`
}

// CPU stat  information struct
type CPUStat struct {
	Cores       int     `json:"cores"`
	Index       int     `json:"index"`
	UsedPercent float64 `json:"usedPercent"`
}

// Memory stat information struct
type MemoryStat struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// Disk stat information struct
type DiskStat struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Path        string  `json:"path"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// Hardware stat information struct
type HardwareStat struct {
	MemoryStat *MemoryStat `json:"memory"`
	DickStat   []*DiskStat `json:"disk"`
	GPUStat    []*GPUStat  `json:"gpu"`
	CPUStat    []*CPUStat  `json:"cpu"`
	Timestamp  int64       `json:"timestamp"`
	Mac        string      `json:"mac"`
}

// Hardware generator
// get hardware information e.g: <- hardware.C
type HardwareGenerate struct {
	C <-chan *HardwareStat
}

// cpu,gpu peak value
type PeakValue struct {
	GpuMax       uint    `gorm:"column:gpuMax" json:"gpuMax"`
	GpuMin       uint    `gorm:"column:gpuMin" json:"gpuMin"`
	CpuMax       float64 `gorm:"column:cpuMax" json:"cpuMax"`
	CpuMin       float64 `gorm:"column:cpuMin" json:"cpuMin"`
	MemMax       float64 `gorm:"column:memMax" json:"memMax"`
	MemMin       float64 `gorm:"column:memMin" json:"memMin"`
	DiskMax      float64 `gorm:"column:diskMax" json:"diskMax"`
	DiskMin      float64 `gorm:"column:diskMin" json:"diskMin"`
	CpuCores     int     `gorm:"column:cpuCores" json:"cpuCores"`
	MemTotal     uint64  `gorm:"column:memTotal" json:"memTotal"`
	GpuTotal     uint64  `gorm:"column:gpuTotal" json:"gpuTotal"`
	ImageSize    int64   `gorm:"column:imageSize" json:"imageSize"`
	DiskTotal    uint64  `gorm:"column:diskTotal" json:"diskTotal"`
	Denominator  uint    `gorm:"column:gpuDenominator" json:"gpuDenominator"`
	NumeratorMax uint    `gorm:"column:numeratorMax" json:"gpuNumeratorMax"`
	NumeratorMin uint    `gorm:"column:numeratorMin" json:"gpuNumeratorMin"`
}

////////////////////////
// job status
///////////////////////
type JobStatus struct {
	Msg       string     `json:"msg"`
	Mac       string     `json:"mac"`
	Code      int        `json:"code"`
	Status    int        `json:"status"`
	TaskId    string     `json:"taskId"`
	JobId     string     `json:"jobId"`
	PeakInfo  *PeakValue `json:"peakInfo"`
	Timestamp int64      `json:"timestamp"`
}

////////////////////////
// job info
///////////////////////
type EdgeCmJob struct {
	ID               uint   `gorm:"primary_key AUTO_INCREMENT;not null;column:id" json:"-"`
	Mac              string `gorm:"-" json:"mac"`              // mac 地址
	RDS              string `gorm:"column:rds" json:"rds"`     // 结果数据存放路径
	DDM              string `gorm:"column:ddm" json:"ddm"`     // 拆分类型
	DDS              string `gorm:"column:dds" json:"dds"`     // 数据集文件地址
	NCP              string `gorm:"column:ncp" json:"ncp"`     // 用户程序名称
	Type             string `gorm:"column:type" json:"type"`   // 任务类型
	Index            int    `gorm:"column:index" json:"index"` // 索引
	Count            int    `gorm:"column:count" json:"count"` // 拆分数量
	JobId            string `gorm:"column:jobid;not null;unique" json:"jobId"`
	TaskId           string `gorm:"column:taskid;not null;" json:"taskId"`
	NcpArgs          string `gorm:"column:ncp_args" json:"ncpArgs"` // 用户程序参数
	JobStatus        int    `gorm:"column:job_status;default:0" json:"-"`
	Dockerfile       string `gorm:"column:dockerfile;not null" json:"dockerfile"`
	NcpCommand       string `gorm:"ncp_command" json:"ncpCommand"`              // 用户启动命令
	DsmCommand       string `gorm:"dsm_command" json:"dsmCommand"`              // 数据拆分启动命令
	RdsFileLocal     string `gorm:"rds_file_local" json:"rds_file_local"`       // 结果集本地输出路径
	SdsFileLocal     string `gorm:"sds_file_local" json:"sds_file_local"`       // 数据集本地拆分路径
	ReceiveTimestamp int64  `gorm:"receive_timestamp" json:"receive_timestamp"` // 接收到任务的时间戳
}

func (e *EdgeCmJob) String() string {
	return fmt.Sprintf("jobid=%s;status=%d", e.JobId, e.JobStatus)
}

////////////////////////
// CWA info
///////////////////////
type CWA struct {
	ID            uint    `gorm:"primary_key AUTO_INCREMENT;not null;column:id" json:"-"`
	CPU           float64 `gorm:"column:cpu" json:"cpu"`
	GPU           uint64  `gorm:"column:gpu" json:"gpu"`
	MEM           uint64  `gorm:"column:mem" json:"mem"`
	Disk          uint64  `gorm:"column:disk" json:"disk"`
	JobId         string  `gorm:"column:jobid" json:"jobid"`
	CpuCores      int     `gorm:"column:cpu_cores" json:"cpu_cores"`
	MemTotal      uint64  `gorm:"column:mem_total" json:"mem_total"`
	GpuTotal      uint64  `gorm:"column:gpu_total" json:"gpu_total"`
	ImageSize     int64   `gorm:"column:image_size" json:"image_size"`
	DiskTotal     uint64  `gorm:"column:disk_total" json:"disk_total"`
	DetectionTime int64   `gorm:"column:detection_time" json:"detection_time"`
	Numerator     uint    `gorm:"column:gpu_numerator" json:"gpu_numerator"`
	Denominator   uint    `gorm:"column:gpu_denominator" json:"gpu_denominator"`
}

////////////////////////
// current gpu stat count
///////////////////////
type CurrentGpu struct {
	Used        uint64
	Total       uint64
	Numerator   uint
	Denominator uint
}

////////////////////////
// failed message
///////////////////////
type FailedMessage struct {
	ID      uint   `gorm:"primary_key AUTO_INCREMENT;not null;column:id" json:"-"`
	Message []byte `gorm:"column:message" json:"message"`
	Topic   string `gorm:"column:topic" json:"topic"`
}
