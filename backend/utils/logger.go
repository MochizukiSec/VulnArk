package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	// 日志文件
	logFile *os.File
	// 错误日志
	ErrorLogger *log.Logger
	// 信息日志
	InfoLogger *log.Logger
)

// 初始化日志
func init() {
	// 创建日志目录
	if err := os.MkdirAll("./logs", 0755); err != nil {
		log.Printf("创建日志目录失败: %v", err)
	}

	// 创建或打开日志文件
	currentDate := time.Now().Format("2006-01-02")
	var err error
	logFile, err = os.OpenFile(fmt.Sprintf("./logs/app-%s.log", currentDate), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("打开日志文件失败: %v", err)
	}

	// 初始化日志记录器
	ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)

	// 同时输出到控制台
	ErrorLogger.SetOutput(os.Stdout)
	InfoLogger.SetOutput(os.Stdout)
}

// LogError 记录错误日志
func LogError(msg string) {
	ErrorLogger.Println(msg)
}

// LogInfo 记录信息日志
func LogInfo(msg string) {
	InfoLogger.Println(msg)
}

// 在程序退出时关闭日志文件
func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}
