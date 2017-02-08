package main

import (
	"fmt"
	"log"

	"github.com/vincentserpoul/cottonwood/lib/gotmclibusb"
)

func main() {
	cottonwoodDevices, err := gotmclibusb.GetCottonwoodDevices()
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

		handle, err := cottonwoodDevice.Open()
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer func() {
			err = handle.Close()
			if err != nil {
				log.Fatalf("Error closing handle: %s", err)
			}
		}()

		// configDescriptor, err := cottonwoodDevice.GetActiveConfigDescriptor()
		// if err != nil {
		// 	log.Fatalf("Failed getting the active config: %v", err)
		// }
		// fmt.Printf("=> Max Power = %d mA\n",
		// 	configDescriptor.MaxPowerMilliAmperes)
		// var singularPlural string
		// if configDescriptor.NumInterfaces == 1 {
		// 	singularPlural = "interface"
		// } else {
		// 	singularPlural = "interfaces"
		// }
		// fmt.Printf("=> Found %d %s\n",
		// 	configDescriptor.NumInterfaces, singularPlural)
		// fmt.Printf("=> The first interface has %d alternate settings.\n",
		// 	configDescriptor.SupportedInterfaces[0].NumAltSettings)
		// firstDescriptor := configDescriptor.SupportedInterfaces[0].InterfaceDescriptors[0]
		// fmt.Printf("=> The first interface descriptor has a length of %d.\n", firstDescriptor.Length)
		// fmt.Printf("=> The first interface descriptor is interface number %d.\n", firstDescriptor.InterfaceNumber)
		// fmt.Printf("=> The first interface descriptor has %d endpoint(s).\n", firstDescriptor.NumEndpoints)
		// fmt.Printf(
		// 	"   => USB-IF class %d, subclass %d, protocol %d.\n",
		// 	firstDescriptor.InterfaceClass, firstDescriptor.InterfaceSubClass, firstDescriptor.InterfaceProtocol,
		// )
		// for i, endpoint := range firstDescriptor.EndpointDescriptors {
		// 	fmt.Printf(
		// 		"   => Endpoint index %d on Interface %d has the following properties:\n",
		// 		i, firstDescriptor.InterfaceNumber)
		// 	fmt.Printf("     => Address: %d (b%08b)\n", endpoint.EndpointAddress, endpoint.EndpointAddress)
		// 	fmt.Printf("       => Endpoint #: %d\n", endpoint.Number())
		// 	fmt.Printf("       => Direction: %s (%d)\n", endpoint.Direction(), endpoint.Direction())
		// 	fmt.Printf("     => Attributes: %d (b%08b) \n", endpoint.Attributes, endpoint.Attributes)
		// 	fmt.Printf("       => Transfer Type: %s (%d) \n", endpoint.TransferType(), endpoint.TransferType())
		// 	fmt.Printf("     => Max packet size: %d\n", endpoint.MaxPacketSize)
		// }

		err = handle.ClaimInterface(0)
		if err != nil {
			log.Printf("Error claiming interface %s", err)
		}
		// // Send USBTMC message to Agilent 33220A
		// bulkOutput := firstDescriptor.EndpointDescriptors[0]
		// address := bulkOutput.EndpointAddress
		// fmt.Printf("Set frequency/amplitude on endpoint address %d\n", address)
		// data := createGotmcMessage("apply:sinusoid 2340, 0.1, 0.0")
		// transferred, err := handle.BulkTransfer(address, data, len(data), 5000)
		// if err != nil {
		// 	log.Printf("Error on bulk transfer %s", err)
		// }
		// fmt.Printf("Sent %d bytes to 33220A\n", transferred)
		err = handle.ReleaseInterface(0)
		if err != nil {
			log.Printf("Error releasing interface %s", err)
		}
	}
}
