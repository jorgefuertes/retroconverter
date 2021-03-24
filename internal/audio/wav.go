package audio

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/youpy/go-wav"
)

// Wav type
type Wav struct {
	Data    []int
	Highest int
	Factor  float64
	Reader  *wav.Reader
}

func (w *Wav) Load(inputFile string) error {
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// AudioFormat   uint16
	// NumChannels   uint16
	// SampleRate    uint32
	// ByteRate      uint32
	// BlockAlign    uint16
	// BitsPerSample uint16
	w.Reader = wav.NewReader(f)
	format, err := w.Reader.Format()
	if err != nil {
		return err
	}
	fmt.Printf("Format:%v Channels:%d Rate:%v BPS:%d \n",
		format.AudioFormat, format.NumChannels, format.SampleRate, format.BitsPerSample)

	var i int64
	for {
		samples, err := w.Reader.ReadSamples()
		if err == io.EOF {
			break
		}

		for _, sample := range samples {
			w.Data = append(w.Data, sample.Values[0]-128)
			if w.Data[i] > w.Highest {
				w.Highest = w.Data[i]
			}
			i++
		}
	}
	w.Factor = 127 / float64(w.Highest)

	fmt.Printf("Samples:%v Highest:%+d Factor:%.3f\n", i, w.Highest, w.Factor)

	return nil
}

// Normalize volume levels
func (w *Wav) Normalize() {
	var h int = 0
	var l int = 0
	fmt.Print("Normalization ")
	for i := 0; i < len(w.Data); i++ {
		w.Data[i] = int(math.Round(float64(w.Data[i]) * w.Factor))
		if w.Data[i] < -128 {
			w.Data[i] = -128
		} else if w.Data[i] > 127 {
			w.Data[i] = 127
		}
		if w.Data[i] > h {
			h = w.Data[i]
		}
		if w.Data[i] < l {
			l = w.Data[i]
		}
	}
	fmt.Printf("LO:%+d HI:%+d\n", l, h)
}
