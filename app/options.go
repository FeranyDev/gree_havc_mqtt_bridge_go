package app

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/feranydev/gree_havc_mqtt_bridge_go/bemfa"
	"github.com/feranydev/gree_havc_mqtt_bridge_go/config"
	"github.com/feranydev/gree_havc_mqtt_bridge_go/gree"
	"github.com/labstack/gommon/log"
)

type AppOptions struct {
	HvacHost        net.IP
	MqttTopicPrefix string
	BemfaTopic      string
	MqttRetain      bool
	DeviceState     map[string]string
	Commands        gree.Command
	Command         func()
	mqttClient      mqtt.Client
	BemfaClient     mqtt.Client
	UdpClient       *gree.DeviceFactory
	BemfaState      *Bemfa
}

type Bemfa struct {
	Power       bool
	Mode        int
	Temperature int
	FanSpeed    int
	SwingHor    int
	SwingVert   int
}

func (options *AppOptions) UpdateStatusToMqtt(device *gree.Device) {
	if options.mqttClient != nil && options.MqttTopicPrefix != "" {
		options.publishIfChanged("Temperature", strconv.Itoa(device.Props[options.Commands.Temperature.Code]), "/temperature/get")
		options.publishIfChanged("fanSpeed", getKeyByValue(options.Commands.FanSpeed.Value, device.Props[options.Commands.FanSpeed.Code]), "/fanspeed/get")
		options.publishIfChanged("swingHor", getKeyByValue(options.Commands.SwingHor.Value, device.Props[options.Commands.SwingHor.Code]), "/swinghor/get")
		options.publishIfChanged("swingVert", getKeyByValue(options.Commands.SwingVert.Value, device.Props[options.Commands.SwingVert.Code]), "/swingvert/get")
		options.publishIfChanged("power", getKeyByValue(options.Commands.Power.Value, device.Props[options.Commands.Power.Code]), "/power/get")
		options.publishIfChanged("health", getKeyByValue(options.Commands.Health.Value, device.Props[options.Commands.Health.Code]), "/health/get")
		options.publishIfChanged("powerSave", getKeyByValue(options.Commands.PowerSave.Value, device.Props[options.Commands.PowerSave.Code]), "/powersave/get")
		options.publishIfChanged("lights", getKeyByValue(options.Commands.Lights.Value, device.Props[options.Commands.Lights.Code]), "/lights/get")
		options.publishIfChanged("quiet", getKeyByValue(options.Commands.Quiet.Value, device.Props[options.Commands.Quiet.Code]), "/quiet/get")
		options.publishIfChanged("blow", getKeyByValue(options.Commands.Blow.Value, device.Props[options.Commands.Blow.Code]), "/blow/get")
		options.publishIfChanged("air", getKeyByValue(options.Commands.Air.Value, device.Props[options.Commands.Air.Code]), "/air/get")
		options.publishIfChanged("sleep", getKeyByValue(options.Commands.Sleep.Value, device.Props[options.Commands.Sleep.Code]), "/sleep/get")
		options.publishIfChanged("turbo", getKeyByValue(options.Commands.Turbo.Value, device.Props[options.Commands.Turbo.Code]), "/turbo/get")

		extandeMode := "off"
		if device.Props[options.Commands.Power.Code] == 1 {
			extandeMode = getKeyByValue(options.Commands.Mode.Value, device.Props[options.Commands.Mode.Code])
		}
		options.publishIfChanged("mode", extandeMode, "/mode/get")
	}
	if options.BemfaClient != nil && options.BemfaTopic != "" {
		go options.bemfaGet(device)
	}
}
func (options *AppOptions) OnUpdate(_ *gree.Device) {
	// log.Infof("[MQTT] Device Status Update: %v\n", device)
}
func (options *AppOptions) OnConnected() {
	time.Sleep(time.Second * 2)
	mqttTopicPrefix := options.MqttTopicPrefix
	callBack := options.callBack
	if options.mqttClient != nil && options.MqttTopicPrefix != "" {
		options.mqttClient.Subscribe(mqttTopicPrefix+"/temperature/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/fanspeed/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/swinghor/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/swingvert/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/power/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/health/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/powersave/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/lights/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/quiet/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/blow/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/air/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/sleep/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/turbo/set", 0, callBack)
		options.mqttClient.Subscribe(mqttTopicPrefix+"/mode/set", 0, callBack)
	}
	if options.BemfaClient != nil && options.BemfaTopic != "" {
		options.BemfaClient.Subscribe(options.BemfaTopic, 0, callBack)
	}
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

func (options *AppOptions) publishIfChanged(stateProp string, newValue string, mqttTopic string) {
	if options.DeviceState[stateProp] != newValue {
		options.DeviceState[stateProp] = newValue
		options.mqttClient.Publish(options.MqttTopicPrefix+mqttTopic, 0, options.MqttRetain, newValue)
	}
}

func (options *AppOptions) callBack(_ mqtt.Client, message mqtt.Message) {
	data := string(message.Payload())

	log.Infof("[MQTT] Received Message \"%s\" received for %s\n", data, message.Topic())
	switch message.Topic() {
	case options.MqttTopicPrefix + "/temperature/set":
		float, err := strconv.ParseFloat(data, 32)
		if err != nil {
			log.Errorf("[MQTT] Error Temperature %s to int\n", data)
		}
		options.UdpClient.SetTemperature(int(float), 0)
		return
	case options.MqttTopicPrefix + "/mode/set":
		if data == "off" {
			options.UdpClient.SetPower(false)
		} else {
			if options.DeviceState["power"] == "off" {
				options.UdpClient.SetPower(true)
			}
			options.UdpClient.SetMode(options.Commands.Mode.Value[data])
		}
		return
	case options.MqttTopicPrefix + "/fanspeed/set":
		options.UdpClient.SetFanSpeed(options.Commands.FanSpeed.Value[data])
		return
	case options.MqttTopicPrefix + "/swinghor/set":
		options.UdpClient.SetSwingHor(options.Commands.SwingHor.Value[data])
		return
	case options.MqttTopicPrefix + "/swingvert/set":
		options.UdpClient.SetSwingVert(options.Commands.SwingVert.Value[data])
		return
	case options.MqttTopicPrefix + "/power/set":
		options.UdpClient.SetPower(data == "1")
		return
	case options.MqttTopicPrefix + "/health/set":
		options.UdpClient.SetHealthMode(data == "1")
		return
	case options.MqttTopicPrefix + "/powersave/set":
		options.UdpClient.SetPowerSave(data == "1")
		return
	case options.MqttTopicPrefix + "/lights/set":
		options.UdpClient.SetLights(data == "1")
		return
	case options.MqttTopicPrefix + "/quiet/set":
		intVar, err := strconv.Atoi(data)
		if err != nil {
			log.Errorf("[MQTT] Error Quiet Mode %s to int\n", data)
			return
		}
		options.UdpClient.SetQuietMode(intVar)
		return
	case options.MqttTopicPrefix + "/blow/set":
		options.UdpClient.SetBlow(data == "1")
		return
	case options.MqttTopicPrefix + "/air/set":
		options.UdpClient.SetAir(data == "1")
		return
	case options.MqttTopicPrefix + "/sleep/set":
		options.UdpClient.SetSleepMode(data == "1")
		return
	case options.MqttTopicPrefix + "/turbo/set":
		options.UdpClient.SetTurbo(data == "1")
		return
	case options.BemfaTopic + "":
		go options.bemfaSet(data)
	}
}

func (options *AppOptions) bemfaGet(device *gree.Device) {
	newStatus := &Bemfa{}
	if device.Props[options.Commands.Power.Code] == 0 {
		newStatus.Power = false
	} else {
		newStatus.Power = true
	}
	newStatus.Mode = bemfa.Commands().Mode[getKeyByValue(options.Commands.Mode.Value, device.Props[options.Commands.Mode.Code])]
	if device.Props[options.Commands.PowerSave.Code] == 1 {
		newStatus.Mode = 7
	} else if device.Props[options.Commands.Sleep.Code] == 1 {
		newStatus.Mode = 6
	}
	newStatus.Temperature = device.Props[options.Commands.Temperature.Code]
	newStatus.FanSpeed = bemfa.Commands().FanSpeed[getKeyByValue(options.Commands.FanSpeed.Value, device.Props[options.Commands.FanSpeed.Code])]
	if device.Props[options.Commands.SwingHor.Code] == 1 {
		newStatus.SwingHor = 1
	} else {
		newStatus.SwingHor = 0
	}
	if device.Props[options.Commands.SwingVert.Code] == 1 {
		newStatus.SwingVert = 1
	} else {
		newStatus.SwingVert = 0
	}
	if !reflect.DeepEqual(*options.BemfaState, *newStatus) {
		options.BemfaState = newStatus
		power := "off"
		if newStatus.Power {
			power = "on"
		}
		value := fmt.Sprintf("%s#%d#%d#%d#%d#%d", power, newStatus.Mode, newStatus.Temperature, newStatus.FanSpeed, newStatus.SwingHor, newStatus.SwingVert)
		options.BemfaClient.Publish(options.BemfaTopic+"/set", 0, options.MqttRetain, value)
	}
}

func (options *AppOptions) bemfaSet(data string) {
	datas := strings.Split(data, "#")
	comms := make([]string, 0)
	values := make([]int, 0)
	if datas[0] == "off" {
		if options.BemfaState.Power {
			options.UdpClient.SetPower(false)
			options.BemfaState.Power = false
			return
		}
	} else {
		if !options.BemfaState.Power {
			comms = append(comms, options.Commands.Power.Code)
			values = append(values, 1)
			options.BemfaState.Power = true
		}
	}
	switch len(datas) {
	case 6:
		tmp, err := strconv.Atoi(datas[5])
		if err != nil {
			log.Errorf("[MQTT] Error SwingVert %s to int\n", datas[5])
			return
		}
		if options.BemfaState.SwingVert != tmp {
			comms = append(comms, options.Commands.SwingVert.Code)
			values = append(values, options.Commands.SwingVert.Value[getKeyByValue(bemfa.Commands().SwingVert, tmp)])
			options.BemfaState.SwingVert = tmp
		}
		fallthrough
	case 5:
		tmp, err := strconv.Atoi(datas[4])
		if err != nil {
			log.Errorf("[MQTT] Error SwingHor %s to int\n", datas[4])
			return
		}
		if options.BemfaState.SwingHor != tmp {
			comms = append(comms, options.Commands.SwingHor.Code)
			values = append(values, options.Commands.SwingVert.Value[getKeyByValue(bemfa.Commands().SwingHor, tmp)])
			options.BemfaState.SwingHor = tmp
		}
		fallthrough
	case 4:
		tmp, err := strconv.Atoi(datas[3])
		if err != nil {
			log.Errorf("[MQTT] Error FanSpeed %s to int\n", datas[3])
			return
		}
		if options.BemfaState.FanSpeed != tmp {
			comms = append(comms, options.Commands.FanSpeed.Code)
			values = append(values, options.Commands.FanSpeed.Value[getKeyByValue(bemfa.Commands().FanSpeed, tmp)])
			options.BemfaState.FanSpeed = tmp
		}
		fallthrough
	case 3:
		tmp, err := strconv.Atoi(datas[2])
		if err != nil {
			log.Errorf("[MQTT] Error Temperature %s to int\n", datas[2])
			return
		}
		if options.BemfaState.Temperature != tmp {
			comms = append(comms, options.Commands.Temperature.Code)
			values = append(values, tmp)
			comms = append(comms, options.Commands.TemperatureUnit.Code)
			values = append(values, 0)
			options.BemfaState.Temperature = tmp
		}
		fallthrough
	case 2:
		tmp, err := strconv.Atoi(datas[1])
		if err != nil {
			log.Errorf("[MQTT] Error Mode %s to int\n", datas[1])
			return
		}
		if options.BemfaState.Mode != tmp {
			switch tmp {
			case 7:
				comms = append(comms, options.Commands.PowerSave.Code)
				values = append(values, 1)
			case 6:
				comms = append(comms, options.Commands.Sleep.Code)
				values = append(values, 1)
			case 5, 4, 3, 2, 1:
				comms = append(comms, options.Commands.Mode.Code)
				values = append(values, options.Commands.Mode.Value[getKeyByValue(bemfa.Commands().Mode, tmp)])
			}
			options.BemfaState.Mode = tmp
		}
	}
	if len(comms) != 0 || len(values) != 0 {
		options.UdpClient.SetBFCommand(comms, values)
	}
}

func (options *AppOptions) Start(mqtt mqtt.Client, bemfa mqtt.Client, greeConfig *config.Gree) {

	deviceFactory := gree.DeviceFactory{
		Host:        options.HvacHost,
		OnStatus:    options.UpdateStatusToMqtt,
		OnConnected: options.OnConnected,
		OnUpdate:    options.OnUpdate,
	}

	options.Commands = gree.Commands()
	options.HvacHost = greeConfig.Host
	options.mqttClient = mqtt
	options.MqttTopicPrefix = greeConfig.HavcTopic
	options.BemfaTopic = greeConfig.BemfaTopic
	options.BemfaClient = bemfa
	options.UdpClient = gree.Create(&deviceFactory)
	options.DeviceState = make(map[string]string)
	options.BemfaState = &Bemfa{}
	go options.UdpClient.ConnectToDevice(options.HvacHost)

	// select {}
}

func Create() *AppOptions {
	return &AppOptions{}
}
