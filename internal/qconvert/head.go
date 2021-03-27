package qconvert

import "time"

const NEGATIVE int8 = 0
const POSITIVE int8 = 1
const HUSH int8 = 2
const UINT24_MAX uint = 16777216

type Pulse struct {
	Positive bool
	Duration uint
}

type Block struct {
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
	Highest int
	Lowest  int
	Factor  float64
	Data    []int
	Blocks  []Block
}
