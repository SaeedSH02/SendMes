package handler

import (
	"context"
	"fmt"
	"sendMes/config"
	logger "sendMes/logs"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v4"
)

func OnTextHandler(c tele.Context) error {
	userID := c.Sender().ID
	username := c.Sender().Username

	log := logger.Gl.With(
		"user_id", userID,
		"username", username,
		"handler", "startHamdler",
		"instance_id", uuid.New().String(),
	)

	stateKey := fmt.Sprintf("state:%d", userID)
	state, err := config.Rdb.Get(context.Background(), stateKey).Result()
	if err != nil {
		log.Error("failed to get user state from redis", "err", err)
		return c.Send("🚫 خطا در تنظیم وضعیت 🚫")
	}
	if handler, exist := StateHandlers[state]; exist {
		log.Info("handling user state", "state", state)
		return handler.Handle(c)
	} else {
		log.Warn("no handler found for user state", "state", state)
		return c.Send("🚫 وضعیت نامعتبر 🚫")
	}

}
