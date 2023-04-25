package main

import (
	"log"
	"os"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

func main() {
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
			log.Println("Lightbulb is on")
		} else {
			log.Println("Lightbulb is off")
		}
	})

	// configure the ip transport
	config := hc.Config{
		Pin:         "00000009",
		StoragePath: "/usr/local/etc/orchestration-manager/" + info.Name + "_" + os.Args[1],
	}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
