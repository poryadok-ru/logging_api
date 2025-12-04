package sentry

import (
	"fmt"
	"log"
	"time"

	"logging_api/configs"

	"github.com/getsentry/sentry-go"
)

func InitSentry(config *configs.SentryConfig) func() {
	if config.DSN == "" {
		log.Println("Sentry DSN не задан, Sentry отключен")
		return func() {}
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.DSN,
		Environment:      config.Environment,
		Release:          config.Release,
		TracesSampleRate: 1.0,
	})

	if err != nil {
		log.Printf("Ошибка инициализации Sentry: %v", err)
		return func() {}
	}

	log.Println("Sentry успешно инициализирован")
	return func() {
		sentry.Flush(time.Second * 2)
	}
}

func SendLog(botID string, projectCode string, botName string, status, msg string, createdAt time.Time) {
	if status != "Error" && status != "Critical" {
		return
	}

	var level sentry.Level
	switch status {
	case "Error":
		level = sentry.LevelError
	case "Critical":
		level = sentry.LevelFatal
	default:
		return
	}

	err := fmt.Errorf("%s", msg)

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)

		scope.SetTag("project_code", projectCode)
		scope.SetTag("bot_id", botID)
		scope.SetTag("bot_name", botName)

		scope.SetContext("log", map[string]interface{}{
			"bot_id": botID,
			"status": status,
		})

		// Дополнительная информация
		scope.SetExtra("bot_id", botID)
		scope.SetExtra("project_code", projectCode)
		scope.SetExtra("bot_name", botName)
		scope.SetExtra("log_status", status)
		scope.SetExtra("created_at", createdAt.Format(time.RFC3339))

		// Используем CaptureException для правильного отображения сообщения
		sentry.CaptureException(err)
	})
}
