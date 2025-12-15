package handler

import (
	"fmt"
	logger "sendMes/logs"
	tele "gopkg.in/telebot.v4"
)

func GetID(c tele.Context) error {
	userID := c.Sender().ID
	username := c.Sender().Username

	log := logger.Gl.With(
		"user_id", userID,
		"username", username,
		"handler", "getID",
	)
	log.Info("user requested telegram id")
	return c.Send(fmt.Sprintf("Your ID: %d", c.Sender().ID))
}