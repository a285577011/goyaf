package goyaf

import (
	"fmt"
	"log"
	"os"
	"time"
)

var ErrorLog *log.Logger

//日志记录
func Log(v ...interface{}) {
	systemTime := time.Now().Format("2006-01-02 15:04:05")

	fmt.Print(systemTime, " [GOYAF LOG] ")
	fmt.Println(v...)
}

//调试信息
func Debug(v ...interface{}) {
	if GetConfigByKey("debugmode") == "0" {
		return
	}
	systemTime := time.Now().Format("2006-01-02 15:04:05")

	fmt.Print(systemTime, " [GOYAF DEBUG] ")
	fmt.Println(v...)
}

//记录错误信息
func Error(v ...interface{}) {
	ErrorLog.Println(v)
}

func init() {
	ErrorLog = log.New(os.Stderr, "[GOYAF ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
