package config

import (
	"flag"
)

var serverAddress *string
var finalAddress *string

func ParseAddresses() {
	serverAddress = flag.String("a", ":8080", "server address")
	finalAddress = flag.String("b", "localhost", "destination folder")
	flag.Parse()
}

func GetFinAddr() string {
	return *finalAddress
}

func GetSrvAddr() string {
	return *serverAddress
}
