package pcl

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()
	return router
}

func Run(router *http.ServeMux, number, storatePath string) {

	port, err := strconv.Atoi(number)
	if err != nil {
		panic("unable to convert port, please use a number as the second parameter")
	}

	ip, err := getIPAddress()
	if err != nil {
		panic(err)
	}

	if err := addIPToDescription(ip, storatePath); err != nil {
		panic(err)
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

func addIPToDescription(ipAddress, storatePath string) error {
	var description string = storatePath + "/wotm/description.json"

	t, err := template.ParseFiles(storatePath + "/wotm/description.json")
	if err != nil {
		return errors.New(fmt.Sprint("PCL: Failed to parse template: ", err))
	}

	output, err := os.OpenFile(description, os.O_RDWR, 0666)
	if err != nil {
		return errors.New(fmt.Sprint("PCL: Failed to open description file: ", err))
	}
	defer output.Close()

	err = t.Execute(output, ipAddress)
	if err != nil {
		return errors.New(fmt.Sprint("PCL: Failed to write into template: ", err))
	}
	return nil
}

func getIPAddress() (string, error) {
	hostname, _ := os.Hostname()

	ipAddr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		return "", errors.New(fmt.Sprint("PCL: Failed wo obtain Host-IP Address", err))
	}
	return ipAddr.IP.String(), nil
}
