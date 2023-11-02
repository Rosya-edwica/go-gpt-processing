package processing

import (
	"fmt"
	"go-gpt-processing/pkg/logger"
	"go-gpt-processing/pkg/telegram"
	"log"
	"strings"
	"time"
)

const processingPrefix = "processing: "

// TODO: Заполнить структуру, перенести вопросы в конфиг и добавить методы к промту
type Promt struct {
	PositionAbout string
}

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

	case strings.Contains(err.Error(), "context deadline exceeded"):
		Pause(30)
		fmt.Println(err.Error(), "Программа продолжит выполнение через 30 секунд")
		return
	case strings.Contains(err.Error(), "GPT не знает что ответить"):
		Pause(30)
		fmt.Println(err.Error(), "Программа продолжит выполнение через 30 секунд")
		return
	case strings.Contains(err.Error(), "invalid character '<'"):
		Pause(30)
		fmt.Println(err.Error(), "Программа продолжит выполнение через 30 секунд")

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

	logger.LogError.Println(processingPrefix + message)
	telegram.ErrorMessageMailing(message)
	log.Fatal(message)

}

func Pause(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
}
