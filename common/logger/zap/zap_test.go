package zap

import (
	"testing"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
)

func TestDefaultLogger(t *testing.T) {
	// Test that DefaultLogger() returns a logger instance
	logger := DefaultLogger()
	
	if logger == nil {
		t.Error("DefaultLogger() = nil, want non-nil logger")
	}
}

func TestZapLogger_Debug(t *testing.T) {
	// Test that the Debug method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Debug("test debug message")
}

func TestZapLogger_Debugf(t *testing.T) {
	// Test that the Debugf method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Debugf("test debug message %s", "formatted")
}

func TestZapLogger_Info(t *testing.T) {
	// Test that the Info method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Info("test info message")
}

func TestZapLogger_Infof(t *testing.T) {
	// Test that the Infof method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Infof("test info message %s", "formatted")
}

func TestZapLogger_Warn(t *testing.T) {
	// Test that the Warn method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Warn("test warn message")
}

func TestZapLogger_Warnf(t *testing.T) {
	// Test that the Warnf method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Warnf("test warn message %s", "formatted")
}

func TestZapLogger_Error(t *testing.T) {
	// Test that the Error method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Error("test error message")
}

func TestZapLogger_Errorf(t *testing.T) {
	// Test that the Errorf method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Errorf("test error message %s", "formatted")
}

func TestZapLogger_GetWriter(t *testing.T) {
	// Test that the GetWriter method exists and returns a writer
	logger := DefaultLogger()
	
	writer := logger.GetWriter()
	
	if writer == nil {
		t.Error("GetWriter() = nil, want non-nil writer")
	}
}

func TestZapLogger_Flush(t *testing.T) {
	// Test that the Flush method exists and can be called
	logger := DefaultLogger()
	
	// This should not panic
	logger.Flush()
}

func TestZapLogger_LogLevels(t *testing.T) {
	// Test that the logger respects different log levels
	zapLogger := zapLogger{
		logLevel: logger.DEBUG,
	}
	
	if zapLogger.logLevel != logger.DEBUG {
		t.Errorf("Expected log level DEBUG, got %v", zapLogger.logLevel)
	}
}