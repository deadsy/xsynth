//-----------------------------------------------------------------------------
/*

MIDI Keyboard Driver

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"math"

	"github.com/deadsy/libusb"
)

//-----------------------------------------------------------------------------

const NOTES_IN_OCTAVE = 12
const MIDI_NOTE_A5 = 69
const A5_FREQUENCY = 440

// return the frequency of the midi note
func midi_to_frequency(note uint) float32 {
	return float32(A5_FREQUENCY * math.Pow(2.0, float64(int(note)-MIDI_NOTE_A5)/NOTES_IN_OCTAVE))
}

// return an octave number
func midi_to_octave(note uint) uint {
	return note / NOTES_IN_OCTAVE
}

//-----------------------------------------------------------------------------

func major_chord(root uint) [3]uint {
	return [3]uint{root, root + 4, root + 7}
}

func minor_chord(root uint) [3]uint {
	return [3]uint{root, root + 3, root + 7}
}

//-----------------------------------------------------------------------------

// Korg NanoKey2
const usb_vid = 0x0944
const usb_pid = 0x0115

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
