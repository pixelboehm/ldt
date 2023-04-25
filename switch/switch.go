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
		Name:             "Exciting Switch",
		SerialNumber:     "051AC-23AAM1",
		Manufacturer:     "Apple",
		Model:            "AB",
		FirmwareRevision: "1.0.1",
	}
	ac := accessory.NewSwitch(info)
	ac.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			log.Println("Switch is on")
		} else {
			log.Println("Switch is off")
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
