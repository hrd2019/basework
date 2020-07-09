package log

import (
	cfg "github.com/fuloge/basework/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	core zapcore.Core
	//Loger *ZapLog
	ZapLoger *zap.Logger
)

//type ZapLog struct {
//	Logger *zap.Logger
//}

func init() {
	logfile := cfg.EnvConfig.Log.Logfile

	hook := lumberjack.Logger{
		Filename:   logfile, // 日志文件路径
		MaxSize:    128,     // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 300,     // 日志文件最多保存多少个备份
		MaxAge:     120,     // 文件最多保存多少天
		Compress:   true,    // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	field := zap.Fields(zap.String("serviceName", cfg.EnvConfig.Authkey.Subject))
	// 构造日志
	ZapLoger = zap.New(core, caller, development, field)
	ZapLoger.Info("log 初始化成功")

	//Loger = &ZapLog{
	//	Logger: ZapLoger,
	//}
}

func New() (logger *zap.Logger) {
	return ZapLoger
}

//func (z *ZapLog) Info(methodName, msg string, err error) {
//	z.Logger.Info(msg, zap.String(methodName, ""))
//}
//
//func (z *ZapLog) Error(methodName, msg string, err error) {
//	z.Logger.Error(msg, zap.String(methodName, err.Error()))
//}
//
//func (z *ZapLog) Debug(methodName, msg string, err error) {
//	z.Logger.Debug(msg, zap.String(methodName, err.Error()))
//}
//
//func (z *ZapLog) Fatal(methodName, msg string, err error) {
//	z.Logger.Fatal(msg, zap.String(methodName, err.Error()))
//}
