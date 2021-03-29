package qconvert

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/go-audio/wav"
)

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
	fmt.Printf("Channels: %d Rate: %d Bits: %d Duration: %s\n",
		w.Format.NumChans,
		w.Format.SampleRate,
		w.Format.BitDepth,
		w.Format.Duration,
	)

	switch w.Format.SampleRate {
	case 22050:
		w.TStates = TS_22050
	case 44100:
		w.TStates = TS_44100
	default:
		return fmt.Errorf("invalid sampling freq: %d (use 22050 or 44100)", w.Format.SampleRate)
	}
	fmt.Printf("T-states/sample: %d ", w.TStates)

	w.cast2Signed()

	return nil
}

// cast2Signed - Transform to signed numbers
func (w *Wav) cast2Signed() {
	for i := 0; i < len(w.Data); i++ {
		w.Data[i] -= 127
		if w.Data[i] < -127 {
			w.Data[i] = -127
			continue
		}
		if w.Data[i] > 127 {
			w.Data[i] = 127
		}
	}
}

// CalcLevels - Calc the actual HI, LOW and FACTOR
func (w *Wav) CalcLevels() {
	w.Highest = -127
	w.Lowest = 127
	// find the upper volume while copy and cast to signed
	for _, b := range w.Data {
		if b > w.Highest {
			w.Highest = b
		}
		if b < w.Lowest {
			w.Lowest = b
		}
	}
	// calc correction factor needed
	w.Factor = 255 / float64(w.Highest)
}

// FixLevels - Volume normalization
func (w *Wav) FixLevels() {
	for i, b := range w.Data {
		w.Data[i] = int(math.Round(float64(b) * w.Factor))
		if w.Data[i] > 127 {
			w.Data[i] = 127
		} else if w.Data[i] < -127 {
			w.Data[i] = -127
		}
	}
	w.CalcLevels()
}

// Stats - Calc and display an histogram
func (w *Wav) Stats() {
	stats := make([]int, 256)
	for _, b := range w.Data {
		stats[b+127]++
	}

	fmt.Println("\nSTATS:")
	for i, v := range stats {
		if v != 0 {
			fmt.Printf("%+ 4d: %d\n", i-127, v)
		}
	}
}
