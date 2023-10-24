package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	LogInfo    *log.Logger
	LogWarning *log.Logger
	LogError   *log.Logger
)

const logFilePath = "info.log"

func init() {
	file := getLogFile()

	LogInfo = log.New(file, "INFO:", log.LstdFlags|log.Lshortfile)
	LogWarning = log.New(file, "WARNING:", log.LstdFlags|log.Lshortfile)
	LogError = log.New(file, "ERROR:", log.LstdFlags|log.Lshortfile)
	fmt.Println("Logfile: " + logFilePath)
}

// Возвращает существующий файл, в которой можно дописывать логи или создает новый файл
func getLogFile() *os.File {
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			panic(err.Error())
		}
		return file
	}

	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		panic(err.Error())
	}
	return file
}
