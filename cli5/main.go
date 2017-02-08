package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/flynn/hid"
)

func main() {
	devices, err := Devices()
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

		res, err := dev.OutAntennaPower([]byte{0xff})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("message %v", res)
	}
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
func Devices() ([]*hid.DeviceInfo, error) {
	devices, err := hid.Devices()
	if err != nil {
		return nil, err
	}

	res := make([]*hid.DeviceInfo, 0, len(devices))
	for _, d := range devices {
		if d.VendorID == cottonVendorID && d.ProductID == cottonProductID {
			res = append(res, d)
		}
	}

	return res, nil
}

// Open initializes a communication channel with a U2F HID device.
func Open(info *hid.DeviceInfo) (*Device, error) {
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

	info   *hid.DeviceInfo
	device hid.Device

	mtx    sync.Mutex
	readCh <-chan []byte
	buf    []byte
}

// Command sends a command and associated data to the device and returns the
// response.
func (d *Device) Command(cmd byte, data []byte) ([]byte, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if len(data)+2 > maxMessageLen {
		return nil, fmt.Errorf("device %s Command: command too long", d.info.Product)
	}

	d.buf = []byte{cmd}
	d.buf = append(d.buf, byte(len(data)+2))
	d.buf = append(d.buf, data...)

	if err := d.device.Write(d.buf); err != nil {
		return nil, err
	}

	return d.readResponse()
}

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
