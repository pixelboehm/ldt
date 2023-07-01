package pcl

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type Data struct {
	Device_Name string
	Device_IPv4 string
	Device_MAC  string
	Ldt_IPv4    string
}

var port int

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()
	return router
}

func Run(router *http.ServeMux, ip, port_number string) {
	var err error
	port, err = strconv.Atoi(port_number)
	if err != nil {
		panic("unable to convert port, please use a number as the second parameter")
	}

	var addr string = fmt.Sprintf(":%d", port)
	log.Printf("HTTP serve at %s%s\n", ip, addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}

func AddHTTPHandler(router *http.ServeMux, route string, handler func(w http.ResponseWriter, r *http.Request)) {
	router.HandleFunc(route, handler)
}

func WriteAddressesToDescription(ldt_address, device_Name, device_IPv4, device_MAC, storatePath string) error {
	var description string = storatePath + "/wotm/description.json"

	Temp := Data{
		Device_Name: device_Name,
		Device_IPv4: device_IPv4,
		Device_MAC:  device_MAC,
		Ldt_IPv4:    ldt_address + ":" + strconv.Itoa(port),
	}

	t, err := template.ParseFiles(description)
	if err != nil {
		return errors.New(fmt.Sprint("PCL: Failed to parse template: ", err))
	}

	output, err := os.Create(description)
	if err != nil {
		return errors.New(fmt.Sprint("PCL: Failed to open description file: ", err))
	}
	defer output.Close()

	err = t.Execute(output, Temp)
	if err != nil {
		return errors.New(fmt.Sprint("PCL: Failed to write into template: ", err))
	}
	return nil
}

func GetIPAddress() (string, error) {
	hostname, _ := os.Hostname()

	ipAddr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		return "", errors.New(fmt.Sprint("PCL: Failed wo obtain Host-IP Address", err))
	}
	return ipAddr.IP.String(), nil
}
