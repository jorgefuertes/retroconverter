package tzx

import (
	"errors"
	"fmt"

	"github.com/dustin/go-humanize"
)

// BlockStats - Pulse Block listing
func (t *TZX) BlockStats() {
	fmt.Println("> BLOCK LISTING:")
	for i, block := range t.Blocks {
		fmt.Printf("  [%03d/%03d] SAMPLES:%s PULSES:%s PAUSE:%d\n",
			(i + 1), len(t.Blocks),
			humanize.Comma(int64(block.SampleCount)),
			humanize.Comma(int64(len(block.Pulses))),
			block.Pause)
	}
}

// return next pulse
func (t *TZX) getPulse(offset int) Pulse {
	pulse := Pulse{}
	pulse.Positive = (t.Data[offset] > TRESHOLD)
	for i := offset; i < len(t.Data); i++ {
		if (t.Data[i] > TRESHOLD) == pulse.Positive {
			pulse.Duration++
		} else {
			break
		}
	}

	return pulse
}

// ToPulses - Crunch it to pulses
func (t *TZX) ToPulses() {
	block := SampleBlock{}
	var offset int
	for offset < len(t.Data) {
		pulse := t.getPulse(offset)
		block.Pulses = append(block.Pulses, pulse)
		block.SampleCount += pulse.Duration
		offset += int(pulse.Duration)
	}
	t.Blocks = append(t.Blocks, block)
}

// DownSample - Downsample 44100->22050->11025
func (t *TZX) DownSample(hz uint32) error {
	if hz != 22050 && hz != 11025 {
		return errors.New("downsample 44100->22050->11025 only")
	}
	factor := t.Format.SampleRate / hz
	if factor < 2 {
		return fmt.Errorf("can't downsample from %d to %d", t.Format.SampleRate, hz)
	}
	t.TStates = t.TStates * uint16(factor)
	t.Format.SampleRate = t.Format.SampleRate / factor

	for b := 0; b < len(t.Blocks); b++ {
		t.Blocks[b].Pause = t.Blocks[b].Pause / uint(factor)
		t.Blocks[b].SampleCount = 0
		for p := 0; p < len(t.Blocks[b].Pulses); p++ {
			t.Blocks[b].Pulses[p].Duration = t.Blocks[b].Pulses[p].Duration / uint(factor)
			t.Blocks[b].SampleCount += t.Blocks[b].Pulses[p].Duration
		}
	}

	return nil
}
