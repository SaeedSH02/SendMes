package main

import (
	db "sendMes/database"
	log "sendMes/logs"
	tel "sendMes/telegram"
)

func main() {

	log.Initialize()

	db.DbConn()

	tel.StartBot()

}
