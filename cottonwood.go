package cottonwood

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/vincentserpoul/cottonwood/hid"
	"github.com/vincentserpoul/cottonwood/ioctl"
)

// A Device is used to communicate with a U2F HID device.
type Device struct {
	Info   *hid.DeviceInfo
	handle *os.File
}

const (
	cottonVendorID  = 0x1325
	cottonProductID = 0xc029

	messageLen = 64

	responseTimeout = 3 * time.Second
)

// GetDevices lists available HID cottonwood devices
func GetDevices() ([]*Device, error) {
	devices, err := hid.Devices()
	if err != nil {
		return nil, err
	}

	res := make([]*Device, 0, len(devices))
	for _, d := range devices {
		if d.VendorID == cottonVendorID && d.ProductID == cottonProductID {
			res = append(res, &Device{Info: d})
		}
	}

	return res, nil
}

// Close closes the device and frees associated resources.
func (dev *Device) Close() error {
	return dev.handle.Close()
}

// Open the device saves the file handle to struct
func (dev *Device) Open() (err error) {
	dev.handle, err = dev.Info.Open()
	return err
}

// command will run a report on the device
func (dev *Device) Command(reportID byte, payload []byte) ([]byte, error) {

	/* Set Feature */
	bufS := make([]byte, messageLen)
	bufS[0] = reportID
	bufS[1] = byte(len(payload) + 2)
	for i, payloadItm := range payload {
		bufS[i+2] = payloadItm
	}

	fd := dev.handle.Fd()
	err := ioctl.Ioctl(fd, hid.IoctlHIDIOCSFEATURE(messageLen), uintptr(unsafe.Pointer(&bufS[0])))
	if err != nil {
		return nil, err
	}

	// wait a little
	time.Sleep(100 * time.Millisecond)

	bufG := make([]byte, 64)
	bufG[0] = reportID
	err = ioctl.Ioctl(fd, hid.IoctlHIDIOCGFEATURE(messageLen), uintptr(unsafe.Pointer(&bufG[0])))
	if err != nil {
		return nil, err
	}

	return bufG, nil
}

// command will run a report on the device
func (dev *Device) command(reportID byte, payload []byte, waitingTime time.Duration) ([]byte, error) {

	/* Set Feature */
	bufS := make([]byte, messageLen)
	bufS[0] = reportID
	bufS[1] = byte(len(payload) + 2)
	for i, payloadItm := range payload {
		bufS[i+2] = payloadItm
	}

	fd := dev.handle.Fd()
	err := ioctl.Ioctl(fd, hid.IoctlHIDIOCSFEATURE(messageLen), uintptr(unsafe.Pointer(&bufS[0])))
	if err != nil {
		return nil, err
	}

	// wait a little
	time.Sleep(waitingTime)

	bufG := make([]byte, 64)
	bufG[0] = reportID
	err = ioctl.Ioctl(fd, hid.IoctlHIDIOCGFEATURE(messageLen), uintptr(unsafe.Pointer(&bufG[0])))
	if err != nil {
		return nil, err
	}

	return bufG, nil
}

const (
	cmdOutFirmHardID    = 0x10
	cmdOutAntennaPower  = 0x18
	cmdOutInventory     = 0x31
	cmdOutInventoryRSSI = 0x43
)

// OutFirmHardID sends the hardware command
// data possible values
// 0x00 Firmware
// 0x01 Hardware
func (dev *Device) OutFirmHardID(data []byte) ([]byte, error) {
	return dev.command(cmdOutFirmHardID, data, 100*time.Millisecond)
}

// OutAntennaPower changes the antenna power.
// data possible values
// 0x00 Power OFF
// 0x01 - 0xFE Reserved to change the output level in later versions.
// 0xFF Power ON
func (dev *Device) OutAntennaPower(data []byte) ([]byte, error) {
	return dev.command(cmdOutAntennaPower, data, 100*time.Millisecond)
}

// outInventory ask for inventory
// data possible values
// 0x01 Start inventory round
// 0x02 Next Tag information (should not be sent anymore since v1.3.0)
//
// RESPONSE
// Byte 0: report ID
// Byte 1: Frame length
// Byte 2: Number of found tags
// Byte 3: Length of EPC byte
// Byte 4-Byte xx:  EPC 1…x
// Bytexx+1..Byte: EPC 1…x rfu
func (dev *Device) outInventory(data []byte) ([]byte, error) {
	return dev.command(cmdOutInventory, data, 100*time.Millisecond)
}

// outInventory ask for inventory
// data possible values
// 0x01 Start inventory round
// 0x02 Next Tag information (should not be sent anymore since v1.3.0)
//
// RESPONSE
// Byte 0: report ID
// Byte 1: Frame length
// Byte 2: Number of found tags
// Byte 3: Length of EPC byte
// Byte 4-Byte xx:  EPC 1…x
// Bytexx+1..Byte: EPC 1…x rfu
func (dev *Device) outInventoryRSSI(data []byte) ([]byte, error) {
	return dev.command(cmdOutInventory, data, 100*time.Millisecond)
}

// Tag represents a tag
type Tag struct {
	ID  []byte
	RFU []byte
}

// OutInventory get the inventory
func (dev *Device) OutInventory() ([]Tag, error) {
	response, err := dev.outInventory([]byte{0x01})
	if err != nil {
		return nil, err
	}

	if response[0] != 0x32 {
		return nil, fmt.Errorf("OutInventory: %X", response)
	}

	tagCount := int(response[2])
	if tagCount == 0 {
		return nil, nil
	}

	tags := make([]Tag, tagCount)

	frameLen := int(response[1])
	tagLen := int(response[3])
	tags[0] = Tag{
		ID:  response[6 : 6+tagLen],
		RFU: response[6+tagLen : frameLen],
	}

	return tags, nil

}
