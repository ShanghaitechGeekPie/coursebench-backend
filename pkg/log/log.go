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

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...any) {
	log.Print(v...)
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...any) {
	log.Printf(format, v...)
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...any) {
	log.Println(v...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...any) {
	log.Fatal(v...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...any) {
	log.Fatalln(v...)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...any) {
	log.Panic(v...)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...any) {
	log.Panicf(format, v...)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...any) {
	log.Panicln(v...)
}
