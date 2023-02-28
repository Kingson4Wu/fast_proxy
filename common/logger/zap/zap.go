package zap

import (
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func DefaultLogger() logger.Logger {
	var coreArr []zapcore.Core

	// Get encoder
	encoderConfig := zap.NewProductionEncoderConfig()            // NewJSONEncoder() outputs logs in JSON format, while NewConsoleEncoder() outputs plain text format.
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Specify the time format.
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Show logs in different colors according to the severity level. If unnecessary, use zapcore.CapitalLevelEncoder instead.
	// encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        // Show the full file path.
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Log level
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // Error level
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // Info and Debug levels, Debug level is the lowest.
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	// Info file writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./work/log/info.log", // Directory to store the log files. If the folder does not exist, it will be created automatically.
		MaxSize:    2,                     // File size limit, in MB.
		MaxBackups: 100,                   // Maximum number of retained log files.
		MaxAge:     30,                    // Number of days to keep log files.
		Compress:   false,                 // Whether to compress the log files.
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority) // The third and subsequent parameters are the log levels written to the log file. ErrorLevel mode only logs errors.
	// Error file writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./work/log/error.log", // Directory to store the log files.
		MaxSize:    1,                      // File size limit, in MB.
		MaxBackups: 5,                      // Maximum number of retained log files.
		MaxAge:     30,                     // Number of days to keep log files.
		Compress:   false,                  // Whether to compress the log files.
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) // The third and subsequent parameters are the log levels written to the log file. ErrorLevel mode only logs errors.

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	log := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) // zap.AddCaller() shows the file name and line number, and can be omitted.

	return zapLogger{logger: log, logLevel: logger.INFO, w: os.Stdout}
}

type zapLogger struct {
	logger   *zap.Logger
	logLevel logger.Level
	w        io.Writer
}

func (pl zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	if pl.logLevel <= logger.DEBUG {
		pl.logger.Sugar().Debugw(msg, keysAndValues...)
	}
}

func (pl zapLogger) Debugf(template string, args ...interface{}) {
	if pl.logLevel <= logger.DEBUG {
		pl.logger.Sugar().Debugf(template, args...)
	}
}

func (pl zapLogger) Info(msg string, keysAndValues ...interface{}) {
	if pl.logLevel <= logger.INFO {
		pl.logger.Sugar().Infow(msg, keysAndValues...)
	}
}

func (pl zapLogger) Infof(template string, args ...interface{}) {
	if pl.logLevel <= logger.INFO {
		pl.logger.Sugar().Infof(template, args...)
	}
}

func (pl zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	if pl.logLevel <= logger.WARN {
		pl.logger.Sugar().Warnw(msg, keysAndValues...)
	}
}

func (pl zapLogger) Warnf(template string, args ...interface{}) {
	if pl.logLevel <= logger.WARN {
		pl.logger.Sugar().Warnf(template, args...)
	}
}

func (pl zapLogger) Error(msg string, keysAndValues ...interface{}) {
	if pl.logLevel <= logger.ERROR {
		pl.logger.Sugar().Errorw(msg, keysAndValues...)
	}

}

func (pl zapLogger) Errorf(template string, args ...interface{}) {
	if pl.logLevel <= logger.ERROR {
		pl.logger.Sugar().Errorf(template, args...)
	}
}

func (pl zapLogger) GetWriter() io.Writer {
	return pl.w
}

func (pl zapLogger) Flush() {
	pl.logger.Sync()
}
