package app

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
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
}

type UDPInfo struct {
	T    string `json:"t"`
	I    int    `json:"i"`
	Uid  int    `json:"uid"`
	Cid  string `json:"cid"`
	Tcid string `json:"tcid"`
	Pack string `json:"pack"`
}

var device = Device{}

func Create(option DeviceFactory) *DeviceFactory {
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
	options.Conn = conn
	if err != nil {
		log.Errorf("[UDP} Error: %s Waite 10s!", err)
		time.Sleep(time.Second * 10)
		options.ConnectToDevice(address)
	}
	message := []byte("{\"t\":\"scan\"}")

	go options.handleResponse(conn)
	//log.Infof("[UDP] Send connect request: %s", message)
	_, err = conn.Write(message)
	if err != nil {
		log.Errorf("[UDP] Error: %s", err)
	}
}

func (options *DeviceFactory) sendBindRequest() {
	message := struct {
		Mac string `json:"mac"`
		T   string `json:"t"`
		Uid int    `json:"uid"`
	}{
		Mac: device.Id,
		T:   "bind",
		Uid: 0,
	}
	//log.Infof("[UDP] Send bind request: %v", message)
	encryptedBoundMessage := Encrypt(message, "")

	request := UDPInfo{
		Cid:  "app",
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
	//log.Infof("[UDP] Send bind request: %s", requestJson)
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
		Mac: device.Id,
		T:   "status",
	}
	options.sendRequest(message)
	time.Sleep(time.Second * 3)
	options.requestDeviceStatus()
}

func (options *DeviceFactory) sendRequest(message interface{}) {
	encryptedMessage := Encrypt(message, device.Key)
	request := UDPInfo{
		Tcid: "",
		Cid:  "app",
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
	//props := make(Props)
	for {
		data := make([]byte, 1024)
		read, err := conn.Read(data)
		if err != nil {
			log.Errorf("[UDP] Error: %s", err)
		}
		//log.Infof("[UDP] Received: %s from %s", data[:read], conn.RemoteAddr())

		udpInfo := UDPInfo{}
		err = json.Unmarshal(data[:read], &udpInfo)
		if err != nil {
			log.Errorf("[JSON] Error: %s", err)
		}

		pack := Decrypt(udpInfo, device.Key)

		if pack.T == "dev" {
			log.Infof("[UDP] Nw Device Registered: %s", pack.Name)
			device.Id = "app"
			//device.Id = pack.Mac
			device.Name = pack.Name
			device.Address = conn.RemoteAddr().String()
			device.Port = 7000
			device.Bound = false
			device.Props = nil
			options.sendBindRequest()
			continue
		}
		if pack.T == "bindok" {
			log.Infof("[UDP] Bound to %s", device.Name)
			device.Bound = true
			device.Key = pack.Key
			go options.requestDeviceStatus()
			options.OnConnected()
			continue
		}
		if pack.T == "dat" && device.Bound {
			log.Infof("[UDP] Received Data from %s", device.Name)
			device.Props = make(map[string]int)
			for i := 0; i < len(pack.Dat); i++ {
				device.Props[pack.Cols[i]] = pack.Dat[i]
				//props[pack.Cols[i]] = pack.Dat[i]
			}
			options.OnStatus(&device)
			continue
		}
		if pack.T == "res" && device.Bound {
			log.Infof("[UDP] Received Data from %s", device.Name)
			device.Props = make(map[string]int)
			for i := 0; i < len(pack.Dat); i++ {
				device.Props[pack.Cols[i]] = pack.Val[i]
			}
			options.OnUpdate(&device)
			continue
		}
	}

}

// SetPower /**
// * @param {string} power - "on" or "off"
// */
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
