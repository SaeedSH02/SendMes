package config

type Config struct {
	Redis Redis `json:"redis"`
}

type Redis struct {
	Addr     string `json:"address" default:"localhost:6379"`
	Password string `json:"password"`
	DB       int    `json:"db" default:"0"`
}
