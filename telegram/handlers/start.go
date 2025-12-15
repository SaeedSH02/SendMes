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

func StartHandler(c tele.Context) error {
	userID := c.Sender().ID
	username := c.Sender().Username

	log := logger.Gl.With(
		"user_id", userID,
		"username", username,
		"handler", "startHandler",
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
