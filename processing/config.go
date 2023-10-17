package processing

import (
	"fmt"
	"go-gpt-processing/pkg/telegram"
	"strings"
	"time"
)

func checkErr(err error) {
	var message string
	if err == nil {
		return
	}

	switch {
	// gpt не смог ответить на запрос так, как нам нужно. Поэтому пропускаем ошибку и делаем следующий запрос
	case strings.Contains(err.Error(), "Wrong answer"):
		Pause(30)
		return

	// Скорее всего проблема с интернетом или с доступом к openAI
	case strings.Contains(err.Error(), "status: code: 503"):
		message = "Не удается подключиться к openai. Проблема остановлена"

	// Произошло что-то с аккаунтом на openai или закончились деньги на токенах
	case strings.Contains(err.Error(), "status code: 429"):
		message = fmt.Sprintf("Произошла ошибка аутентификации на openai.\nПодробности:\n%s\nПрограма остановлена.", err.Error())

	// Дефолтная ошибка
	default:
		message = fmt.Sprintf("Не удалось определить тип ошибки: %s\nПрограма остановлена.", err.Error())
	}
	telegram.ErrorMessageMailing(message)
	panic(err)

}

func Pause(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
}
