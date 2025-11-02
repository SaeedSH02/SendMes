package tel

import (
	"context"
	"fmt"
	"regexp"
	middle "sendMes/Middleware"
	"sendMes/config"
	log "sendMes/logs"

	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

type StateHandler interface {
	Handle(c telebot.Context) error
}

type PhoneInputState struct{}
// type MessageInputState struct{}

var StateHandlers = map[string]StateHandler{
	"waiting_for_phone": &PhoneInputState{},
	// "waiting_for_email": &MessageInputState{},
}

func isValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^09\d{9}$`)
	return re.MatchString(phone)
}

func (s *PhoneInputState) Handle(c telebot.Context) error {
	phone := c.Message().Text

	if !isValidPhoneNumber(phone) {
		c.Send("شماره تماس معتبر نیست")
		fmt.Printf("User %v Send a unValid Number: %s", c.Chat().Username, phone)
		return nil
	}

	result, err := middle.SendOneSMS(phone)
	if err != nil {
		c.Send("خطا در ارسال پیام")
		log.Gl.Info("err while sending message"+c.Sender().Username, zap.Error(err))
		config.Rdb.Del(context.Background(), fmt.Sprintf("state:%d", c.Sender().ID))
		return nil
	}
	if int(result.MessageIds) > 1000 {
		c.Send("پیام با موفقیت ارسال شد ✅")
		fmt.Printf("User %v Send a Message to %s", c.Chat().Username, phone)
		config.Rdb.Del(context.Background(), fmt.Sprintf("state:%d", c.Sender().ID))
		return nil
	}
	if msg, exist := ErrorMessages[result.MessageIds]; exist {
		log.Gl.Warn("Result with Error :",zap.String("msg", msg))
		fmt.Printf("Result with Error: %s Code: %v", msg, result.MessageIds)
		config.Rdb.Del(context.Background(), fmt.Sprintf("state:%d", c.Sender().ID))
		return nil
	}
	_ = c.Send(fmt.Sprintf("کد خطای ناشناخته: %d ❌", result.MessageIds))
	log.Gl.Warn("Unknown result code" + result.Result, zap.Int(" code", result.MessageIds))
	config.Rdb.Del(context.Background(), fmt.Sprintf("state:%d", c.Sender().ID))
	return nil
}

// func (s *MessageInputState) Handle(c telebot.Context) error {

// 	return nil
// }
