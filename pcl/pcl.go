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
	Device string
	Ldt    string
}

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()
	return router
}

func Run(router *http.ServeMux, ip, port_number string) {
	port, err := strconv.Atoi(port_number)
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

func AddIPToDescription(ldt_address, device_address, storatePath string) error {
	var description string = storatePath + "/wotm/description.json"
	log.Println("writing description: ", description)

	Temp := Data{
		Device: device_address,
		Ldt:    ldt_address,
	}

	t, err := template.ParseFiles(storatePath + "/wotm/description.json")
	if err != nil {
		return errors.New(fmt.Sprint("PCL: Failed to parse template: ", err))
	}

	output, err := os.OpenFile(description, os.O_RDWR, 0666)
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
