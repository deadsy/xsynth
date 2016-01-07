//-----------------------------------------------------------------------------
/*

MIDI Keyboard Driver

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"github.com/deadsy/libusb"
)

//-----------------------------------------------------------------------------

// Korg NanoKey2
const usb_vid = 0x0944
const usb_pid = 0x0115

//-----------------------------------------------------------------------------

func midi_init() error {
	var ctx libusb.Context
	err := libusb.Init(&ctx)
	if err != nil {
		return err
	}
	defer libusb.Exit(ctx)

	list, err := libusb.Get_Device_List(ctx)
	if err != nil {
		return err
	}
	defer libusb.Free_Device_List(list, 1)

	for _, dev := range list {
		dd, err := libusb.Get_Device_Descriptor(dev)
		if err != nil {
			return err
		}
		if dd.IdVendor == usb_vid && dd.IdProduct == usb_pid {
			fmt.Printf("found\n")
			break
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
