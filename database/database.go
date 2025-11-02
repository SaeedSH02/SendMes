package db

import (
	log "sendMes/logs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConn() error {

	dsn := "host=localhost user=postgres password=123 dbname=goTodo port=54321 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Gl.Fatal("Faild to connect to Database!....")
		return err
	}

	pg, err := db.DB()
	if err != nil {
		log.Gl.Fatal("Faild to connect to Database!....")
		return err
	}

	if err := pg.Ping(); err != nil {
		log.Gl.Panic("DB Connection is not established")
		return err
	}

	log.Gl.Info("DB Connected...")
	db.AutoMigrate()
	log.Gl.Info("DB Migrated...")
	return nil
}
