package qconvert

import (
	"encoding/binary"
	"fmt"
	"os"
)

var SIGNATURE []byte = []byte("ZXTape!")
var REVISION []byte = []byte{1, 20}

// SaveTzx - Save to TZX output file
func (w *Wav) SaveTzx(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// header
	f.Write(SIGNATURE)    // [7b] signature
	f.Write([]byte{0x1A}) // [1b] string end
	f.Write(REVISION)     // [2b] Major/minor revision number

	// t-states calc
	tstates := []byte{0, 0}
	binary.LittleEndian.PutUint16(
		tstates, uint16(0.5+(3500000.0/float64(w.Format.SampleRate))))

	// blocks (all ID 15)
	for _, block := range w.Blocks {
		f.Write(tstates) // [2b] T-States per bit
		// pause
		pauseMs := []byte{0, 0}
		binary.LittleEndian.PutUint16(pauseMs, uint16(block.Pause))
		f.Write(pauseMs) // [2b ]Pause after this block in ms
		// alignment
		var totalPulses uint
		for _, pulse := range block.Pulses {
			totalPulses += pulse.Duration
		}
		f.Write([]byte{byte(totalPulses % 8)}) // [1b] Used bits in the last byte
		// length 24bits LSB
		u := make([]byte, 4)
		binary.LittleEndian.PutUint32(u, uint32(totalPulses))
		f.Write(u[0:3]) // [3b] Data len
		var currByte byte
		var bitCt int8
		for _, pulse := range block.Pulses {
			for i := 0; i < int(pulse.Duration); i++ {
				bitCt++
				if pulse.Positive {
					currByte |= (1 << (8 - bitCt))
				}

				if bitCt == 8 {
					f.Write([]byte{currByte}) // Write a byte of pulses
					bitCt = 0
					currByte = 0
				}
			}
			// check for remaining bits
			if bitCt < 8 {
				f.Write([]byte{currByte}) // Write last and incomplete byte
			}
		}
	}
	fmt.Println("TZX data written to:", f.Name())

	return nil
}
