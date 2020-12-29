package conf

////////////////////////
// topic config
///////////////////////
const (
	HardwarePushTopic   = "ec_container_node_dynamic"   // 硬件动态信息推送Topic
	JobStatusPushTopic  = "ec_container_job_status"     // Job状态推送Topic
	RunningJobPushTopic = "ec_container_job_running"    // 正在运行Job数量推送队列
	JobStatusPullTopic  = "ec_dataset_status"           // Job状态接收Topic
	JobCommandPullTopic = "ec_schedule_job_command"     // Job操作命令接收Topic
	JobInfoPullTopic    = "ec_schedule_job_distributor" // 任务文件,dockerfile文件 等信息接收队列
)

////////////////////////
// status code
///////////////////////
const (
	ImageBuildFailedCode        = 501000 //	镜像拉取失败
	ImageBuildingCode           = 501001 // 正在构建镜像
	ImageBuildSuccessCode       = 501001 //	镜像拉取成功
	StartContainerFailedCode    = 502000 //容器(job)启动失败
	StartContainerSuccessCode   = 502001 //容器(job)启动成功
	StopContainerFailedCode     = 503000 //停止容器(job)停止失败
	StopContainerSuccessCode    = 503001 //容器(job)停止成功
	RestartContainerFailedCode  = 504000 //容器(job)重启失败
	RestartContainerSuccessCode = 504001 //容器(job)重启成功
	RemoveImageFailedCode       = 505000 //	镜像删除失败
	RemoveImageSuccessCode      = 505001 //	镜像删除成功
	GetDockerfileFailed         = 506000 // mml code 获取dockerfile失败
	GetDockerfileSuccess        = 506001 // mml code 获取dockerfile成功
	ContainerExitedCode         = 507000 // 容易异常退出
)

////////////////////////
// job status
///////////////////////
const (
	NotRunning               = -1 // 还未开始
	ImageBuilding            = 10 // 正在构建镜像
	ImageBuilt               = 11 // 镜像构建成功
	ImageBuildFailure        = 12 // 镜像构建失败
	StartContainerSuccess    = 13 // 启动容器成功
	StartContainerFailure    = 14 // 启动容器失败
	Decomposing              = 20 // 数据集正在拆分
	Decomposed               = 21 //数据集拆分成功
	DecompositionFailure     = 22 //数据集拆分失败
	DataSetDownloading       = 23 // 数据集正在下载
	DataSetDownloaded        = 24 // 数据集下载完成
	DataSetDownloadFailure   = 25 // 数据集下载失败
	DataSetDecompression     = 26 //数据集解压开始
	DataSetDecompressSuccess = 27 //数据集解压成功
	DataSetDecompressFailure = 28 //数据集解压失败
	Calculating              = 30 // ncp正在计算
	Calculated               = 31 // ncp计算成功
	CalculateFailure         = 32 // ncp计算失败
	Uploading                = 40 // 正在上传结果
	Uploaded                 = 41 //上传结果成功
	UploadFailure            = 42 //上传结果失败
	Stopping                 = 50 // 停止中
	StopSuccess              = 51 // 停止成功
	StopFailure              = 52 // 停止失败
	SystemError              = 60 // 系统错误
	RestartSuccess           = 61 // 重启成功
	RestartFailure           = 62 // 重启失败
	ContainerExited          = 63 // 容器异常退出
)

////////////////////////
// job command
///////////////////////
const (
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"
	DELETE  = "delete"
)

////////////////////////
// job type
///////////////////////
const (
	AI   = "1" // ai任务
	CR   = "2" // 云渲染任务
	RL   = "3" // 增强学习
	BD   = "4" // 大数据
	VASP = "5" // vasp
	P2P  = "6" // p2p
)
