package main

import (
	"fmt"
	"os"

	"github.com/boombuler/hid"
)

func main() {
	cot := NewCottonRFID()
	cot.CMDFirm()

	stop := make(chan struct{})
	reports, errors := cot.Listen(stop)
	go func() {
		for {
			select {
			case rep := <-reports:
				fmt.Println("Report:", rep)
			case err := <-errors:
				fmt.Println("Error:", err)
			case <-stop:
				return
			}
		}
	}()
	fmt.Scanln()
	close(stop)

	fmt.Printf("%v", cot)
}

const (
	cottonVendorID  = 0x1325
	cottonProductID = 0xc029
)

const (
	cmdOutFirmHardID = 0x10
)

type CottonRFID struct {
	devices []hid.Device
	bytes   []byte
}

func NewCottonRFID() (cotton CottonRFID) {
	cotton.devices = findDevices()
	cotton.bytes = []byte{}
	return
}

func findDevices() []hid.Device {
	devices := []hid.Device{}
	deviceInfos := hid.Devices()
	for {
		info, more := <-deviceInfos
		if more {
			device, error := info.Open()
			if error != nil {
				fmt.Println(error)
			}
			if !isCottonRFIDDevice(*info) {
				fmt.Printf("%s %s is not a CottonRFID device.\n", info.Manufacturer, info.Product)
			} else {
				devices = append(devices, device)
				fmt.Printf("%s %s is a CottonRFID device.\n", info.Manufacturer, info.Product)
			}
		} else {
			break
		}
	}
	if len(devices) == 0 {
		fmt.Println("No CottonRFIDs found.")
		os.Exit(1)
	}
	return devices
}

func isCottonRFIDDevice(deviceInfo hid.DeviceInfo) bool {
	// from forums: "Blync creates 2 HID devices and the only way to find out the right device is the MaxFeatureReportLength = 0"
	if deviceInfo.VendorId == cottonVendorID && deviceInfo.ProductId == cottonProductID && deviceInfo.FeatureReportLength == 0 {
		return true
	}
	return false
}

func (b CottonRFID) CMDFirm() {
	for _, device := range b.devices {
		error := device.Write([]byte{cmdOutFirmHardID, 0x03, 0x00})
		if error != nil {
			fmt.Println(error)
		}
	}
}
