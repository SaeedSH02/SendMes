package tel

import (
	"context"
	"fmt"
	"sendMes/config"

	tele "gopkg.in/telebot.v4"
)

func startHandler(c tele.Context) error {
	config.Rdb.Set(context.Background(), fmt.Sprintf("state:%d", c.Sender().ID), "waiting_for_phone", 0)
	return c.Send("لطفا شماره تماس خود را وارد کنید :")
}

func onTextHandler(c tele.Context) error {
	userID := c.Sender().ID

	state, err := config.Rdb.Get(context.Background(), fmt.Sprintf("state:%d", userID)).Result()
	if  err != nil {
		return err
	}
	if handler, exist := StateHandlers[state]; exist{
		return handler.Handle(c)
	}
	return nil
}