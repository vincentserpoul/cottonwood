package main

import (
	"fmt"

	"github.com/deadsy/libusb"
)

func cottonwood() int {
	var ctx libusb.Context
	err := libusb.Init(&ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return -1
	}
	defer libusb.Exit(ctx)

	midi_device(ctx, 0x0944, 0x0115) // Korg Nano Key 2
	//midi_device(ctx, 0x041e, 0x3f0e) // Creative Technology, E-MU XMidi1X1 Tab

	return 0

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("\n%s\n", sig)
		quit = true
	}()

	os.Exit(cottonwood())
}