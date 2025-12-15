package models

import (
	"gorm.io/gorm"
)

type ApiSender struct {
	gorm.Model
}

const (
	OrginalMsg       string = "dfsf \n لغو11"
	OwnSender        string = "30006708282828"
	SuccessThreshold        = 1000
)
