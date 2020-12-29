package tools

import (
	"github.com/DOSNetwork/core/edge_cm/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

type Log struct {
}

var logger = Log{}

func (Log) LoadConf(conf *conf.Config) error {
	hook := lumberjack.Logger{
		Filename:   conf.Log.LogPath,       // 日志文件路径
		MaxSize:    conf.Log.LogMaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: conf.Log.LogMaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     conf.Log.LogMaxAge,     // 文件最多保存多少天
		Compress:   conf.Log.LogCompress,   // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "Logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch conf.Log.LogLevel {
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	default:
		atomicLevel.SetLevel(zap.DebugLevel)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 编码器配置NewConsoleEncoder,NewJSONEncoder
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	//filed := zap.Fields(zap.String("Edge", "CM"))
	// 构造日志
	Logger = zap.New(core, caller, development)
	return nil
}
