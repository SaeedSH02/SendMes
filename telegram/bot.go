package tel

import (
	"log"
	"os"
	middle "sendMes/Middleware"
	logger "sendMes/logs"

	"time"

	tele "gopkg.in/telebot.v4"
)

var b *tele.Bot

func StartBot() {

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	var err error
	b, err = tele.NewBot(pref)
	if err != nil {
		logger.Gl.Error("failed to create bot", "err", err)
		log.Fatal(err)
		return
	}
	b.Use(middle.UserAllowed)
	registerHandlers(b)
	

	logger.Gl.Info("bot started!")
	b.Start()

}

func registerHandlers(b *tele.Bot) {
	b.Handle("/start", startHandler)
	b.Handle("/custom", send_custom_message)
	b.Handle(tele.OnText, onTextHandler)
	logger.Gl.Info("handlers registered")
}
