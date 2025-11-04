package tel

import (
	"context"
	"fmt"
	"sendMes/config"

	tele "gopkg.in/telebot.v4"
)

func startHandler(c tele.Context) error {
	config.Rdb.Set(context.Background(), fmt.Sprintf("state:%d", c.Sender().ID), "send_without_wait", 0)
	return c.Send("لطفا شماره تماس خود را وارد کنید :")
}

func send_custom_message(c tele.Context) error {
	config.Rdb.Set(context.Background(), fmt.Sprintf("state:%d", c.Sender().ID), "waiting_for_phone", 0)
	return c.Send("لطفا شماره تماس خود را وارد کنید :")
}




func onTextHandler(c tele.Context) error {
	userID := c.Sender().ID

	state, err := config.Rdb.Get(context.Background(), fmt.Sprintf("state:%d", userID)).Result()
	if  err != nil {
		return c.Send("🚫 دستور اشتباه است 🚫")
	}
	if handler, exist := StateHandlers[state]; exist{
		return handler.Handle(c)
	}
	return nil
}

// func reply_key(c tele.Context) error {
// 	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

// 	btnSendCustom := menu.Text("ارسال گروهی پیام ✉️")
// 	btnHelp := menu.Text("راهنما ℹ️")

// 	menu.Reply(
// 		menu.Row(btnSendCustom, btnHelp),
// 		menu.Row(btnHelp),
// 	)
// 	return c.Send("یک گزینه را انتخاب کنید:", menu)
// }