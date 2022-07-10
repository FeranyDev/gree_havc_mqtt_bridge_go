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
	Props   map[string]string
	Key     string
}

type DeviceFactory struct {
	Host        net.IP
	OnStatus    func(device *Device)
	OnUpdate    func(device *Device)
	OnConnected func()
}

var device = Device{}

var conn = net.UDPConn{}

func (options *DeviceFactory) create(option DeviceFactory) {
	options.Host = option.Host
	options.OnStatus = option.OnStatus
	options.OnUpdate = option.OnUpdate
	options.OnConnected = option.OnConnected

	options.connectToDevice(options.Host)

}

func (options *DeviceFactory) connectToDevice(address net.IP) {
	scrAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: address, Port: 5000}
	conn, err := net.DialUDP("udp", scrAddr, dstAddr)
	if err != nil {
		log.Errorf("[UDP} Error: %s", err)
		time.Sleep(time.Second * 10)
	}
	//defer conn.Close()
	message, err := json.Marshal("{ \"t\": \"scan\" }")
	if err != nil {
		log.Errorf("[UDP] Error: %s", err)
	}
	conn.Write(message)

	go options.handleResponse(conn)
	//options.OnConnected()
}

func (options *DeviceFactory) sendBindRequest() {
	message := struct {
		mac string
		t   string
		uid int
	}{
		mac: device.Id,
		t:   "bind",
		uid: 0,
	}
	encryptedBoundMessage := encrypt(message, "")

	request := struct {
		cid  string
		i    int
		t    string
		uid  int
		pack string
	}{
		cid:  "app",
		i:    1,
		t:    "pack",
		uid:  0,
		pack: encryptedBoundMessage,
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		log.Errorf("[UDP] Error: %s", err)
	}
	_, err = conn.Write(requestJson)
	if err != nil {
		log.Errorf("[UDP] Error: %s", err)
	}
}

func (options *DeviceFactory) confirmBinding() {

}

func (options *DeviceFactory) handleResponse(conn net.Conn) {

	data := make([]byte, 1024)

	for {
		read, err := conn.Read(data)
		if err != nil {
			log.Errorf("[UDP] Error: %s", err)
		}
		log.Infof("[UDP] Received: %s from %s", data[:read], conn.RemoteAddr())

		pack := decrypt(string(data), device.Key)

		if pack.T == "dev" {
			log.Infof("[UDP] Nw Device Registered: %s", pack.Name)
			device.Id = pack.Mac
			device.Name = pack.Name
			device.Address = conn.RemoteAddr().String()
			device.Port = 7000
			device.Bound = false
			device.Props = nil
			options.OnStatus(&device)
			options.sendBindRequest()
			return
		}
		if pack.T == "bindok" {
			log.Infof("[UDP] Bound to %s", device.Name)
			device.Bound = true
			device.Key = pack.Key
			return
		}
	}
}
