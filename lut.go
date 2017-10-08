//-----------------------------------------------------------------------------
/*

Lookup Tables for Wave Generation

*/
//-----------------------------------------------------------------------------

package main

import "math"

//-----------------------------------------------------------------------------

var cos_table []float32
var sawtooth_table []float32
var square_table []float32

func init() {
	var n int

	n = 512
	cos_table = make([]float32, n)
	for i := range cos_table {
		cos_table[i] = float32(math.Cos(float64(i) * 2.0 * math.Pi / float64(n)))
	}

	n = 512
	sawtooth_table = make([]float32, n)
	k := float32(2.0 / (float32(n) - 1.0))
	for i := range sawtooth_table {
		sawtooth_table[i] = (k * float32(i)) - 1.0
	}

	n = 128
	square_table = make([]float32, n)
	for i := 0; i < n; i++ {
		if i <= n/2 {
			square_table[i] = -1.0
		} else {
			square_table[i] = 1.0
		}
	}

}

//-----------------------------------------------------------------------------

type LUT struct {
	table  []float32
	xrange float32
	x      float32
	step   float32
}

func (t *LUT) SetTable(table []float32) {
	t.table = table
	t.xrange = float32(len(t.table))
}

func (t *LUT) SetStep(f float32, rate int) {
	t.step = f * float32(len(t.table)) / float32(rate)
}

func (t *LUT) Sample() float32 {
	// linear interpolation
	x0 := int(math.Floor(float64(t.x)))
	y0 := t.table[x0]
	var y1 float32
	if x0 == len(t.table)-1 {
		y1 = t.table[0]
	} else {
		y1 = t.table[x0+1]
	}
	y := y0 + ((t.x - float32(x0)) * (y1 - y0))
	// step the x position
	t.x += t.step
	if t.x >= t.xrange {
		t.x -= t.xrange
	}
	return y
}

//-----------------------------------------------------------------------------

func NewLUT_Sine(f float32, rate int) *LUT {
	t := &LUT{}
	t.SetTable(cos_table)
	t.SetStep(f, rate)
	return t
}

func NewLUT_Sawtooth(f float32, rate int) *LUT {
	t := &LUT{}
	t.SetTable(sawtooth_table)
	t.SetStep(f, rate)
	return t
}

func NewLUT_Square(f float32, rate int) *LUT {
	t := &LUT{}
	t.SetTable(square_table)
	t.SetStep(f, rate)
	return t
}

//-----------------------------------------------------------------------------
