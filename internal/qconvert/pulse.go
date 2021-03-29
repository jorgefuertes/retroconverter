package qconvert

import (
	"errors"
	"fmt"
)

// pulseType - Given a sample level returns the signal 0,1 or 2
func (w *Wav) pulseType(b int) int8 {
	// hush
	if b >= w.Filters.SLo && b <= w.Filters.SHi {
		return HUSH
	}

	// positive
	if b >= w.Filters.PHi {
		return POSITIVE
	}

	// negative
	return NEGATIVE
}

// BlockStats - Pulse Block listing
func (w *Wav) BlockStats() {
	fmt.Println("BLOCK LISTING:")
	for i, block := range w.Blocks {
		var loCt int
		var hiCt int
		for _, pulse := range block.Pulses {
			if pulse.Positive {
				hiCt++
			} else {
				loCt++
			}
		}
		fmt.Printf("[BLOCK %04d/%04d] LO_PULSES:%d HI_PULSES:%d TOTAL:%d PAUSE_AFTER:%d\n",
			i, len(w.Blocks), loCt, hiCt, len(block.Pulses), block.Pause)
	}
}

// CalcPulse - Crunch it to pulses
func (w *Wav) CalcPulse() error {
	var signal int8 = HUSH // default to hush
	var count uint
	var pulseCount uint
	block := SampleBlock{}
	for _, b := range w.Data {
		if w.pulseType(b) == signal {
			// No change, just increase pulse length
			count++
			continue
		}

		// Pulse change from last signal
		switch signal {
		case HUSH:
			// Was HUSH, close block with pause if significant and block in course
			if count >= 10 && len(w.Blocks) > 0 {
				block.Pause = count
				// Closing block
				w.Blocks = append(w.Blocks, block)
				block = SampleBlock{}
			} else {
				// Begin with new signal discarding short hush
				if len(block.Pulses) > 0 {
					if pulseCount > UINT24_MAX {
						return errors.New("24bits pulse limit reached")
					}
					block.Pulses[len(block.Pulses)-1].Duration += count
				}
				// Just discard
			}
			// Begin with new signal
			signal = w.pulseType(b)
			count = 1
			continue
		default:
			// Change from last signal
			// Append the pulse
			block.Pulses = append(block.Pulses, Pulse{(w.pulseType(b) == POSITIVE), count})
			// Begin with new signal
			signal = w.pulseType(b)
			count = 1
		}
	}
	// Close last block if not empty
	if len(block.Pulses) > 0 || block.Pause > 0 {
		w.Blocks = append(w.Blocks, block)
	}

	return nil
}
