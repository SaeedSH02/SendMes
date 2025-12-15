package main

import (
	log "sendMes/logs"
	tel "sendMes/telegram"
)

func main() {

	log.Initialize()


	tel.StartBot()

}
