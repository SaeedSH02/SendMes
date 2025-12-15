package handler

import (
	"context"
	"fmt"
	"sendMes/config"
	logger "sendMes/logs"
	"time"
	"github.com/google/uuid"
	tele "gopkg.in/telebot.v4"
)


func Send_custom_message(c tele.Context) error {
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
