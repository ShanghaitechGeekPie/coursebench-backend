package log

import (
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

type LogConfig struct {
	LogFile string `mapstructure:"log_file"`
}

var logConfig LogConfig
var file *os.File
var LogWriter io.Writer

func InitLog() {
	config := viper.Sub("log")
	if config == nil {
		log.Fatalln("Log config not found")
	}
	err := config.Unmarshal(&logConfig)
	if err != nil {
		panic(err)
	}
	file, err = os.OpenFile(logConfig.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	LogWriter = io.MultiWriter(os.Stdout, file)

	log.SetOutput(LogWriter)
	log.Println("Log initialized")
}

func WriteLog(str string) {
	log.Println(str)
}
