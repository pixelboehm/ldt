package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	pcl "go-ldts/pcl"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

var device_address string
var config hc.Config
var ip_address string

func main() {
	var ldt_specific_folder string = os.Args[1]
	var port string = os.Args[2]
	device_address = os.Args[3]

	router := pcl.SetupRouter()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	pcl.AddHTTPHandler(router, "/register", registerDevice)

	// create an accessory
	info := accessory.Info{
		Name:             "Awesome Lightbulb",
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
		StoragePath: "/usr/local/etc/orchestration-manager/" + ldt_specific_folder,
	}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
		os.Exit(1)
	})

	ip_address, err = pcl.GetIPAddress()
	if err != nil {
		panic(err)
	}
	go pcl.Run(router, ip_address, port)
	go printDeviceAddress()

	t.Start()
}

func turnOn(client *http.Client) {
	req, err := http.NewRequest(http.MethodGet, "http://"+device_address+"/on", nil)
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
	req, err := http.NewRequest(http.MethodGet, "http://"+device_address+"/off", nil)
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
	log.Println("A new Device wants to register")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body")
		return
	}
	device_address = string(body)
	pcl.AddIPToDescription(ip_address, device_address, config.StoragePath)
	w.Write([]byte("ack"))
}

func printDeviceAddress() {
	for {
		ticker := time.NewTicker(4 * time.Second)
		if device_address == "" {
			log.Printf("Device Address: <none>")
		} else {
			log.Printf("Device Address: %s", device_address)
		}
		<-ticker.C
	}
}
