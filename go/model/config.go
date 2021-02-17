package model

type Config struct {

	Grpc struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"grpc"`

	Logs struct {
		MaxSize    int `json:"maxSize"`
		MaxBackUps int `json:"maxBackUps"`
		MaxAge     int `json:"maxAge"`
	} `json:"logs"`

	MQTT struct {
		Protocol     string `json:"protocol"`
		Host         string `json:"host"`
		Port         int    `json:"port"`
		Qos          int    `json:"qos"`
		ClientID     string `json:"clientID"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		Suppliername string `json:"suppliername"`
	} `json:"mqtt"`
}
