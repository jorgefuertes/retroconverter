package qconvert

import (
	"time"
)

const NEGATIVE int8 = 0
const POSITIVE int8 = 1
const HUSH int8 = 2
const UINT24_MAX uint = 16777216
const TS_22050 uint16 = 158
const TS_44100 uint16 = 59

type Pulse struct {
	Positive bool
	Duration uint
}

type SampleBlock struct {
	Pause  uint // Pause after milliseconds
	Pulses []Pulse
}

// Wav file data
type Wav struct {
	Format struct {
		NumChans   uint16
		SampleRate uint32
		BitDepth   uint16
		Duration   time.Duration
	}
	Filters struct {
		PHi int
		PLo int
		SHi int
		SLo int
	}
	TStates uint16
	Highest int
	Lowest  int
	Factor  float64
	Data    []int
	Blocks  []SampleBlock
}
