package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pcl "go-ldts/pcl"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

var device_IPv4 string
var ldt_IPv4 string
var config hc.Config
var ldt_name string

func main() {
	var ldt_specific_folder string = os.Args[1]
	var port string = os.Args[2]
	device_IPv4 = os.Args[3]
	ldt_name = os.Args[4]

	router := pcl.SetupRouter()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	pcl.AddHTTPHandler(router, "/register", registerDevice)

	// create an accessory
	info := accessory.Info{
		Name:             "Lightbulb " + ldt_name,
		SerialNumber:     "042DH-01BCN3",
		Manufacturer:     "Apple",
		Model:            "AB",
		FirmwareRevision: "1.0.1",
	}
	ac := accessory.NewLightbulb(info)
	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			turnOn(client)
		} else {
			turnOff(client)
		}
	})

	// configure the ip transport
	config = hc.Config{
		Pin:         "00000009",
		StoragePath: ldt_specific_folder,
	}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
		os.Exit(1)
	})

	ldt_IPv4, err = pcl.GetIPAddress()
	if err != nil {
		panic(err)
	}
	go pcl.Run(router, ldt_IPv4, port)
	go printDeviceAddress()

	t.Start()
}

func turnOn(client *http.Client) {
	req, err := http.NewRequest(http.MethodGet, "http://"+device_IPv4+"/on", nil)
	if err != nil {
		log.Println("Failed to create request")
		return
	}
	_, err = client.Do(req)
	if err != nil {
		log.Println(fmt.Sprint("Failed to do the request", err))
	}
	log.Println("Lightbulb: Turn On")
}

func turnOff(client *http.Client) {
	req, err := http.NewRequest(http.MethodGet, "http://"+device_IPv4+"/off", nil)
	if err != nil {
		log.Println("Failed to create request")
		return
	}
	_, err = client.Do(req)
	if err != nil {
		log.Println(fmt.Sprint("Failed to do the request", err))
	}
	log.Println("Lightbulb: Turn Off")
}

func registerDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("A new connection request from a Device arrived!")

	type RequestPayload struct {
		Device_IPv4 string
		Device_MAC  string
	}

	var payload RequestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Println("Lightbulb: Decoding Error: ", err)
	}
	defer r.Body.Close()

	device_IPv4 = payload.Device_IPv4
	pcl.WriteAddressesToDescription(ldt_IPv4, ldt_name, payload.Device_IPv4, payload.Device_MAC, config.StoragePath)
	w.Write([]byte("ack"))
}

func printDeviceAddress() {
	for {
		ticker := time.NewTicker(4 * time.Second)
		if device_IPv4 == "" {
			log.Printf("Device Address: <none>")
		} else {
			log.Printf("Device Address: %s", device_IPv4)
		}
		<-ticker.C
	}
}
