package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-ldts/pcl"
	"go-ldts/wotm"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

var device_IPv4 string
var ldt_IPv4 string
var ldt_identifier string
var port string
var ldt_specific_folder string

func main() {
	ldt_specific_folder = os.Args[1]
	port = os.Args[2]
	device_IPv4 = os.Args[3]
	ldt_identifier = os.Args[4]

	router := pcl.SetupRouter()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	pcl.AddHTTPHandler(router, "/register", registerDevice)
	ldt_IPv4, err := pcl.GetIPAddress()
	if err != nil {
		panic(err)
	}

	go pcl.Run(router, ldt_IPv4, port)
	go printDeviceAddress()

	lightbulb := accessory.NewLightbulb(accessory.Info{
		Name: ldt_identifier + " " + "Lamp",
	})
	fs := hap.NewFsStore(ldt_specific_folder + "/db")
	server, err := hap.NewServer(fs, lightbulb.A)
	if err != nil {
		log.Panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()

	lightbulb.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			turnOn(client)
			log.Println("Lightbulb is on")
		} else {
			turnOff(client)
			log.Println("Lightbulb is off")
		}
	})
	server.ListenAndServe(ctx)
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
	wotm.WriteAddressesToDescription(ldt_IPv4, ldt_identifier, payload.Device_IPv4, payload.Device_MAC, port, ldt_specific_folder)
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
