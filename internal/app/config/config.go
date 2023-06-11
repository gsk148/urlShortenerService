package config

import (
	"flag"
	"strconv"
	"strings"
)

var ServerAddress string
var FinalAddress string

type NetAddress struct {
	Host string
	Port int
	set  bool
}

func (a *NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		if hp[0] == "http" {
			a.Host = hp[1][2:]
			port, err := strconv.Atoi(hp[2])
			if err != nil {
				return err
			}
			a.Port = port
			a.set = true
			return nil
		}
	}

	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	a.set = true
	return nil
}

func (a *NetAddress) setDefault() {
	a.Host = "localhost"
	a.Port = 8080
}

func (a *NetAddress) checkFlagProvided() {
	if !a.set {
		a.setDefault()
	}
}

func ParseAddresses() {
	serverAddr := new(NetAddress)
	finalAddr := new(NetAddress)
	flag.Var(serverAddr, "a", "Server run: net address host:port")
	flag.Var(finalAddr, "b", "Returned address: net address host:port")
	flag.Parse()
	serverAddr.checkFlagProvided()
	finalAddr.checkFlagProvided()
	ServerAddress = serverAddr.String()
	FinalAddress = finalAddr.String()
}

func GetFinAddr() string {
	return FinalAddress
}

func GetSrvAddr() string {
	return ServerAddress
}