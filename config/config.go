package config

import (
	"flag"
	"github.com/kkyr/fig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	usageConfig    = "use this flag for set path to configuration file"
	usagePath      = "use this flag for set path to source file"
	usageDelimiter = "use this flag for set delimiter in to path to source file"
)

// Config structure for settings of application
type Config struct {
	App struct {
		FilePath        string        `fig:"filePath" default:"../test/employees.csv"` // path to CSV file
		Delimiter       string        `fig:"delimiter" default:";"`                    // delimiter in CSV file
		TimeoutRequest  time.Duration `fig:"timeoutRequest" default:"10"`              // request timeout in seconds
		CountGoroutines int           `fig:"countGoroutines" default:"10"`             // count of goroutines
		Logger          *zap.Logger   // logger for use, don't load from configuration file
		LogLevel        int           `fig:"logLevel" default:"0"` // flag for log level (0-info, 1-warn, -1-debug, 2-error, 4-panic, 5-fatal)
		ConfigPath      string        // path to configuration file, don't load from configuration file
	} `fig:"app"`
}

// InitConfig function for initialize Config structure
func InitConfig() (*Config, error) {
	useConfig := flag.String("config", "../config/config.yml", usageConfig)
	usePath := flag.String("file", "", usagePath)
	useDelimiter := flag.String("delimiter", "", usageDelimiter)
	flag.Parse()

	var cfg = Config{}
	err := fig.Load(&cfg, fig.File(*useConfig))
	if err != nil {
		err = fig.Load(&cfg, fig.File("config.yml"))
		if err != nil {
			log.Fatalf("can't load configuration file: %s", err)
			return nil, err
		}
		*useConfig = "config.yml"
	}

	cfg.App.ConfigPath = *useConfig

	if len(*usePath) > 0 {
		cfg.App.FilePath = *usePath
	}

	if len(*useDelimiter) > 0 {
		cfg.App.Delimiter = *useDelimiter
	}

	if err := cfg.setABSPath(); err != nil {
		log.Fatalf("error on get ABS path: %v\n", err)
		return nil, err
	}

	//Set log level
	atomicLevel := zap.NewAtomicLevel()
	switch cfg.App.LogLevel {
	case 0:
		{
			atomicLevel.SetLevel(zap.InfoLevel)
		}
	case 1:
		{
			atomicLevel.SetLevel(zap.WarnLevel)
		}
	case -1:
		{
			atomicLevel.SetLevel(zap.DebugLevel)
		}
	case 2:
		{
			atomicLevel.SetLevel(zap.ErrorLevel)
		}
	case 4:
		{
			atomicLevel.SetLevel(zap.PanicLevel)
		}
	case 5:
		{
			atomicLevel.SetLevel(zap.FatalLevel)
		}
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	), zap.AddCaller())

	cfg.App.Logger = logger

	return &cfg, err
}

// setABSPath method prepare all paths to ABS
func (c *Config) setABSPath() error {
	// get absolut filepath
	sourcePath, err := filepath.Abs(c.App.FilePath)
	if err != nil {
		return err
	}
	configPath, err := filepath.Abs(c.App.ConfigPath)
	if err != nil {
		return err
	}

	c.App.FilePath = sourcePath
	c.App.ConfigPath = configPath

	return nil
}
