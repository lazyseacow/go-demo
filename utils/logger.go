package utils

import (
	"os"
	"path/filepath"

	"github.com/demo/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

// InitLogger 初始化日志系统
func InitLogger() error {
	cfg := config.GetConfig().Log

	// 确保日志目录存在
	logDir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 日志级别
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// 日志编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色输出
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 文件日志编码器配置（不使用彩色）
	fileEncoderConfig := encoderConfig
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 控制台编码器（开发模式）
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 文件编码器（JSON格式，便于日志分析）
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	// 日志文件配置（支持日志轮转）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,    // MB
		MaxBackups: cfg.MaxBackups, // 保留旧文件的个数
		MaxAge:     cfg.MaxAge,     // 保留旧文件的天数
		Compress:   cfg.Compress,   // 是否压缩
		LocalTime:  true,           // 使用本地时间
	})

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 创建 Core
	core := zapcore.NewTee(
		// 文件输出（JSON格式）
		zapcore.NewCore(fileEncoder, fileWriter, level),
		// 控制台输出（彩色格式）
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
	)

	// 创建 Logger
	Logger = zap.New(core,
		zap.AddCaller(),                       // 添加调用者信息
		zap.AddCallerSkip(1),                  // 跳过一层调用栈
		zap.AddStacktrace(zapcore.ErrorLevel), // Error 级别记录堆栈
	)

	// 创建 SugaredLogger（更灵活的API）
	Sugar = Logger.Sugar()

	return nil
}

// Sync 刷新日志缓冲区
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
	if Sugar != nil {
		_ = Sugar.Sync()
	}
}

// LogDebug 调试日志
func LogDebug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// LogInfo 信息日志
func LogInfo(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// LogWarn 警告日志
func LogWarn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// LogError 错误日志
func LogError(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// LogFatal 致命错误日志
func LogFatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// LogDebugf 格式化调试日志
func LogDebugf(template string, args ...any) {
	Sugar.Debugf(template, args...)
}

// LogInfof 格式化信息日志
func LogInfof(template string, args ...any) {
	Sugar.Infof(template, args...)
}

// LogWarnf 格式化警告日志
func LogWarnf(template string, args ...any) {
	Sugar.Warnf(template, args...)
}

// LogErrorf 格式化错误日志
func LogErrorf(template string, args ...any) {
	Sugar.Errorf(template, args...)
}

// LogFatalf 格式化致命错误日志
func LogFatalf(template string, args ...any) {
	Sugar.Fatalf(template, args...)
}
