package gree

import (
	"encoding/json"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Device struct {
	Id      string
	Name    string
	Address string
	Port    int
	Bound   bool
	Props   map[string]int
	Key     string
}

type DeviceFactory struct {
	Host        net.IP
	OnStatus    func(device *Device)
	OnUpdate    func(device *Device)
	OnConnected func()
	Conn        net.Conn
	Device      Device
}

type UDPInfo struct {
	T    string `json:"t"`
	I    int    `json:"i"`
	Uid  int    `json:"uid"`
	Cid  string `json:"cid"`
	Tcid string `json:"tcid"`
	Pack string `json:"pack"`
}

func Create(option *DeviceFactory) *DeviceFactory {
	return &DeviceFactory{
		Host:        option.Host,
		OnStatus:    option.OnStatus,
		OnUpdate:    option.OnUpdate,
		OnConnected: option.OnConnected,
		Conn:        nil,
	}
}

func (options *DeviceFactory) ConnectToDevice(address net.IP) {
	scrAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: address, Port: 7000}

	conn, err := net.DialUDP("udp", scrAddr, dstAddr)
	if err != nil {
		log.Errorf("[UDP} Error: %s Waite 10s!", err)
		time.Sleep(time.Second * 10)
		options.ConnectToDevice(address)
	}

	options.Conn = conn
	message := []byte("{\"t\":\"scan\"}")

	go options.handleResponse(conn)

	if _, err = conn.Write(message); err != nil {
		log.Errorf("[UDP] Error: %s", err)
	}
}

func (options *DeviceFactory) sendBindRequest() {
	message := struct {
		Mac string `json:"mac"`
		T   string `json:"t"`
		Uid int    `json:"uid"`
	}{
		Mac: options.Device.Id,
		T:   "bind",
		Uid: 0,
	}
	encryptedBoundMessage := Encrypt(message, "")

	request := UDPInfo{
		Cid:  "gree",
		I:    1,
		T:    "pack",
		Uid:  0,
		Pack: encryptedBoundMessage,
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		log.Errorf("[JSON] Error: %s", err)
	}
	_, err = options.Conn.Write(requestJson)
	if err != nil {
		log.Errorf("[UDP] Seed Error: %s", err)
	}
}

func (options *DeviceFactory) requestDeviceStatus() {
	message := struct {
		Cols []string `json:"cols"`
		Mac  string   `json:"mac"`
		T    string   `json:"t"`
	}{
		Cols: []string{
			"Pow",
			"Mod",
			"TemUn",
			"SetTem",
			"WdSpd",
			"Air",
			"Blo",
			"Health",
			"SwhSlp",
			"Lig",
			"SwingLfRig",
			"SwUpDn",
			"Quiet",
			"Tur",
			"SvSt",
		},
		Mac: options.Device.Id,
		T:   "status",
	}
	options.sendRequest(message)
	time.Sleep(time.Second * 3)
	options.requestDeviceStatus()
}

func (options *DeviceFactory) sendRequest(message interface{}) {
	encryptedMessage := Encrypt(message, options.Device.Key)
	request := UDPInfo{
		Tcid: "",
		Cid:  "gree",
		I:    0,
		T:    "pack",
		Uid:  0,
		Pack: encryptedMessage,
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		log.Errorf("[UDP] Error: %s", err)
	}
	_, err = options.Conn.Write(requestJson)
	if err != nil {
		log.Errorf("[UDP] Error: %s", err)
	}
}

// SendCommand /**
func (options *DeviceFactory) sendCommand(commends []string, values []int) {
	message := struct {
		Opt []string `json:"opt"`
		P   []int    `json:"p"`
		T   string   `json:"t"`
	}{
		Opt: commends,
		P:   values,
		T:   "cmd",
	}
	options.sendRequest(message)
}

