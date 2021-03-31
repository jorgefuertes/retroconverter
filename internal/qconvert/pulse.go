package qconvert

import (
	"errors"
	"fmt"

	"github.com/dustin/go-humanize"
)

// BlockStats - Pulse Block listing
func (w *Wav) BlockStats() {
	fmt.Println("> BLOCK LISTING:")
	for i, block := range w.Blocks {
		fmt.Printf("  [%03d/%03d] SAMPLES:%s PULSES:%s PAUSE:%d\n",
			(i + 1), len(w.Blocks),
			humanize.Comma(int64(block.SampleCount)),
			humanize.Comma(int64(len(block.Pulses))),
			block.Pause)
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

// DownSample - Downsample 44100->22050->11025
func (w *Wav) DownSample(hz uint32) error {
	if hz != 22050 && hz != 11025 {
		return errors.New("I can downsample 44100->22050->11025 only")
	}
	factor := w.Format.SampleRate / hz
	if factor < 2 {
		return fmt.Errorf("I can't downsample from %d to %d", w.Format.SampleRate, hz)
	}
	w.TStates = w.TStates * uint16(factor)
	w.Format.SampleRate = w.Format.SampleRate / factor

	for b, _ := range w.Blocks {
		w.Blocks[b].Pause = w.Blocks[b].Pause / uint(factor)
		w.Blocks[b].SampleCount = 0
		for p, _ := range w.Blocks[b].Pulses {
			w.Blocks[b].Pulses[p].Duration = w.Blocks[b].Pulses[p].Duration / uint(factor)
			w.Blocks[b].SampleCount += w.Blocks[b].Pulses[p].Duration
		}
	}

	return nil
}
