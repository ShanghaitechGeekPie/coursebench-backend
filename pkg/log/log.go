// Copyright (C) 2021-2024 ShanghaiTech GeekPie
// This file is part of CourseBench Backend.
//
// CourseBench Backend is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// CourseBench Backend is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with CourseBench Backend.  If not, see <http://www.gnu.org/licenses/>.

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
