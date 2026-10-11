package common

import (
	"fmt"
	"os"
	"time"

	"github.com/eryajf/go-ldap-admin/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局日志变量
//var Log *zap.Logger
var Log *zap.SugaredLogger

/**
 * 初始化日志
 * filename 日志文件路径
 * level 日志级别
 * maxSize 每个日志文件保存的最大尺寸 单位：M
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名
 * 由于zap不具备日志切割功能, 这里使用lumberjack配合
 */
func InitLogger() {
	now := time.Now()
	infoLogFileName := fmt.Sprintf("%s/info/%04d-%02d-%02d.log", config.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	errorLogFileName := fmt.Sprintf("%s/error/%04d-%02d-%02d.log", config.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	var coreArr []zapcore.Core

	// 获取编码器
	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 指定时间格式
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // ，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	////encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        // 显示完整文件路径
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "file",
		FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		//EncodeTime: zapcore.ISO8601TimeEncoder, // ISO8601 UTC 时间格式
		//EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		//	enc.AppendInt64(int64(d) / 1000000)
		//},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		//EncodeCaller: zapcore.FullCallerEncoder,
		//EncodeName:       nil,
		//ConsoleSeparator: "",
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 日志级别
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	})

	// 当yml配置中的等级大于Error时，lowPriority级别日志停止记录
	if config.Conf.Logs.Level >= 2 {
		lowPriority = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return false
		})
	}

	// info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogFileName,             //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    config.Conf.Logs.MaxSize,    //文件大小限制,单位MB
		MaxAge:     config.Conf.Logs.MaxAge,     //日志文件保留天数
		MaxBackups: config.Conf.Logs.MaxBackups, //最大保留日志文件数量
		LocalTime:  false,
		Compress:   config.Conf.Logs.Compress, //是否压缩处理
	})
	// 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)

	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFileName,            //日志文件存放目录
		MaxSize:    config.Conf.Logs.MaxSize,    //文件大小限制,单位MB
		MaxAge:     config.Conf.Logs.MaxAge,     //日志文件保留天数
		MaxBackups: config.Conf.Logs.MaxBackups, //最大保留日志文件数量
		LocalTime:  false,
		Compress:   config.Conf.Logs.Compress, //是否压缩处理
	})
	// 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	logger := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
	Log = logger.Sugar()
	Log.Info("zap logger initialized")
}
