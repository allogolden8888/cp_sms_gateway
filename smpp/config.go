package smpp

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	To       string `json:"to"`
	Message  string `json:"message"`
	BindType string `json:"bindType"`
	Encoding string `json:"encoding"`
	Validity string `json:"validity"`
	Register int    `json:"register"`
	Priority int    `json:"priority"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(f, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
