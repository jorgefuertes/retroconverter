package qconvert

import (
	"time"
)

const TRESHOLD int = 127
const NEGATIVE int8 = 0
const POSITIVE int8 = 1
const UINT24_MAX uint = 16777216
const TS_22050 uint16 = 158
const TS_44100 uint16 = 79

type Pulse struct {
	Positive bool
	Duration uint
}

type SampleBlock struct {
	Pause       uint // Pause after milliseconds
	Pulses      []Pulse
	SampleCount uint
}

// Wav file data
type Wav struct {
	Format struct {
		NumChans   uint16
		SampleRate uint32
		BitDepth   uint16
		Duration   time.Duration
	}
	TStates uint16
	Data    []int
	Blocks  []SampleBlock
}
