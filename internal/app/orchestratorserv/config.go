package orchestratorserv

type Config struct {
	HttpAddr string
	LogLevel string
}

func NewConfig() *Config {
	return &Config{
		HttpAddr: ":8080",
		LogLevel: "debug",
	}
}
