package agent

import "time"

type Config struct {
	LogLevel        string
	ComputingPower  int
	PollingInterval time.Duration
}

func NewConfig() *Config {
	return &Config{
		LogLevel:        "debug",                        // Уровень логирования
		ComputingPower:  3,                              // Количетво горутин
		PollingInterval: time.Duration(1) * time.Second, //Частота вызова оркестратора
	}
}
