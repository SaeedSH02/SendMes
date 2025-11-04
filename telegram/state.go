package tel

import (
	"context"
	"fmt"
	"regexp"
	middle "sendMes/Middleware"
	models "sendMes/Models"
	"sendMes/config"
	log "sendMes/logs"

	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

type StateHandler interface {
	Handle(c telebot.Context) error
}

var ctx = context.Background()

var StateHandlers = map[string]StateHandler{
	"send_without_wait":   &SendOne{},
	"waiting_for_phone":   &PhoneInputState{},
	"waiting_for_message": &MessageInputState{},
}

func isValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^09\d{9}$`)
	return re.MatchString(phone)
}

// Clear State
func clearState(userId int64) {
	config.Rdb.Del(ctx, fmt.Sprintf("state:%d", userId))
}

type SendOne struct{}

func (s *SendOne) Handle(c telebot.Context) error {
	phone := c.Message().Text
	userID := c.Sender().ID

	if !isValidPhoneNumber(phone) {
		c.Send("شماره تماس معتبر نیست")
		fmt.Printf("User %v Send a unValid Number: %s", c.Chat().Username, phone)
		return nil
	}
	bodyMessage := models.OrginalMsg

	result, err := middle.SendMessage(phone, bodyMessage)
	if err != nil {
		c.Send("خطا در ارسال پیام")
		log.Gl.Info("err while sending message"+c.Sender().Username, zap.Error(err))
		clearState(userID)
		return nil
	}
	if int(result.MessageIds) > models.SuccessThreshold {
		c.Send("پیام با موفقیت ارسال شد ✅")
		fmt.Printf("User %v Send a Message to %s", c.Chat().Username, phone)
		clearState(userID)
		return nil
	} else if msg, exist := ErrorMessages[result.MessageIds]; exist {
		log.Gl.Warn("Result with Error :", zap.String("msg", msg))
		fmt.Printf("Result with Error: %s Code: %v", msg, result.MessageIds)
		clearState(userID)
		return nil
	}
	_ = c.Send(fmt.Sprintf("کد خطای ناشناخته: %d ❌", result.MessageIds))
	log.Gl.Warn("Unknown result code"+result.Result, zap.Int(" code", result.MessageIds))
	clearState(userID)
	return nil
}

type PhoneInputState struct{}

func (s *PhoneInputState) Handle(c telebot.Context) error {
	phone := c.Message().Text
	userID := c.Sender().ID
	username := c.Sender().Username

	if !isValidPhoneNumber(phone) {
		c.Send("شماره تماس معتبر نیست")
		log.Gl.Warn("Invalid phone number received",
			zap.Int64("userID", userID),
			zap.String("username", username),
			zap.String("phone", phone),
		)
		return nil
	}

	config.Rdb.Set(ctx, fmt.Sprintf("phone:%d", userID), phone, 0)
	config.Rdb.Set(ctx, fmt.Sprintf("state:%d", userID), "waiting_for_message", 0)

	return c.Send(" لطفاً متن پیام خود را وارد کنید: ")
}

type MessageInputState struct{}

func (s *MessageInputState) Handle(c telebot.Context) error {
	message := c.Text() + "\nلغو11"
	userID := c.Sender().ID
	username := c.Sender().Username

	logger := log.Gl.With(
		zap.Int64("userID", userID),
		zap.String("username", username),
	)

	defer clearState(userID)

	phone, err := config.Rdb.Get(ctx, fmt.Sprintf("phone:%d", userID)).Result()
	if err != nil {
		c.Send("خطایی رخ داده است لطفا بعدا تلاش کنید")
		logger.Warn("Error while get phone from redis...", zap.Error(err))
		return nil
	}

	result, err := middle.SendMessage(phone, message)
	if err != nil {
		c.Send("خطا در ارسال پیام")
		logger.Error("Failed to send message via middle.SendMessage", zap.Error(err))
	}

	if int(result.MessageIds) > models.SuccessThreshold {
		c.Send("پیام با موفقیت ارسال شد ✅")
		logger.Info("Message sent successfully", zap.String("phone", phone))
		return nil
	} else if msg, exist := ErrorMessages[result.MessageIds]; exist {
		logger.Warn("Message send failed with known error code",
			zap.String("error_message", msg),
			zap.String("Entered number", phone),
			zap.Any("status_code", result.MessageIds),
		)
		return nil
	} else {
		_ = c.Send(" خطای ناشناخته ❌")
		logger.Warn("Message send failed with unknown error code",
			zap.String("raw_result", result.Result),
		)
	}

	return nil
}
