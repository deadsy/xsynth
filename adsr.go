//-----------------------------------------------------------------------------
/*

Attack Decay Sustain Release Envelope

*/
//-----------------------------------------------------------------------------

package main

import (
	"errors"
	"math"
)

//-----------------------------------------------------------------------------

// Example: Suppose we want an attack that goes from 0.0 to 1.0 in 0.1 seconds.
// An exponential rise will only approach 1.0, so we pick a level slightly
// below this to be achieved in the required time. The trigger_epsilon value
// controls this.
const trigger_epsilon = 0.01

// Return a k value to give the exponential rise/fall in the required time.
func get_k(t float32, rate int) float32 {
	if t <= 0 {
		return 1.0
	}
	return float32(1.0 - math.Exp(math.Log(trigger_epsilon)/(float64(t)*float64(rate))))
}

//-----------------------------------------------------------------------------

type ADSRState int

const (
	idle ADSRState = iota
	attack
	decay
	sustain
	release
)

var state_txt = map[ADSRState]string{
	idle:    "idle",
	attack:  "attack",
	decay:   "decay",
	sustain: "sustain",
	release: "release",
}

func (x ADSRState) String() string {
	return state_txt[x]
}

//-----------------------------------------------------------------------------

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
	e.ka = get_k(a, rate)
	e.kd = get_k(d, rate)
	e.kr = get_k(r, rate)

	e.d_trigger = 1.0 - trigger_epsilon
	e.s_trigger = e.s + (1.0-e.s)*trigger_epsilon
	e.i_trigger = e.s * trigger_epsilon

	return e, nil
}

//-----------------------------------------------------------------------------

// Enter attack state.
func (e *ADSR) Attack() {
	e.state = attack
}

// Enter release state.
func (e *ADSR) Release() {
	if e.state != idle {
		e.state = release
	}
}

// Reset to idle.
func (e *ADSR) Idle() {
	e.val = 0
	e.state = idle
}

//-----------------------------------------------------------------------------

// Return a sample value for the ADSR envelope.
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
