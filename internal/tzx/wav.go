package tzx

import (
	"errors"
	"fmt"
	"os"

	"git.martianoids.com/queru/retroconverter/internal/cfg"
	"github.com/dustin/go-humanize"
	"github.com/go-audio/wav"
)

// Load input wav file
func (t *TZX) Load(inFileName string) error {
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
	t.Data = buf.AsIntBuffer().Data
	if cfg.Main.Verbose {
		fmt.Printf("  [INF] %s read\n", humanize.Bytes(uint64(len(t.Data))))
	}

	t.Format.Duration, _ = d.Duration()
	t.Format.NumChans = d.NumChans
	t.Format.SampleRate = d.SampleRate
	t.Format.BitDepth = d.BitDepth
	if cfg.Main.Verbose {
		fmt.Printf("  [INF] Channels: %d Rate: %d Bits: %d Duration: %s\n",
			t.Format.NumChans,
			t.Format.SampleRate,
			t.Format.BitDepth,
			t.Format.Duration,
		)
	}

	switch t.Format.SampleRate {
	case 11025:
		t.TStates = TS_11025
	case 22050:
		t.TStates = TS_22050
	case 44100:
		t.TStates = TS_44100
	default:
		return fmt.Errorf("invalid sampling freq: %d (use 11025, 22050 or 44100)",
			t.Format.SampleRate)
	}
	if cfg.Main.Verbose {
		fmt.Printf("  [INF] T-states/sample: %d\n", t.TStates)
		fmt.Printf("  [INF] Sample count: %s\n", humanize.Comma(int64(len(t.Data))))
	}

	return nil
}
