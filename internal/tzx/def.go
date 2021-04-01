package tzx

import (
	"os"
	"time"
)

const TRESHOLD int = 127
const UINT24_MAX uint = 16_777_216
const TS_11025 uint16 = 316
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

// Tzx data
type TZX struct {
	InFilename  string
	OutFilename string
	OutFile     *os.File
	Written     int
	Format      struct {
		NumChans   uint16
		SampleRate uint32
		BitDepth   uint16
		Duration   time.Duration
	}
	TStates uint16
	Data    []int
	Blocks  []SampleBlock
}
