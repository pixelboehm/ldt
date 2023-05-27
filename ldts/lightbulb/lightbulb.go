package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	pcl "go-ldts/pcl"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

var Device_address string

func main() {
	router := pcl.SetupRouter()
	pcl.AddHTTPHandler(router, "/register", registerDevice)

	go pcl.Run(router, os.Args[2])
	go printDeviceAddress()

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
			turnOn()
		} else {
			turnOff()
		}
	})

	// configure the ip transport
	config := hc.Config{
		Pin:         "00000009",
		StoragePath: "/usr/local/etc/orchestration-manager/" + os.Args[1],
	}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
		os.Exit(1)
	})

	t.Start()
}

func turnOn() {
	req, err := http.NewRequest(http.MethodGet, Device_address+"/on", nil)
	if err != nil {
		log.Println("Failed to create request")
		return
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Failed to do the request")
	}
	log.Println("Send turn On Command")
}

func turnOff() {
	req, err := http.NewRequest(http.MethodGet, Device_address+"/off", nil)
	if err != nil {
		log.Println("Failed to create request")
		return
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Failed to do the request")
	}
	log.Println("Send turn off command")
}

func registerDevice(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body")
		return
	}
	Device_address = string(body)
	w.Write([]byte("ack"))
}

func printDeviceAddress() {
	for {
		ticker := time.NewTicker(2 * time.Second)
		log.Printf("Device Address: %s", Device_address)
		<-ticker.C
	}
}
