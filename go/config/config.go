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
	getConfigMqtt()
}

func getConfigMqtt() {
	// easy way to instance an object with nested struct
	mqtt := Config.MQTT

	mqtt.Protocol = os.Getenv("MQTT_PROTOCOL")
	mqtt.Host = os.Getenv("MQTT_HOST")
	mqtt.Port = strEnvToInt("MQTT_PORT")
	mqtt.Qos = strEnvToInt("MQTT_QOS")
	mqtt.ClientID = os.Getenv("MQTT_CLIENT_ID")
	mqtt.Username = os.Getenv("MQTT_USERNAME")
	mqtt.Password = os.Getenv("MQTT_PASSWORD")

	Config.MQTT = mqtt
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
