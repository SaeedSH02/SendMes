package models

import (
	"gorm.io/gorm"
)

type ApiSender struct {
	gorm.Model
}

type BodyMessage struct {
	Message  string `json:"message"`
	Receptor string `json:"receptor"`
	Sender   string `json:"sender"`
}
type Results struct {
	Result     string `json:"result"`
	MessageIds int  `json:"messageids"`
}

const (
	OrginalMsg string = "dfsf \n لغو11"
	OwnSender  string = "30006708282828"
)


