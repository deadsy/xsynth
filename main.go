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

const SAMPLE_RATE = 22500

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
	t0 := NewLUT_Sine(261.63, SAMPLE_RATE)
	t1 := NewLUT_Sine(329.63, SAMPLE_RATE)
	t2 := NewLUT_Sine(392.00, SAMPLE_RATE)

	//t := NewLUT_Sawtooth(440.0, SAMPLE_RATE)
	//t := NewLUT_Square(440.0, SAMPLE_RATE)

	for {
		for i, _ := range samples {
			y := t0.Sample() + t1.Sample()*t2.Sample()
			samples[i] = y * amp
		}
		st.Write(samples, pulsego.SEEK_RELATIVE)
	}
}

//-----------------------------------------------------------------------------

func main2() {

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

func main() {

	e, err := NewADSR(0.01, 0.05, 0.5, 0.04, 600)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	for i := 0; i < 100; i++ {

		if i == 5 {
			e.Attack()
		}

		if i == 70 {
			e.Release()
		}

		fmt.Printf("%d, %f\n", i, e.Sample())
		//fmt.Printf("%d: %s %f\n", i, e.state, e.Sample())
	}

}

//-----------------------------------------------------------------------------
