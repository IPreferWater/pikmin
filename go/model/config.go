package model

type Config struct {
	DatabaseConfig DatabaseConfig
	MQTTConfig MQTTConfig
	
	Grpc struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"grpc"`

	Logs struct {
		MaxSize    int `json:"maxSize"`
		MaxBackUps int `json:"maxBackUps"`
		MaxAge     int `json:"maxAge"`
	} `json:"logs"`

	
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

type MQTTConfig struct {
	Protocol     string `json:"protocol"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Qos          int    `json:"qos"`
	ClientID     string `json:"clientID"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Suppliername string `json:"suppliername"`
}
