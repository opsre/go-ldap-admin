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

// Global logger variable
//var Log *zap.Logger
var Log *zap.SugaredLogger

/**
 * Initialize logging
 * filename is the log file path
 * level is the log level
 * maxSize is the maximum size of each log file, in MB
 * maxBackups is the maximum number of log backups to keep
 * maxAge is the maximum number of days to retain log files
 * compress indicates whether to compress log files
 * serviceName is the service name
 * Zap does not support log rotation, so lumberjack is used alongside it
 */
func InitLogger() {
	now := time.Now()
	infoLogFileName := fmt.Sprintf("%s/info/%04d-%02d-%02d.log", config.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	errorLogFileName := fmt.Sprintf("%s/error/%04d-%02d-%02d.log", config.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	var coreArr []zapcore.Core

	// Get the encoder
	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Use the specified time format
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Use zapcore.CapitalLevelEncoder if color is not needed
	////encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        // Show the full file path
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
		//EncodeTime: zapcore.ISO8601TimeEncoder, // ISO8601 UTC time format
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

	// Log level
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	})

	// When the level in the YAML config is higher than Error, stop recording lowPriority logs
	if config.Conf.Logs.Level >= 2 {
		lowPriority = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return false
		})
	}

	// Info file writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogFileName,             //Log file directory; create it if it does not exist
		MaxSize:    config.Conf.Logs.MaxSize,    //File size limit, in MB
		MaxAge:     config.Conf.Logs.MaxAge,     //Number of days to retain log files
		MaxBackups: config.Conf.Logs.MaxBackups, //Maximum number of log files to retain
		LocalTime:  false,
		Compress:   config.Conf.Logs.Compress, //Whether to compress files
	})
	// The third and following arguments are log levels written to the file; ErrorLevel mode records only error-level logs
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)

	// Error file writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFileName,            //Log file directory
		MaxSize:    config.Conf.Logs.MaxSize,    //File size limit, in MB
		MaxAge:     config.Conf.Logs.MaxAge,     //Number of days to retain log files
		MaxBackups: config.Conf.Logs.MaxBackups, //Maximum number of log files to retain
		LocalTime:  false,
		Compress:   config.Conf.Logs.Compress, //Whether to compress files
	})
	// The third and following arguments are log levels written to the file; ErrorLevel mode records only error-level logs
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	logger := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
	Log = logger.Sugar()
	Log.Info("zap logger initialized")
}
