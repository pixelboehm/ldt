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

func GetIPAddress() (string, error) {
	hostname, _ := os.Hostname()

	ipAddr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		return "", errors.New(fmt.Sprint("PCL: Failed wo obtain Host-IP Address", err))
	}
	return ipAddr.IP.String(), nil
}
