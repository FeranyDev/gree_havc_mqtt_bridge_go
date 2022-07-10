package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"gree_havc_mqtt_bridge_go/app"
	"net"
)

type Device app.Device

type DeviceOption struct {
	Host        net.IP
	OnStart     func(device *Device)
	OnUpdate    func(device *Device)
	OnConnected func()
}

var hvacHost = net.IPv4(192, 168, 100, 249)
var mqttBrokerUrl = "192.168.100.1"
var mqttBrokerPort = 1883
var mqttTopicPrefix = "home/greehvac1"
var mqttUser = "admin"
var mqttPass = "admin"
var mqttRetain = false

var deviceState map[string]string

var commands = app.Commands()

var client mqtt.Client

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Infof("[MQTT] Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Infof("[MQTT] Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Infof("[MQTT] Connect lost: %v\n", err)
}

var deviceOptions = DeviceOption{
	Host: hvacHost,
	OnStart: func(device *Device) {
		publishIfChanged("Temperature", device.Props[commands.Temperature.Code], "/temperature/get")
		publishIfChanged("fanSpeed", getKeyByValue(commands.FanSpeed.Value, device.Props[commands.FanSpeed.Code]), "/fanspeed/get")
		publishIfChanged("swingHor", getKeyByValue(commands.SwingHor.Value, device.Props[commands.SwingHor.Code]), "/swinghor/get")
		publishIfChanged("swingVert", getKeyByValue(commands.SwingVert.Value, device.Props[commands.SwingVert.Code]), "/swingvert/get")
		publishIfChanged("power", getKeyByValue(commands.Power.Value, device.Props[commands.Power.Code]), "/power/get")
		publishIfChanged("health", getKeyByValue(commands.Health.Value, device.Props[commands.Health.Code]), "/health/get")
		publishIfChanged("powerSave", getKeyByValue(commands.PowerSave.Value, device.Props[commands.PowerSave.Code]), "/powersave/get")
		publishIfChanged("lights", getKeyByValue(commands.Lights.Value, device.Props[commands.Lights.Code]), "/lights/get")
		publishIfChanged("quiet", getKeyByValue(commands.Quiet.Value, device.Props[commands.Quiet.Code]), "/quiet/get")
		publishIfChanged("blow", getKeyByValue(commands.Blow.Value, device.Props[commands.Blow.Code]), "/blow/get")
		publishIfChanged("air", getKeyByValue(commands.Air.Value, device.Props[commands.Air.Code]), "/air/get")
		publishIfChanged("sleep", getKeyByValue(commands.Sleep.Value, device.Props[commands.Sleep.Code]), "/sleep/get")
		publishIfChanged("turbo", getKeyByValue(commands.Turbo.Value, device.Props[commands.Turbo.Code]), "/turbo/get")
		extandeMode := "off"
		if device.Props[commands.Power.Code] == "1" {
			extandeMode = getKeyByValue(commands.Mode.Value, device.Props[commands.Mode.Code])
		}
		publishIfChanged("mode", extandeMode, "/mode/get")
	},
	OnUpdate: func(device *Device) {
		log.Infof("[UDP] Status updated on %s\n", device.Name)
	},
	OnConnected: func() {
		client.Subscribe(mqttTopicPrefix+"/temperature/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/fanspeed/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/swinghor/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/swingvert/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/power/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/health/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/powersave/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/lights/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/quiet/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/blow/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/air/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/sleep/set", 0, callBack)
		client.Subscribe(mqttTopicPrefix+"/turbo/set", 0, callBack)
	},
}

func getKeyByValue(object map[string]int, value string) string {
	return string(rune(object[value]))
}

func publishIfChanged(stateProp, newValue, mqttTopic string) {
	retain := mqttRetain
	deviceState[stateProp] = newValue
	client.Publish(mqttTopicPrefix+mqttTopic, 0, retain, newValue)
}

func callBack(_ mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	data := message.Payload()
	log.Infof("[MQTT] Received Message \"%s\" received for %s\n", data, topic)

}

func main() {
	server := mqttBrokerUrl
	port := mqttBrokerPort
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", server, port))
	opts.SetClientID(mqttTopicPrefix)
	opts.SetUsername(mqttUser)
	opts.SetPassword(mqttPass)
	opts.SetCleanSession(true)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	deviceOptions.OnConnected()

	select {}

	//client.Disconnect(250)
}
