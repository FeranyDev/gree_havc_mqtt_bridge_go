package main

import (
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"gree_havc_mqtt_bridge_go/app"
	"net"
	"os"
	"strconv"
)

var hvacHost net.IP
var mqttBrokerUrl string
var mqttBrokerPort int
var mqttTopicPrefix string
var mqttUser string
var mqttPass string
var mqttRetain bool

var deviceState map[string]string
var commands = app.Commands()
var client mqtt.Client
var udpClient *app.DeviceFactory

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Infof("[MQTT] Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Infof("[MQTT] Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Infof("[MQTT] Connect lost: %v\n", err)
}

var deviceOptions = app.DeviceFactory{
	Host: hvacHost,
	OnStatus: func(device *app.Device) {
		publishIfChanged("Temperature", strconv.Itoa(device.Props[commands.Temperature.Code]), "/temperature/get")
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
		if device.Props[commands.Power.Code] == 1 {
			extandeMode = getKeyByValue(commands.Mode.Value, device.Props[commands.Mode.Code])
		}
		publishIfChanged("mode", extandeMode, "/mode/get")
	},
	OnUpdate: func(device *app.Device) {
		//log.Infof("[MQTT] Device Status Update: %v\n", device)
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
		client.Subscribe(mqttTopicPrefix+"/mode/set", 0, callBack)
	},
}

func getKeyByValue(m map[string]int, value int) (key string) {
	for k, v := range m {
		if v == value {
			key = k
			return
		}
	}
	log.Errorf("[UDP] Key not found for value: %d\n", value)
	return
}

func publishIfChanged(stateProp string, newValue string, mqttTopic string) {
	retain := mqttRetain
	if deviceState[stateProp] != newValue {
		deviceState[stateProp] = newValue
		client.Publish(mqttTopicPrefix+mqttTopic, 0, retain, newValue)
	}
}

func callBack(_ mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	data := string(message.Payload())
	log.Infof("[MQTT] Received Message \"%s\" received for %s\n", data, topic)
	switch topic {
	case mqttTopicPrefix + "/temperature/set":
		float, err := strconv.ParseFloat(data, 32)
		if err != nil {
			log.Errorf("[MQTT] Error Temperature %s to int\n", data)
		}
		udpClient.SetTemperature(int(float), 0)
		return
	case mqttTopicPrefix + "/mode/set":
		if data == "off" {
			udpClient.SetPower(false)
		} else {
			if deviceState["power"] == "off" {
				udpClient.SetPower(true)
			}
			udpClient.SetMode(commands.Mode.Value[data])
		}
		return
	case mqttTopicPrefix + "/fanspeed/set":
		udpClient.SetFanSpeed(commands.FanSpeed.Value[data])
		return
	case mqttTopicPrefix + "/swinghor/set":
		udpClient.SetSwingHor(commands.SwingHor.Value[data])
		return
	case mqttTopicPrefix + "/swingvert/set":
		udpClient.SetSwingVert(commands.SwingVert.Value[data])
		return
	case mqttTopicPrefix + "/power/set":
		udpClient.SetPower(data == "1")
		return
	case mqttTopicPrefix + "/health/set":
		udpClient.SetHealthMode(data == "1")
		return
	case mqttTopicPrefix + "/powersave/set":
		udpClient.SetPowerSave(data == "1")
		return
	case mqttTopicPrefix + "/lights/set":
		udpClient.SetLights(data == "1")
		return
	case mqttTopicPrefix + "/quiet/set":
		intVar, err := strconv.Atoi(data)
		if err != nil {
			log.Errorf("[MQTT] Error Quiet Mode %s to int\n", data)
			return
		}
		udpClient.SetQuietMode(intVar)
		return
	case mqttTopicPrefix + "/blow/set":
		udpClient.SetBlow(data == "1")
		return
	case mqttTopicPrefix + "/air/set":
		udpClient.SetAir(data == "1")
		return
	case mqttTopicPrefix + "/sleep/set":
		udpClient.SetSleepMode(data == "1")
		return
	case mqttTopicPrefix + "/turbo/set":
		udpClient.SetTurbo(data == "1")
		return
	}
}

func main() {

	var ipStr string
	flag.StringVar(&ipStr, "DIR", "", "Device IP Address")
	flag.StringVar(&mqttBrokerUrl, "MBU", "", "MQTT Broker URL")
	flag.IntVar(&mqttBrokerPort, "MBP", 1883, "MQTT Broker Port")
	flag.StringVar(&mqttTopicPrefix, "MTP", "home/greehvac", "MQTT Topic Prefix")
	flag.StringVar(&mqttUser, "MU", "admin", "MQTT User")
	flag.StringVar(&mqttPass, "MP", "admin", "MQTT Password")
	flag.BoolVar(&mqttRetain, "MR", false, "MQTT Retain")

	flag.Parse()
	if ipStr == "" && mqttBrokerUrl == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	hvacHost = net.ParseIP(ipStr)
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
	udpClient = app.Create(deviceOptions)
	deviceState = make(map[string]string)
	udpClient.ConnectToDevice(hvacHost)

	select {}
}
