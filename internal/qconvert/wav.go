package qconvert

import (
	"errors"
	"fmt"
	"os"

	"git.martianoids.com/queru/retroconverter/internal/cfg"
	"github.com/dustin/go-humanize"
	"github.com/go-audio/wav"
)

// Load input wav file
func (w *Wav) Load(inFileName string) error {
	if cfg.Main.Verbose {
		fmt.Println("> WAV Reading:")
	}

	var err error
	f, err := os.Open(inFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	d := wav.NewDecoder(f)
	if !d.IsValidFile() {
		return errors.New("invalid wav file")
	}

	buf, err := d.FullPCMBuffer()
	if err != nil {
		return err
	}
	w.Data = buf.AsIntBuffer().Data
	if cfg.Main.Verbose {
		fmt.Printf("  [INF] %s read\n", humanize.Bytes(uint64(len(w.Data))))
	}

	w.Format.Duration, _ = d.Duration()
	w.Format.NumChans = d.NumChans
	w.Format.SampleRate = d.SampleRate
	w.Format.BitDepth = d.BitDepth
	if cfg.Main.Verbose {
		fmt.Printf("  [INF] Channels: %d Rate: %d Bits: %d Duration: %s\n",
			w.Format.NumChans,
			w.Format.SampleRate,
			w.Format.BitDepth,
			w.Format.Duration,
		)
	}

	switch w.Format.SampleRate {
	case 22050:
		w.TStates = TS_22050
	case 44100:
		w.TStates = TS_44100
	default:
		return fmt.Errorf("invalid sampling freq: %d (use 22050 or 44100)", w.Format.SampleRate)
	}
	if cfg.Main.Verbose {
		fmt.Printf("  [INF] T-states/sample: %d\n", w.TStates)
	}

	return nil
}
