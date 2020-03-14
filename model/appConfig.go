package model

type AppConfig struct {
	App   App   `yaml:"app"`
	Redis Redis `yaml:"redis"`
}

type App struct {
	LogFile       string `yaml:"logFile"`
	HttpPort      string `yaml:"httpPort"`
	WebSocketPort string `yaml:"webSocketPort"`
	WebSocketUrl  string `yaml:"webSocketUrl"`
	HttpUrl       string `yaml:"httpUrl"`
}

type Redis struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"DB"`
	PoolSize     int    `yaml:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns"`
}
