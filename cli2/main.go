// Copyright 2013 Google Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// rawread attempts to read from the specified USB device.
package main

import (
	"log"
	"strconv"
	"time"

	"github.com/thorduri/go-libusb/usb"
)

func main() {
	// Only one context should be needed for an application.  It should always be closed.
	ctx := usb.NewContext()
	defer func() {
		errCl := ctx.Close()
		if errCl != nil {
			log.Fatal(errCl)
		}
	}()

	ctx.Debug(1)

	// ListDevices is used to find the devices to open.
	devs, err := ctx.ListDevices(
		func(desc *usb.Descriptor) bool {
			if desc.Vendor == GetCottonwoodVendor() && desc.Product == GetCottonwoodProduct() {
				return true
			}
			return false
		})

	// All Devices returned from ListDevices must be closed.
	defer func() {
		for _, dev := range devs {
			errCl := dev.Close()
			if errCl != nil {
				log.Fatal(errCl)
			}
		}
	}()

	// ListDevices can occasionally  fail, so be sure to check its return value.
	if err != nil {
		log.Fatalf("list: %s", err)
	}

	for _, dev := range devs {
		// Once the device has been selected from ListDevices, it is opened
		// and can be interacted with.
		// Open up two ep for read and write

		epBulkWrite, err := dev.OpenEndpoint(1, 0, 0, 2|uint8(usb.ENDPOINT_DIR_OUT))
		if err != nil {
			log.Fatalf("OpenEndpoint Write error for %v: %v", dev.Address, err)
		}

		// Poll Firmware/Hardware Version ID

		// AntennaOn
		outAntennaPowerOnCmd := []byte{0x18, 0x03, 0xFF}
		// outFirmIdCmd := []byte{0x10, 0x03, 0x00}
		// outHardIdCmd := []byte{0x10, 0x03, 0x01}

		i, err := epBulkWrite.Write(outAntennaPowerOnCmd)
		if err != nil {
			log.Fatalf("Cannot write command: %v\n", err)
		}
		log.Printf("%v bytes sent", i)

		time.Sleep(1 * time.Second)

		epBulkRead, err := dev.OpenEndpoint(1, 0, 0, 1|uint8(usb.ENDPOINT_DIR_IN))
		if err != nil {
			log.Fatalf("OpenEndpoint Read error for %v: %v", dev.Address, err)
		}

		// readBuffer := make([]byte, 64)
		// n, errR := epBulkRead.Read(readBuffer)
		// if errR != nil {
		// 	fmt.Errorf("Cannot read: %v\n", errR)
		// }
		// log.Printf("%v bytes read\n", n)
		// fmt.Printf("%v bytes read\n", readBuffer)

		for {
			readBuffer := make([]byte, 64)
			n, errRead := epBulkRead.Read(readBuffer)
			log.Printf("read %d bytes: %v", n, readBuffer)
			if errRead != nil {
				log.Printf("error reading: %v\n", errRead)
			}
			time.Sleep(1 * time.Second)
		}

	}
}

const (
	OUT_FIRM_HARDW_ID = 0x10 // Firmware/Hardware version poll command
	// IN_FIRM_HARDW_ID = 0x11   version poll command
)

// GetCottonwoodVendor returns the vendor ID of cottonwood UHF reader
func GetCottonwoodVendor() usb.ID {
	value, err := strconv.ParseUint("1325", 16, 16)
	if err != nil {
		log.Fatal(err)
	}
	return usb.ID(value)
}

// GetCottonwoodProduct returns the product ID of cottonwood UHF reader
func GetCottonwoodProduct() usb.ID {
	value, err := strconv.ParseUint("c029", 16, 16)
	if err != nil {
		log.Fatal(err)
	}
	return usb.ID(value)
}

//   Endpoint Descriptor:
//     bLength                 7
//     bDescriptorType         5
//     bEndpointAddress     0x81  EP 1 IN
//     bmAttributes            3
//       Transfer Type            Interrupt
//       Synch Type               None
//       Usage Type               Data
//     wMaxPacketSize     0x0040  1x 64 bytes
//     bInterval              10
//   Endpoint Descriptor:
//     bLength                 7
//     bDescriptorType         5
//     bEndpointAddress     0x02  EP 2 OUT
//     bmAttributes            3
//       Transfer Type            Interrupt
//       Synch Type               None
//       Usage Type               Data
//     wMaxPacketSize     0x0040  1x 64 bytes
//     bInterval              10
