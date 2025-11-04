package tel

import (
	"os"
	log "sendMes/logs"
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
		log.Gl.Fatal("Cant make bot")
		return
	}
	registerHandlers(b)

	log.Gl.Info("Bot Started...")
	b.Start()

}

func registerHandlers(b *tele.Bot) {
	b.Handle("/start", startHandler)
	b.Handle("/custom", send_custom_message)
	b.Handle(tele.OnText, onTextHandler)
}