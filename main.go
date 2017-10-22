//-----------------------------------------------------------------------------
/*

XSynth

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/xsynth/pulsego"
)

//-----------------------------------------------------------------------------

const SAMPLE_RATE = 44100

//-----------------------------------------------------------------------------

func sine_wave(pa *pulsego.PulseMainLoop) {

	ctx := pa.NewContext("default", 0)
	if ctx == nil {
		fmt.Println("Failed to create a new context")
		return
	}
	defer ctx.Dispose()
	st := ctx.NewStream("default", &pulsego.PulseSampleSpec{
		Format: pulsego.SAMPLE_FLOAT32LE, Rate: SAMPLE_RATE, Channels: 1})
	if st == nil {
		fmt.Println("Failed to create a new stream")
		return
	}
	defer st.Dispose()
	st.ConnectToSink()

	samples := make([]float32, 64)
	amp := float32(0.1)
	chord := major_chord(60)
	t0 := NewLUT_Sine(midi_to_frequency(chord[0]), SAMPLE_RATE)
	t1 := NewLUT_Sine(midi_to_frequency(chord[1]), SAMPLE_RATE)
	t2 := NewLUT_Sine(midi_to_frequency(chord[2]), SAMPLE_RATE)

	//t := NewLUT_Sawtooth(440.0, SAMPLE_RATE)
	//t := NewLUT_Square(440.0, SAMPLE_RATE)

	for {
		for i, _ := range samples {
			y := t0.Sample() + t1.Sample() + t2.Sample()
			samples[i] = y * amp
		}
		st.Write(samples, pulsego.SEEK_RELATIVE)
	}
}

//-----------------------------------------------------------------------------

func main() {

	midi_init()

	pa := pulsego.NewPulseMainLoop()
	defer pa.Dispose()
	pa.Start()

	done := make(chan bool)
	go func() {
		sine_wave(pa)
		done <- true
	}()
	<-done
	close(done)
}

//-----------------------------------------------------------------------------

func mainx() {

	e, err := NewADSR_Envelope(0.04, 0.05, 0.5, 0.04, 600)
	//e, err := NewAD_Envelope(0.04, 0.05, 600)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	for i := 0; i < 200; i++ {

		if i == 5 {
			e.Attack()
		}

		if i == 100 {
			e.Release()
		}

		fmt.Printf("%d: %s %f\n", i, e.state, e.Sample())
	}

}

//-----------------------------------------------------------------------------
