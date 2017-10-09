//-----------------------------------------------------------------------------
/*

Attack Decay Sustain Release Envelope

*/
//-----------------------------------------------------------------------------

package main

import "errors"

//-----------------------------------------------------------------------------

const TRIGGER_LEVEL = 0.02

type ADSRState int

const (
	idle ADSRState = iota
	attack
	decay
	sustain
	release
)

type ADSR struct {
	s         float32   // sustain level
	ka        float32   // attack constant
	kd        float32   // decay constant
	kr        float32   // release constant
	d_trigger float32   // attack->decay trigger level
	s_trigger float32   // decay->sustain trigger level
	i_trigger float32   // release->idle trigger level
	state     ADSRState // envelope state
	val       float32   // output value
}

func NewADSR(
	a float32, // attack time in seconds
	d float32, // decay time in seconds
	s float32, // sustain level
	r float32, // release time in seconds
	rate int, // sample rate
) (*ADSR, error) {
	e := &ADSR{}

	if a < 0 {
		return nil, errors.New("bad attack time")
	}
	if d < 0 {
		return nil, errors.New("bad decay time")
	}
	if s < 0 || s > 1.0 {
		return nil, errors.New("bad sustain level")
	}
	if r < 0 {
		return nil, errors.New("bad release time")
	}

	e.s = s
	e.d_trigger = 1.0 - TRIGGER_LEVEL
	e.s_trigger = e.s + (1.0-e.s)*TRIGGER_LEVEL
	e.i_trigger = e.s * TRIGGER_LEVEL

	return e, nil
}

// Enter attack state.
func (e *ADSR) Start() {
	e.state = attack
}

// Enter release state.
func (e *ADSR) Release() {
	if e.state != idle {
		e.state = release
	}
}

// Rest to idle.
func (e *ADSR) Stop() {
	e.val = 0
	e.state = idle
}

func (e *ADSR) Sample() float32 {

	switch e.state {
	case idle:
		// idle - do nothing
	case attack:
		// attack until 1.0 level
		if e.val < e.d_trigger {
			e.val += e.ka * (1.0 - e.val)
		} else {
			e.val = 1
			e.state = decay
		}
	case decay:
		// decay until sustain level
		if e.val > e.s_trigger {
			e.val += e.kd * (e.s - e.val)
		} else {
			e.val = e.s
			e.state = sustain
		}
	case sustain:
		// sustain - do nothing
	case release:
		// release until idle level
		if e.val > e.i_trigger {
			e.val += e.kr * (0.0 - e.val)
		} else {
			e.val = 0
			e.state = idle
		}
	default:
		panic("bad adsr state")
	}

	return e.val
}

//-----------------------------------------------------------------------------
