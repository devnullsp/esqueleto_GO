package main

import (
	"log"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

const (
	FICHERO_CONFIG = "./<NOMBRE>_cfg.yml
)

var (
	CONFIG      Configuracion
	sugarLogger *zap.SugaredLogger
)

type Configuracion struct {
	Log string `yaml:"log"`
}

func main() {
	loadConfig()
	InitLogger()
	defer sugarLogger.Sync()

}

//--------------------------------------------------------------------------------------------------------------------
func loadConfig() {
	f, err := os.Open(FICHERO_CONFIG)
	if err != nil {
		log.Fatalln("[loadConfig] ERROR open file:", err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&CONFIG)
	if err != nil {
		log.Fatalln("[loadConfig] ERROR decode:", err)
	}
}

func InitLogger() {
	writerSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{Filename: CONFIG.Log,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
