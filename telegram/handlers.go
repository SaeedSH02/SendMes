package tel

import (
	"context"
	"fmt"
	"sendMes/config"
	logger "sendMes/logs"
	"time"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v4"
)

func startHandler(c tele.Context) error {
	userID := c.Sender().ID
	username := c.Sender().Username

	log := logger.Gl.With(
		"user_id", userID,
		"username", username,
		"handler", "startHamdler",
		"instance_id", uuid.New().String(),
	)
	stateKey := fmt.Sprintf("state:%d", userID)
	err := config.Rdb.Set(context.Background(), stateKey, "send_without_wait", 5*time.Minute).Err()
	if err != nil {
		log.Error("failed to set state in redis", "err", err)
		return c.Send("🚫 خطا در تنظیم وضعیت 🚫")
	}
	log.Info("user state set to send_without_wait")
	return c.Send("لطفا شماره تماس خود را وارد کنید :")
}

func send_custom_message(c tele.Context) error {
	userID := c.Sender().ID
	username := c.Sender().Username

	log := logger.Gl.With(
		"user_id", userID,
		"username", username,
		"handler", "startHamdler",
		"instance_id", uuid.New().String(),
	)

	stateKey := fmt.Sprintf("state:%d", c.Sender().ID)
	err := config.Rdb.Set(context.Background(), stateKey, "waiting_for_phone", 5*time.Minute).Err()
	if err != nil {
		log.Error("failed to set state in redis", "err", err)
		return c.Send("🚫 خطا در تنظیم وضعیت 🚫")
	}

	log.Info("user state set to waiting_for_phone")
	return c.Send("لطفا شماره تماس خود را وارد کنید :")
}

func onTextHandler(c tele.Context) error {
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
