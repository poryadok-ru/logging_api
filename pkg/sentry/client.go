package sentry

import (
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

	event := sentry.NewEvent()
	event.Message = msg
	event.Level = level
	event.Timestamp = createdAt

	event.Contexts = map[string]sentry.Context{
		"log": {
			"bot_id": botID,
			"status": status,
		},
	}

	// Теги: только project_code, bot_id и bot_name
	event.Tags = map[string]string{
		"project_code": projectCode,
		"bot_id":       botID,
		"bot_name":     botName,
	}

	event.Extra = map[string]interface{}{
		"bot_id":       botID,
		"project_code": projectCode,
		"bot_name":     botName,
		"log_status":   status,
		"created_at":   createdAt.Format(time.RFC3339),
	}

	sentry.CaptureEvent(event)
}
