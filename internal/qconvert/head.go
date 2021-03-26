package qconvert

import "time"

const NEGATIVE int8 = 0
const POSITIVE int8 = 1
const HUSH int8 = 2
const pulseNames = "NEGATIVE,POSITIVE,HUSH"

type Pulse struct {
	Positive bool
	Duration int
}

type Block struct {
	Pause  int // Pause after milliseconds
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
