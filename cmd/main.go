package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vincentserpoul/cottonwood"
)

func main() {
	devices, err := cottonwood.GetDevices()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(devices))
	for _, dev := range devices {
		errOpen := dev.Open()
		if err != nil {
			log.Fatal(errOpen)
		}
		defer func() {
			errClose := dev.Close()
			if err != nil {
				log.Printf("error closing device %s: %v", dev.Info.Path, errClose)
			}
		}()
		fmt.Println("opened", dev.Info.Path)

		resp, err := dev.OutFirmHardID([]byte{0x00})
		if err != nil {
			log.Fatal(err)
		}

		for i, charact := range resp {
			if i > 1 && charact != 0 {
				fmt.Printf("%c", charact)
			}
		}
		fmt.Print("\n")

		resp, err = dev.OutFirmHardID([]byte{0x01})
		if err != nil {
			log.Fatal(err)
		}

		for i, charact := range resp {
			if i > 1 && charact != 0 {
				fmt.Printf("%c", charact)
			}
		}
		fmt.Print("\n")

		resp, err = dev.OutAntennaPower([]byte{0xff})
		if err != nil {
			log.Fatal(err)
		}

		for i, charact := range resp {
			if i > 1 && charact != 0 {
				fmt.Printf("%c", charact)
			}
		}
		fmt.Print("\n")

		for {
			tags, errInv := dev.OutInventory()
			if errInv != nil {
				log.Println(errInv)
			}
			fmt.Printf("%d tags found: ", len(tags))

			if len(tags) > 0 {
				for _, charact := range tags[0].ID {
					fmt.Printf("0x%X ", charact)
				}
			}
			fmt.Println()
			time.Sleep(2 * time.Second)

		}

		// for {
		// 	fmt.Println("New Scan")
		// 	respInv, errInv := dev.Command(0x31, []byte{0x01})
		// 	if errInv != nil {
		// 		log.Fatal(errInv)
		// 	}

		// 	for _, charact := range respInv {
		// 		fmt.Printf("0x%X ", charact)
		// 	}
		// 	fmt.Print("\n")

		// 	// respInv2, errInv2 := dev.Command(0x43, []byte{0x02})
		// 	// if errInv2 != nil {
		// 	// 	log.Fatal(errInv2)
		// 	// }

		// 	// for _, charact := range respInv2 {
		// 	// 	fmt.Printf("0x%X ", charact)
		// 	// }
		// 	// fmt.Print("\n")

		// 	time.Sleep(2 * time.Second)
		// }

	}
}
