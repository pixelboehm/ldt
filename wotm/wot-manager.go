package wotm

import (
	"errors"
	"fmt"
	"html/template"
	"os"
)

type TemplateData struct {
	Device_Name string
	Device_IPv4 string
	Device_MAC  string
	Ldt_IPv4    string
}

func WriteAddressesToDescription(ldt_address, device_Name, device_IPv4, device_MAC, port, storatePath string) error {
	var description string = storatePath + "/wotm/description.json"

	Temp := TemplateData{
		Device_Name: device_Name,
		Device_IPv4: device_IPv4,
		Device_MAC:  device_MAC,
		Ldt_IPv4:    ldt_address + ":" + port,
	}

	t, err := template.ParseFiles(description)
	if err != nil {
		return errors.New(fmt.Sprint("<WoT-Manager>: Failed to parse template: ", err))
	}

	output, err := os.Create(description)
	if err != nil {
		return errors.New(fmt.Sprint("<WoT-Manager>: Failed to open description file: ", err))
	}
	defer output.Close()

	err = t.Execute(output, Temp)
	if err != nil {
		return errors.New(fmt.Sprint("<WoT-Manager>: Failed to write into template: ", err))
	}
	return nil
}
