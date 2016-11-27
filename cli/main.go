package main

import (
	"fmt"
	"log"

	"github.com/vincentserpoul/cottonwood/lib/usb"
)

func main() {
	cottonwoodDevices, err := usb.GetCottonwoodDevices()
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, cottonwoodDevice := range cottonwoodDevices {
		add, err := cottonwoodDevice.GetDeviceAddress()
		if err != nil {
			log.Fatalf("%v", err)
		}
		port, err := cottonwoodDevice.GetPortNumber()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("found cottonwood device @address %d, port %d\n", add, port)
	}
}
