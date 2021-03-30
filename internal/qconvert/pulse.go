package qconvert

import (
	"fmt"
)

// BlockStats - Pulse Block listing
func (w *Wav) BlockStats() {
	fmt.Println("> BLOCK LISTING:")
	for i, block := range w.Blocks {
		fmt.Printf("  [%03d/%03d] SAMPLES:%d PULSES:%d PAUSE:%d\n",
			(i + 1), len(w.Blocks), block.SampleCount, len(block.Pulses), block.Pause)
	}
}

func levelPositive(b int) bool {
	return b > TRESHOLD
}

func (w *Wav) getPulse(offset int) Pulse {
	pulse := Pulse{}
	pulse.Positive = (w.Data[offset] > TRESHOLD)
	for i := offset; i < len(w.Data); i++ {
		if (w.Data[i] > TRESHOLD) == pulse.Positive {
			pulse.Duration++
		} else {
			break
		}
	}

	return pulse
}

// ToPulses - Crunch it to pulses
func (w *Wav) ToPulses() {
	block := SampleBlock{}
	var offset int
	for offset < len(w.Data) {
		pulse := w.getPulse(offset)
		block.Pulses = append(block.Pulses, pulse)
		block.SampleCount += pulse.Duration
		offset += int(pulse.Duration)
	}
	w.Blocks = append(w.Blocks, block)
}
