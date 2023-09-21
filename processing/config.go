package processing

import (
	"go-gpt-processing/pkg/logger"
	"go-gpt-processing/pkg/telegram"
	"strings"
	"time"
)

func checkErr(err error) {
	if err != nil {
		if strings.HasPrefix(err.Error(), "Неправильный ответ") {
			logger.Log.Printf("ОШИБКА: %s", err)
			return
		} else if strings.Contains(err.Error(), "context deadline exceeded") {
			logger.Log.Printf("ОШИБКА: %s", err)
			time.Sleep(time.Second * 10)
		} else if strings.Contains(err.Error(), "status code: 503") {
			logger.Log.Printf("ОШИБКА: %s", err)
			time.Sleep(time.Second * 10)
		} else {
			telegram.Mailing(err.Error())
			panic(err)
		}
	}
}

func Pause(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
}
