package gotmclibusb

import (
	"fmt"
	"log"
	"strings"

	"github.com/gotmc/libusb"
)

const notavailable = "N/A"

// GetCottonwoodDevices retrieves the list of connected cottonwood devices
func GetCottonwoodDevices() (cottonwoodDevices []*libusb.Device, err error) {

	ctx, err := libusb.Init()
	if err != nil {
		return cottonwoodDevices, fmt.Errorf("Couldn't create USB context. %v", err)
	}
	defer func() {
		errEx := ctx.Exit()
		if errEx != nil {
			log.Fatal(errEx)
		}
	}()
	devices, err := ctx.GetDeviceList()
	if err != nil {
		return cottonwoodDevices, fmt.Errorf("Couldn't get devices. %v", err)
	}
	log.Printf("Found %v USB devices.\n", len(devices))
	for _, device := range devices {
		product, _, manufacturer, errDev := getDeviceDesc(device)
		if errDev == nil && product == "AS3991" && manufacturer == "AMS" {
			cottonwoodDevices = append(cottonwoodDevices, device)
		}

		err = errDev
	}

	if len(cottonwoodDevices) == 0 {
		return cottonwoodDevices, fmt.Errorf(
			"Couldn't find any Cottonwood device.\n %v.\n"+
				"Maybe you can try chown on /dev/bus/usb/<bus>/<device>", err)
	}

	return cottonwoodDevices, nil

}

func getDeviceDesc(device *libusb.Device) (
	product string,
	serialNumber string,
	manufacturer string,
	err error,
) {
	usbDeviceDescriptor, err := device.GetDeviceDescriptor()
	if err != nil {
		return product,
			serialNumber,
			manufacturer,
			fmt.Errorf("Error getting device descriptor: %s", err)
	}
	handle, err := device.Open()
	if err != nil {
		return product,
			serialNumber,
			manufacturer,
			fmt.Errorf("Error opening device: %s", err)
	}
	defer func() {
		err = handle.Close()
		if err != nil {
			log.Fatalf("Error closing handle: %s", err)
		}
	}()
	serialNumber, err = handle.GetStringDescriptorASCII(usbDeviceDescriptor.SerialNumberIndex)
	if err != nil {
		serialNumber = notavailable
	}
	manufacturer, err = handle.GetStringDescriptorASCII(usbDeviceDescriptor.ManufacturerIndex)
	if err != nil {
		manufacturer = notavailable
	}
	product, err = handle.GetStringDescriptorASCII(usbDeviceDescriptor.ProductIndex)
	if err != nil {
		product = notavailable
	}
	return strings.TrimSpace(product),
		strings.TrimSpace(serialNumber),
		strings.TrimSpace(manufacturer),
		nil
}
