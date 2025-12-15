package models

type BodyMessage struct {
	Message  string `json:"message"`
	Receptor string `json:"receptor"`
	Sender   string `json:"sender"`
}
