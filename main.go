//-----------------------------------------------------------------------------
/*

XSynth

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"github.com/deadsy/xsynth/pulsego"
	"math"
)

//-----------------------------------------------------------------------------

const TABLE_SIZE = (2 << 16)
const SAMPLE_RATE = 22500

var cos_table = make([]float32, TABLE_SIZE)

func init_tables() {
	for i, _ := range cos_table {
		cos_table[i] = float32(math.Cos(float64(i) * 2.0 * math.Pi / float64(TABLE_SIZE)))
	}
}

func freq_to_step(f float64) int {
	return int(math.Ceil(f * float64(TABLE_SIZE) / float64(SAMPLE_RATE)))
}

func step_to_freq(step int) float64 {
	return float64(step) * float64(SAMPLE_RATE) / float64(TABLE_SIZE)
}

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

	amp := float32(0.8)
	posn := 0
	step := freq_to_step(440.0)
	fmt.Printf("actual freq is %f\n", step_to_freq(step))

	for {
		for i, _ := range samples {
			samples[i] = cos_table[posn] * amp
			posn = (posn + step) % TABLE_SIZE
		}
		st.Write(samples, pulsego.SEEK_RELATIVE)
	}
}

//-----------------------------------------------------------------------------

func main() {

	init_tables()
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
