package log

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/unknwon/goconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// error logger
var myLogger *zap.SugaredLogger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func init() {
	fileName := fmt.Sprintf("./logs/jd_seckill_%s.log", time.Now().Format("20060102"))
	hook := lumberjack.Logger{
		Filename:   fileName, // 日志文件路径
		MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     7,        // 文件最多保存多少天
		Compress:   true,     // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 时间精度？
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	lvl := "info"
	if cfg, err := goconfig.LoadConfigFile("./conf.ini"); err == nil {
		lvl, _ = cfg.GetValue("config", "log_level")
	}
	fmt.Println(lvl)
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(getLoggerLevel(lvl))

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 日志格式
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel,                                                                     // 日志级别
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.AddCallerSkip(1)

	logger := zap.New(core, caller, development)
	myLogger = logger.Sugar()
}

//兼容 log.Println [INFO]级别
func Println(args ...interface{}) {
	myLogger.Info(args...)
}

//兼容 log.Printf [INFO]级别
func Printf(template string, args ...interface{}) {
	myLogger.Infof(template, args...)
}

func Debug(args ...interface{}) {
	myLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	myLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	myLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	myLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	myLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	myLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	myLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	myLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	myLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	myLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	myLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	myLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	myLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	myLogger.Fatalf(template, args...)
}
