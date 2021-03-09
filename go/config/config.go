package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ipreferwater/pikmin/go/model"
)

// Config is the global model.Config
var Config model.Config

func InitConfig() {
	Config.MQTTConfig = getConfigMQTT()
	Config.DatabaseConfig = getConfigDatabase()
}

func getConfigMQTT() model.MQTTConfig {
	return model.MQTTConfig{
		Protocol:     os.Getenv("MQTT_PROTOCOL"),
		Host:         os.Getenv("MQTT_HOST"),
		Port:         strEnvToInt("MQTT_PORT"),
		Qos:          strEnvToInt("MQTT_QOS"),
		ClientID:     os.Getenv("MQTT_CLIENT_ID"),
		Username:     os.Getenv("MQTT_CLIENT_ID"),
		Password:     os.Getenv("MQTT_USERNAME"),
		Suppliername: os.Getenv("MQTT_PASSWORD"),
	}
}

func getConfigDatabase()model.DatabaseConfig{
	return model.DatabaseConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     strEnvToInt("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func strEnvToInt(envString string) int {
	stringValue := os.Getenv(envString)
	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		panic(err)
	}
	return intValue

}

func strEnvToTimeDuration(envString string) time.Duration {
	sValue := os.Getenv(envString)
	i, err := strconv.Atoi(sValue)
	if err != nil {
		fmt.Printf("can't parse env %s value '%s', => %s\n", envString, sValue, err)
		panic(err)
	}
	return time.Duration(i)
}
