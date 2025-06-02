package logging

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 日志配置
type Config struct {
	Level         string            `yaml:"level"`
	Dir           string            `yaml:"dir"`
	MaxSize       int               `yaml:"max_size"`
	MaxBackups    int               `yaml:"max_backups"`
	MaxAge        int               `yaml:"max_age"`
	Compress      bool              `yaml:"compress"`
	Development   bool              `yaml:"development"`
	DisableCaller bool              `yaml:"disable_caller"`
	UseLocalTime  bool              `yaml:"use_local_time"`
	UseUTCTime    bool              `yaml:"use_utc_time"`
	SplitByLevel  bool              `yaml:"split_by_level"`
	LevelFiles    map[string]string `yaml:"level_files"`
}

// NewLogger 创建日志实例
func NewLogger(cfg Config) (*zap.Logger, error) {
	// 创建日志文件目录
	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return nil, err
	}

	// 设置日志级别
	level := zapcore.InfoLevel
	if err := level.Set(cfg.Level); err != nil {
		level = zapcore.InfoLevel
		log.Printf("Invalid log level: %s", err)
	}

	// 创建编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	if cfg.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	// 设置时间格式
	if cfg.UseLocalTime {
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
	} else if cfg.UseUTCTime {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// 创建编码器
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建核心
	var core zapcore.Core

	if cfg.SplitByLevel {
		// 按级别分割日志
		cores := make([]zapcore.Core, 0, 4)

		// 获取当前日期字符串 (YYYYMMDD)
		dateStr := time.Now().Format("20060102")

		// 创建不同级别的写入器
		levelWriters := map[zapcore.Level]zapcore.WriteSyncer{
			zapcore.DebugLevel: getLevelWriter(cfg, dateStr, "debug"),
			zapcore.InfoLevel:  getLevelWriter(cfg, dateStr, "info"),
			zapcore.WarnLevel:  getLevelWriter(cfg, dateStr, "warn"),
			zapcore.ErrorLevel: getLevelWriter(cfg, dateStr, "error"),
		}

		// 为每个级别创建独立的core
		for level, writer := range levelWriters {
			levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl == level
			})
			cores = append(cores, zapcore.NewCore(fileEncoder, writer, levelEnabler))
		}

		// 如果是开发环境，还需要输出到控制台
		if cfg.Development {
			consoleWriter := zapcore.Lock(os.Stdout)
			consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, level)
			cores = append(cores, consoleCore)
		}

		// 将所有core组合在一起
		core = zapcore.NewTee(cores...)
	} else {
		// 不按级别分割，使用单个文件
		fileWriter := getFileWriter(cfg, "app.log")

		if cfg.Development {
			// 开发环境同时输出到控制台和文件
			consoleWriter := zapcore.Lock(os.Stdout)
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, consoleWriter, level),
				zapcore.NewCore(fileEncoder, fileWriter, level),
			)
		} else {
			// 生产环境只输出到文件
			core = zapcore.NewCore(fileEncoder, fileWriter, level)
		}
	}

	// 创建日志实例
	logger := zap.New(core,
		zap.WithCaller(cfg.DisableCaller),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, nil
}

// getLevelWriter 获取特定级别的写入器
func getLevelWriter(cfg Config, dateStr, levelName string) zapcore.WriteSyncer {
	// 获取该级别配置的文件名模板
	filePattern, exists := cfg.LevelFiles[levelName]
	if !exists {
		filePattern = "{date}_" + levelName + ".log"
	}

	// 替换日期占位符
	filename := strings.ReplaceAll(filePattern, "{date}", dateStr)

	// 完整文件路径
	filePath := filepath.Join(cfg.Dir, filename)

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})
}

// getFileWriter 获取通用文件写入器
func getFileWriter(cfg Config, filename string) zapcore.WriteSyncer {
	filePath := filepath.Join(cfg.Dir, filename)

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})
}
