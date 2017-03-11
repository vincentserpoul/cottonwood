package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"unsafe"
)

func main() {
	devices, err := getDevices()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(devices))
	for _, d := range devices {
		fmt.Printf("device %v\n", d)
		dev, err := Open(d)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			dev.device.Close()
		}()
		fmt.Println("opened", d.Path)

		resp, err := dev.Command(0x10, []byte{0x00})
		if err != nil {
			log.Fatal(err)
		}

		for i, charact := range resp {
			if i > 1 && charact != 0 {
				fmt.Printf("%c", charact)
			}
		}
		fmt.Print("\n")

		resp, err = dev.Command(0x10, []byte{0x01})
		if err != nil {
			log.Fatal(err)
		}

		for i, charact := range resp {
			if i > 1 && charact != 0 {
				fmt.Printf("%c", charact)
			}
		}
		fmt.Print("\n")

		resp, err = dev.Command(0x18, []byte{0xff})
		if err != nil {
			log.Fatal(err)
		}

		for _, charact := range resp {
			fmt.Printf("0x%X ", charact)
		}
		fmt.Print("\n")

		for {
			respInv, errInv := dev.Command(0x43, []byte{0x01})
			if errInv != nil {
				log.Fatal(errInv)
			}

			for _, charact := range respInv {
				fmt.Printf("0x%X ", charact)
			}
			fmt.Print("\n")

			respInv2, errInv2 := dev.Command(0x43, []byte{0x02})
			if errInv2 != nil {
				log.Fatal(errInv2)
			}

			for _, charact := range respInv2 {
				fmt.Printf("0x%X ", charact)
			}
			fmt.Print("\n")

			time.Sleep(5 * time.Second)
		}

	}
}

func (dev *Device) Command(reportID byte, payload []byte) ([]byte, error) {

	/* Set Feature */
	bufS := make([]byte, 64)
	bufS[0] = reportID
	bufS[1] = byte(len(payload) + 2)
	for i, payloadItm := range payload {
		bufS[i+2] = payloadItm
	}

	fd := dev.device.f.Fd()
	err := Ioctl(fd, IoctlHIDIOCSFEATURE(64), uintptr(unsafe.Pointer(&bufS[0])))
	if err != nil {
		return nil, err
	}

	// wait a little
	time.Sleep(10 * time.Millisecond)

	bufG := make([]byte, 64)
	bufG[0] = reportID
	err = Ioctl(fd, IoctlHIDIOCGFEATURE(64), uintptr(unsafe.Pointer(&bufG[0])))
	if err != nil {
		return nil, err
	}

	return bufG, nil
}

const (
	cottonVendorID  = 0x1325
	cottonProductID = 0xc029

	maxMessageLen = 64

	responseTimeout = 3 * time.Second
)

// var errorCodes = map[uint8]string{
// 	1: "invalid command",
// }

// Devices lists available HID devices that advertise the U2F HID protocol.
func getDevices() ([]*DeviceInfo, error) {
	devices, err := Devices()
	if err != nil {
		return nil, err
	}

	res := make([]*DeviceInfo, 0, len(devices))
	for _, d := range devices {
		if d.VendorID == cottonVendorID && d.ProductID == cottonProductID {
			res = append(res, d)
		}
	}

	return res, nil
}

// Open initializes a communication channel with a HID device.
func Open(info *DeviceInfo) (*Device, error) {
	hidDev, err := info.Open()
	if err != nil {
		return nil, err
	}

	d := &Device{
		info:   info,
		device: hidDev,
		readCh: hidDev.ReadCh(),
	}

	return d, nil
}

// A Device is used to communicate with a U2F HID device.
type Device struct {
	ProtocolVersion    uint8
	MajorDeviceVersion uint8
	MinorDeviceVersion uint8
	BuildDeviceVersion uint8

	info   *DeviceInfo
	device *linuxDevice

	mtx    sync.Mutex
	readCh <-chan []byte
	buf    []byte
}

// // Command sends a command and associated data to the device and returns the
// // response.
// func (d *Device) Command(cmd byte, data []byte) ([]byte, error) {
// 	d.mtx.Lock()
// 	defer d.mtx.Unlock()

// 	if len(data)+2 > maxMessageLen {
// 		return nil, fmt.Errorf("device %s Command: command too long", d.info.Product)
// 	}

// 	d.buf = []byte{cmd}
// 	d.buf = append(d.buf, byte(len(data)+2))
// 	d.buf = append(d.buf, data...)

// 	if err := d.device.Write(d.buf); err != nil {
// 		return nil, err
// 	}

// 	return d.readResponse()
// }

// Close closes the device and frees associated resources.
func (d *Device) Close() {
	d.device.Close()
}

const (
	cmdOutFirmHardID   = 0x10
	cmdOutAntennaPower = 0x18
)

// OutFirmHardID sends the hardware command
// data possible values
// 0x00 Firmware
// 0x01 Hardware
func (d *Device) OutFirmHardID(data []byte) ([]byte, error) {
	return d.Command(cmdOutFirmHardID, data)
}

// OutAntennaPower changes the antenna power.
// data possible values
// 0x00 Power OFF
// 0x01 - 0xFE Reserved to change the output level in later versions.
// 0xFF Power ON
func (d *Device) OutAntennaPower(data []byte) ([]byte, error) {
	return d.Command(cmdOutFirmHardID, data)
}

func (d *Device) readResponse() ([]byte, error) {

	timeout := time.After(responseTimeout)

	for {
		select {
		case msg, ok := <-d.readCh:
			if !ok {
				return nil, fmt.Errorf("Read: error reading response, device closed")
			}
			fmt.Printf("%v", msg)
		case <-timeout:
			return nil, fmt.Errorf("Read error: timeout")
		}
	}
}
