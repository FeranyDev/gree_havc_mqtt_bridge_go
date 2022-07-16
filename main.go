package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"gree_havc_mqtt_bridge_go/app"
	"gree_havc_mqtt_bridge_go/config"
	"net"
	"os"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Infof("[MQTT] Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Infof("[MQTT] Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Infof("[MQTT] Connect lost: %v\n", err)
}

var reconnectHandler mqtt.ReconnectHandler = func(client mqtt.Client, options *mqtt.ClientOptions) {
	log.Infof("[MQTT] Reconnected")
}

func main() {

	c := &config.Config{}
	c.Gree = make([]config.Gree, 1)

	var ipStr string
	flag.StringVar(&ipStr, "hvac-host", "", "Device IP Address")
	flag.IntVar(&c.Gree[0].Port, "hvac-port", 7000, "Device Port")

	flag.StringVar(&c.Mqtt.Havc.Host, "mqtt-broker-url", "", "MQTT Broker URL")
	flag.IntVar(&c.Mqtt.Havc.Port, "mqtt-broker-port", 1883, "MQTT Broker Port")
	flag.StringVar(&c.Mqtt.Havc.Username, "mqtt-username", "admin", "MQTT User")
	flag.StringVar(&c.Mqtt.Havc.Password, "mqtt-password", "admin", "MQTT Password")
	flag.BoolVar(&c.Mqtt.Havc.Retain, "mqtt-retain", false, "MQTT Retain")

	flag.StringVar(&c.Gree[0].HavcTopic, "mqtt-topic-prefix", "home/greehvac", "MQTT Topic Prefix")

	flag.StringVar(&c.Mqtt.Bemfa.Host, "bemfa-broker-url", "bemfa.com", "BEMFA Broker URL")
	flag.IntVar(&c.Mqtt.Bemfa.Port, "bemfa-broker-port", 9503, "BEMFA Broker Port")
	flag.StringVar(&c.Mqtt.Bemfa.ClientID, "bemfa-client-id", "", "BEMFA Client ID")
	flag.BoolVar(&c.Mqtt.Bemfa.Retain, "bemfa-retain", false, "BEMFA Retain")
	flag.BoolVar(&c.Mqtt.Bemfa.Tls, "bemfa-tls", true, "BEMFA TLS")

	flag.StringVar(&c.Gree[0].BemfaTopic, "bemfa-topic", "", "BEMFA Topic")

	var configPath string
	flag.StringVar(&configPath, "c", "", "Config Path")

	flag.Parse()

	if (ipStr == "" && (c.Mqtt.Havc.Host == "" || (c.Mqtt.Bemfa.ClientID == "" && c.Gree[0].BemfaTopic == ""))) && configPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if configPath != "" {
		c = config.GetConfig(configPath)
	} else {
		c.Gree[0].Host = net.ParseIP(ipStr)
	}

	var mqttClient mqtt.Client
	var bemfaClient mqtt.Client
	if c.Mqtt.Havc.Host != "" {
		mqttOpts := mqtt.NewClientOptions()
		mqttOpts.SetUsername(c.Mqtt.Havc.Username)
		mqttOpts.SetPassword(c.Mqtt.Havc.Password)
		mqttOpts.SetCleanSession(true)
		mqttOpts.SetDefaultPublishHandler(messagePubHandler)
		mqttOpts.SetReconnectingHandler(reconnectHandler)
		if c.Mqtt.Havc.Tls {
			mqttOpts.SetTLSConfig(&tls.Config{InsecureSkipVerify: false})
			mqttOpts.AddBroker(fmt.Sprintf("mqtts://%s:%d", c.Mqtt.Havc.Host, c.Mqtt.Havc.Port))
		} else {
			mqttOpts.AddBroker(fmt.Sprintf("mqtt://%s:%d", c.Mqtt.Havc.Host, c.Mqtt.Havc.Port))
		}
		mqttOpts.OnConnect = connectHandler
		mqttOpts.OnConnectionLost = connectLostHandler
		mqttClient = mqtt.NewClient(mqttOpts)

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	if c.Mqtt.Bemfa.ClientID != "" {
		bemfaOpts := mqtt.NewClientOptions()
		if c.Mqtt.Bemfa.Tls {
			bemfaOpts.SetTLSConfig(&tls.Config{InsecureSkipVerify: false})
			bemfaOpts.AddBroker(fmt.Sprintf("mqtts://%s:%d", c.Mqtt.Bemfa.Host, c.Mqtt.Bemfa.Port))
		} else {
			bemfaOpts.AddBroker(fmt.Sprintf("mqtt://%s:%d", c.Mqtt.Bemfa.Host, c.Mqtt.Bemfa.Port))
		}
		bemfaOpts.SetCleanSession(true)
		bemfaOpts.SetDefaultPublishHandler(messagePubHandler)
		bemfaOpts.SetClientID(c.Mqtt.Bemfa.ClientID)
		bemfaOpts.SetReconnectingHandler(reconnectHandler)

		bemfaOpts.OnConnect = connectHandler
		bemfaOpts.OnConnectionLost = connectLostHandler
		bemfaClient = mqtt.NewClient(bemfaOpts)

		if token := bemfaClient.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	for i, gree := range c.Gree {
		log.Infof("[APP:%d] Connecting to %s:%d", i+1, gree.Host, gree.Port)
		app.Start(mqttClient, bemfaClient, &gree)
	}

	select {}

}
