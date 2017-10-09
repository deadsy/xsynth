//-----------------------------------------------------------------------------
/*

Attack Decay Sustain Release Envelope

*/
//-----------------------------------------------------------------------------

package main

import "errors"

//-----------------------------------------------------------------------------

type ADSRState int

const (
	idle ADSRState = iota
	attack
	decay
	sustain
	release
)

type ADSR struct {
	ka    float32   // attack constant
	na    int       // attack sample counter
	kd    float32   // decay constant
	nd    int       // decay sample counter
	s     float32   // sustain level
	kr    float32   // release constant
	nr    int       // release sample counter
	state ADSRState // envelope state
	count int       // state counter
	val   float32   // output value
}

func NewADSR(
	a float32, // attack time in seconds
	d float32, // decay time in seconds
	s float32, // sustain level
	r float32, // release time in seconds
	rate int, // sample rate
) (*ADSR, error) {
	e := &ADSR{}

	// sustain level
	if s < 0 || s > 1.0 {
		return nil, errors.New("bad sustain value")
	}
	e.s = s

	if a < 0 {
		return nil, errors.New("bad attack value")
	}
	e.na = int(a * float32(rate))

	if d < 0 {
		return nil, errors.New("bad decay value")
	}
	e.nd = int(d * float32(rate))

	if r < 0 {
		return nil, errors.New("bad release value")
	}
	e.nr = int(r * float32(rate))

	return e, nil
}

func (e *ADSR) Reset() {
	e.val = 0
	e.count = 0
	e.state = idle
}

func (e *ADSR) Start() {
	e.state = attack
}

func (e *ADSR) Stop() {
	e.state = release
}

func (e *ADSR) Sample() float32 {

	switch e.state {
	case idle:
		// idle - do nothing
	case attack:
		// attack exponentially until 1.0
		e.count += 1
		if e.count > e.na {
			e.val = 1.0
			e.count = 0
			e.state = decay
		} else {
			e.val += e.ka * (1.0 - e.val)
		}
	case decay:
		// decay exponentially until sustain
		e.count += 1
		if e.count > e.nd {
			e.val = e.s
			e.count = 0
			e.state = sustain
		} else {
			e.val += e.kd * (e.s - e.val)
		}
	case sustain:
		// sustain - do nothing
	case release:
		// release exponentially until 0.0
		e.count += 1
		if e.count > e.nr {
			e.val = 0
			e.count = 0
			e.state = idle
		} else {
			e.val += e.kr * (0.0 - e.val)
		}
	default:
		panic("bad adsr state")
	}

	return e.val
}

//-----------------------------------------------------------------------------
