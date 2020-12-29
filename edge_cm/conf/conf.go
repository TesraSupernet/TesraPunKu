package conf

type LoadConf interface {
	LoadConf(conf *Config) error
}
type LogConfig struct {
	LogMaxAge     int    `yaml:"log_max_age"`     // 日志最大储存天数
	LogMaxBackups int    `yaml:"log_max_backups"` // 最大日志备份数量
	LogMaxSize    int    `yaml:"log_max_size"`    // 单个最大日志占用多少兆
	LogCompress   bool   `yaml:"log_compress"`    // 是否压缩日志
	LogLevel      string `yaml:"log_level"`       // 日志等级
	LogPath       string `yaml:"log_path"`        // 日志储存路径
}

type PulsarConfig struct {
	PulsarURL         string `yaml:"url"`                // pulsar地址
	ConnectionTimeout int    `yaml:"connection_timeout"` // pulsar客户端TCP连接超时时间
	OperationTimeout  int    `yaml:"operation_timeout"`  // pulsar客户端操作超时时间
	Token             string `yaml:"token"`
}

type IntervalConfig struct {
	PushHardwareInterval    int  `yaml:"push_hardware_interval"`    // 硬件信息推送间隔
	StoreCWAInterval        int  `yaml:"store_cwa_interval"`        // 峰值统计时间间隔
	PushRunningJobInterval  int  `yaml:"push_running_job_interval"` // 上传正在运行job间隔
	StopContainerTimeout    uint `yaml:"stop_container_timeout"`    // 停止容器超时时间
	RestartContainerTimeout uint `yaml:"restart_container_timeout"` // 重启容器超时时间
	ResendMessageInterval   int  `yaml:"resend_message_interval"`   // 监测失败消息并重新发送时间间隔
}

type ImageConfig struct {
	MaxImages         int   `yaml:"max_images"`           // job image max number
	AutoRmImgInterval int   `yaml:"auto_remove_interval"` // auto remove image interval
	MaxDealTimeout    int64 `yaml:"max_deal_timeout"`     // 接受任务到容器启动的最大处理时间
}

type InterfaceConfig struct {
	JOB_INFO_URL string `yaml:"url"` // get dockerfile url
}

type ContainerConfig struct {
	AutoRm    bool `yaml:"auto_rm"`    // auto remove container
	CheckLive int  `yaml:"check_live"` // check container live interval
}

type DBConfig struct {
	Path string `yaml:"path"`
}

type MMLServerConfig struct {
	SERVER_ADDR string `yaml:"addr"`     // mml server address
	MAC_ADDR    string `yaml:"mac_addr"` // mml server address
}

type P2pClientConfig struct {
	Enabled        bool   `yaml:"enabled"`
	DefaultVersion string `yaml:"default_version"`
}

type MinioConfig struct {
	Endpoint string `yaml:"endpoint"`
	Access   string `yaml:"access"`
	Secret   string `yaml:"secret"`
	UseSsl   bool   `yaml:"use_ssl"`
}

type Config struct {
	SeverPowerUrl string          `yaml:"sever_power_url"`
	SeverLoginUrl string          `yaml:"sever_login_url"`
	Log           LogConfig       `yaml:"log"`
	Pulsar        PulsarConfig    `yaml:"pulsar"`
	Interval      IntervalConfig  `yaml:"interval"`
	Image         ImageConfig     `yaml:"image"`
	Interface     InterfaceConfig `yaml:"interface"`
	MMLServer     MMLServerConfig `yaml:"server"`
	Container     ContainerConfig `yaml:"container"`
	DB            DBConfig        `yaml:"db"`
	P2pClient     P2pClientConfig `yaml:"p2p_client"`
	Minio         MinioConfig     `yaml:"minio"`
}

var Conf Config

func init() {
	logConfig := LogConfig{
		LogMaxAge:     7,
		LogMaxBackups: 30,
		LogMaxSize:    128,
		LogCompress:   true,
		LogLevel:      "debug",
		LogPath:       "./logs/edge_cm.log",
	}

	pulsarConfig := PulsarConfig{
		PulsarURL:         "pulsar://10.10.2.6:9001",
		ConnectionTimeout: 5,
		OperationTimeout:  5,
		Token:             "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJKb2UifQ.ipevRNuRP6HflG8cFKnmUPtypruRC4fb1DWtoLL62SY",
	}

	intervalConfig := IntervalConfig{
		PushHardwareInterval:    30,
		StoreCWAInterval:        5,
		PushRunningJobInterval:  30,
		StopContainerTimeout:    10,
		RestartContainerTimeout: 10,
		ResendMessageInterval:   5,
	}

	imageConfig := ImageConfig{
		MaxImages:         5,
		AutoRmImgInterval: 30,
		MaxDealTimeout:    int64(3600),
	}

	interfaceConfig := InterfaceConfig{
		JOB_INFO_URL: "http://10.10.1.22:8082/task/pull/",
	}

	mmlServerConfig := MMLServerConfig{
		SERVER_ADDR: ":9090",
		MAC_ADDR:    "eno1",
	}
	containerConfig := ContainerConfig{
		AutoRm:    true,
		CheckLive: 60,
	}
	dbConfig := DBConfig{
		Path: "/home/db/",
	}
	p2pClient := P2pClientConfig{
		DefaultVersion: "1.0.0",
		Enabled:        false,
	}
	minio := MinioConfig{
		Endpoint: "http://113.204.194.92:19000",
		Access:   "tesra",
		Secret:   "tesra123456",
		UseSsl:   false,
	}
	Conf.Log = logConfig
	Conf.Pulsar = pulsarConfig
	Conf.Interval = intervalConfig
	Conf.Image = imageConfig
	Conf.Interface = interfaceConfig
	Conf.MMLServer = mmlServerConfig
	Conf.Container = containerConfig
	Conf.DB = dbConfig
	Conf.P2pClient = p2pClient
	Conf.Minio = minio
}