func (options *DeviceFactory) handleResponse(conn net.Conn) {
	for {
		data := make([]byte, 1024)
		read, err := conn.Read(data)
		if err != nil {
			log.Errorf("[UDP] Error: %s", err)
		}

		udpInfo := UDPInfo{}
		err = json.Unmarshal(data[:read], &udpInfo)
		if err != nil {
			log.Errorf("[JSON] Error: %s", err)
		}

		pack := Decrypt(udpInfo, options.Device.Key)

		if pack.T == "dev" {
			options.Device.Id = pack.Mac
			options.Device.Name = pack.Name
			options.Device.Address = conn.RemoteAddr().String()
			options.Device.Port = 7000
			options.Device.Bound = false
			options.Device.Props = nil
			options.sendBindRequest()
			continue
		}
		if pack.T == "bindok" {
			log.Infof("[UDP] Bound to %s", options.Device.Name)
			options.Device.Bound = true
			options.Device.Key = pack.Key
			go options.requestDeviceStatus()
			options.OnConnected()
			continue
		}
		if pack.T == "dat" && options.Device.Bound {
			options.Device.Props = make(map[string]int)
			for i := 0; i < len(pack.Dat); i++ {
				options.Device.Props[pack.Cols[i]] = pack.Dat[i]
			}
			options.OnStatus(&options.Device)
			continue
		}
		if pack.T == "res" && options.Device.Bound {
			options.Device.Props = make(map[string]int)
			for i := 0; i < len(pack.Dat); i++ {
				options.Device.Props[pack.Cols[i]] = pack.Val[i]
			}
			options.OnUpdate(&options.Device)
			continue
		}
		log.Errorf("[UDP] Unknown Message of type %s: %v, %v", pack.T, data[:read], pack)
	}
}

func (options *DeviceFactory) SetPower(value bool) {
	if value {
		options.sendCommand([]string{Commands().Power.Code}, []int{1})
	} else {
		options.sendCommand([]string{Commands().Power.Code}, []int{0})
	}
}

// SetTemperature /**
func (options *DeviceFactory) SetTemperature(value int, unit int) {
	options.sendCommand(
		[]string{
			Commands().TemperatureUnit.Code,
			Commands().Temperature.Code,
		},
		[]int{
			unit,
			value,
		})
}

func (options *DeviceFactory) SetMode(value int) {
	options.sendCommand([]string{Commands().Mode.Code}, []int{value})
}

func (options *DeviceFactory) SetFanSpeed(value int) {
	options.sendCommand([]string{Commands().FanSpeed.Code}, []int{value})
}

func (options *DeviceFactory) SetSwingHor(value int) {
	options.sendCommand([]string{Commands().SwingHor.Code}, []int{value})
}

func (options *DeviceFactory) SetSwingVert(value int) {
	options.sendCommand([]string{Commands().SwingVert.Code}, []int{value})
}

func (options *DeviceFactory) SetQuietMode(value int) {
	options.sendCommand([]string{Commands().Quiet.Code}, []int{value})
}

func (options *DeviceFactory) SetSleepMode(value bool) {
	if value {
		options.sendCommand([]string{Commands().Sleep.Code}, []int{1})
	} else {
		options.sendCommand([]string{Commands().Sleep.Code}, []int{0})
	}
}

func (options *DeviceFactory) SetLights(value bool) {
	if value {
		options.sendCommand([]string{Commands().Lights.Code}, []int{1})
	} else {
		options.sendCommand([]string{Commands().Lights.Code}, []int{0})
	}
}

func (options *DeviceFactory) SetTurbo(value bool) {
	if value {
		options.sendCommand([]string{Commands().Turbo.Code}, []int{1})
	} else {
		options.sendCommand([]string{Commands().Turbo.Code}, []int{0})
	}
}

func (options *DeviceFactory) SetBlow(value bool) {
	if value {
		options.sendCommand([]string{Commands().Blow.Code}, []int{1})
	} else {
		options.sendCommand([]string{Commands().Blow.Code}, []int{0})
	}
}

func (options *DeviceFactory) SetHealthMode(value bool) {
	if value {
		options.sendCommand([]string{Commands().Health.Code}, []int{1})
	} else {
		options.sendCommand([]string{Commands().Health.Code}, []int{0})
	}
}

func (options *DeviceFactory) SetAir(value bool) {
	if value {
		options.sendCommand([]string{Commands().Air.Code}, []int{1})
	} else {
		options.sendCommand([]string{Commands().Air.Code}, []int{0})
	}
}

func (options *DeviceFactory) SetPowerSave(value bool) {
	if value {
		options.sendCommand([]string{Commands().PowerSave.Code}, []int{1})
	} else {
		options.sendCommand([]string{Commands().PowerSave.Code}, []int{0})
	}
}

func (options *DeviceFactory) SetBFCommand(comm []string, value []int) {
	options.sendCommand(comm, value)
}
