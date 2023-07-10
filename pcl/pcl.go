package pcl

import (
	"errors"
	"fmt"
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

func Run(router *http.ServeMux, ip, port_number string) {
	var err error
	port, err := strconv.Atoi(port_number)
	if err != nil {
		panic("<PCL>: Unable to convert port, please use a number as the second parameter")
	}

	var addr string = fmt.Sprintf(":%d", port)
	log.Printf("<PCL>: HTTP serve at %s%s\n", ip, addr)
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
		return "", errors.New(fmt.Sprint("<PCL>: Failed wo obtain Host-IP Address", err))
	}
	return ipAddr.IP.String(), nil
}
