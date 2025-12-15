package handler

import (
	"context"
	"fmt"
	"regexp"
	"sendMes/config"
	logger "sendMes/logs"
	models "sendMes/models"

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
		logger.Gl.Warn(
			"invalid phone number received",
			"user_id", c.Sender().ID,
			"invalid_phone", phone,
		)
		return nil
	}
	bodyMessage := models.OrginalMsg

	result, err := SendMessage(phone, bodyMessage)
	if err != nil {
		c.Send("خطا در ارسال پیام")
		logger.Gl.Warn(
			"err while sending message",
			"sender", c.Sender().ID,
			"error", err,
		)
		clearState(userID)
		return nil
	}
	if int(result.MessageIds) > models.SuccessThreshold {
		c.Send("پیام با موفقیت ارسال شد ✅")
		logger.Gl.Info(
			"message sent successfully",
			"sender", c.Sender().ID,
			"phone_suffix", phone[len(phone)-4:],
		)
		clearState(userID)
		return nil
	} else if msg, exist := models.ErrorMessages[result.MessageIds]; exist {
		logger.Gl.Warn(
			"message send failed",
			"error_message", msg,
			"entered_number", phone,
			"status_code", result.MessageIds,
		)
		clearState(userID)
		return nil
	}
	_ = c.Send(fmt.Sprintf("کد خطای ناشناخته: %d ❌", result.MessageIds))
	logger.Gl.Warn(
		"message send failed with unknown error code",
		"raw_result", result.Result,
		"code", result.MessageIds,
	)
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
		logger.Gl.Warn(
			"invalid phone number received",
			"userID", userID,
			"username", username,
			"phone", phone,
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

	log := logger.Gl.With(
		"userID", userID,
		"username", username,
	)

	defer clearState(userID)

	phone, err := config.Rdb.Get(ctx, fmt.Sprintf("phone:%d", userID)).Result()
	if err != nil {
		c.Send("خطایی رخ داده است لطفا بعدا تلاش کنید")
		log.Error(
			"failed to get phone from redis",
			"error", err,
			"redis_key", fmt.Sprintf("phone:%d", userID),
			"step", "get_phone",
		)
		return nil
	}

	result, err := SendMessage(phone, message)
	if err != nil {
		c.Send("خطا در ارسال پیام")
		log.Error(
			"failed to send message",
			"err", err,
			"phone", phone,
			"step", "send_message",
		)
	}

	if int(result.MessageIds) > models.SuccessThreshold {
		c.Send("پیام با موفقیت ارسال شد ✅")
		log.Info(
			"message sent successfully",
			"step", "send_message",
			"phone_suffix", phone[len(phone)-4:],
		)
		return nil
	} else if msg, exist := models.ErrorMessages[result.MessageIds]; exist {
		log.Warn(
			"message send failed with known error code",
			"error_message", msg,
			"entered_number", phone,
			"status_code", result.MessageIds,
		)
		return nil
	} else {
		_ = c.Send(" خطای ناشناخته ❌")
		log.Warn(
			"message send failed with unknown error code",
			"raw_result", result.Result,
			"status_code", result.MessageIds,
		)
	}

	return nil
}
