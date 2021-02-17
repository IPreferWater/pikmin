package mqtt

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ipreferwater/pikmin/go/config"
	log "github.com/sirupsen/logrus"
)

// InitMqtt connect the mosquitto server from the robot
//
// subscribe to the topics
//
// handle the auto-reconnect
func InitMqtt() {
	config := config.Config.MQTT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	url := fmt.Sprintf("%s://%s:%d", config.Protocol, config.Host, config.Port)
	connOpts := MQTT.NewClientOptions().AddBroker(url).SetClientID(config.ClientID).SetCleanSession(true).SetAutoReconnect(true).SetConnectionLostHandler(connectionLostHandler)
	connOpts.SetUsername(config.Username)
	connOpts.SetPassword(config.Password)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		log.Infof("Connected to %s\n", url)

		topicPikmins := "/pikmins"
		if token := c.Subscribe(topicPikmins, byte(config.Qos), mqttPikmins); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		topicFood := "/food"
		if token := c.Subscribe(topicFood, byte(config.Qos), mqttTreasure); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		topicTreasure := "/treasure"
		if token := c.Subscribe(topicTreasure, byte(config.Qos), mqttFood); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	//TODO replace by cron
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("mqtt connection failed, try again in 5 seconds", token.Error())
		time.Sleep(5 * time.Second)
		InitMqtt()
	}

	<-c
}

func connectionLostHandler(client MQTT.Client, err error) {
	log.Errorf("Connection lost, reason: %v", err)
}

//example xx/xx/xx
func mqttPikmins(client MQTT.Client, message MQTT.Message) {
	payload := string(message.Payload())
	fmt.Println(payload)
}

//example xx/xx/xx
func mqttFood(client MQTT.Client, message MQTT.Message) {
	payload := string(message.Payload())
	fmt.Println(payload)
}

//example xx/xx/xx
func mqttTreasure(client MQTT.Client, message MQTT.Message) {
	payload := string(message.Payload())
	fmt.Println(payload)
}
