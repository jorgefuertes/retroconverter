package qconvert

import (
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/go-audio/wav"
)

// Wav file data
type Wav struct {
	Format struct {
		NumChans   uint16
		SampleRate uint32
		BitDepth   uint16
		Duration   time.Duration
	}
	Highest int
	Factor  float64
	Data    []int
}

// Load input wav file
func (w *Wav) Load(inFileName string) error {
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

	w.Format.Duration, _ = d.Duration()
	w.Format.NumChans = d.NumChans
	w.Format.SampleRate = d.SampleRate
	w.Format.BitDepth = d.BitDepth

	w.calcAudioValues()
	w.fixLevels()

	return nil
}

func (w *Wav) calcAudioValues() {
	// find the upper volume while copy and cast to signed
	for i, b := range w.Data {
		w.Data[i] = b - 127
		if w.Data[i] > w.Highest {
			w.Highest = w.Data[i]
		}
	}
	// calc correction factor needed
	w.Factor = 255 / float64(w.Highest)

	fmt.Printf("Channels: %d Rate: %d Bits: %d Duration: %s\n",
		w.Format.NumChans,
		w.Format.SampleRate,
		w.Format.BitDepth,
		w.Format.Duration,
	)

	fmt.Printf("Volume highest: %+d Correction factor: %.3f\n", w.Highest, w.Factor)
}

func (w *Wav) fixLevels() {
	for i, b := range w.Data {
		w.Data[i] = int(math.Round(float64(b) * w.Factor))
		if w.Data[i] > 127 {
			w.Data[i] = 127
		} else if w.Data[i] < -128 {
			w.Data[i] = -128
		}
	}
}

func (w *Wav) Stats() {
	var percent int
	var last int
	stats := make([]int, 256)
	fmt.Printf("Calc stats…%d%%", percent)
	for i, b := range w.Data {
		percent = (i / len(w.Data)) * 100
		if percent != last {
			fmt.Printf("\rCalc stats…%d%%", percent)
			last = percent
		}
		stats[b+128]++
	}

	fmt.Println("\nSTATS:")
	for i, v := range stats {
		if v != 0 {
			fmt.Printf("%+03d: %d\n", i-127, v)
		}
	}
}
