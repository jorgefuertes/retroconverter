package qconvert

import (
	"encoding/binary"
	"fmt"
	"os"
)

func (tzx *Tzx) create() error {
	var err error
	tzx.F, err = os.Create(tzx.Filename)
	return err
}

func (tzx *Tzx) close() error {
	return tzx.F.Close()
}

func (tzx *Tzx) write(b []byte) error {
	c, err := tzx.F.Write(b)
	tzx.Written += c
	return err
}

// SaveTzx - Save to TZX output file
func (w *Wav) SaveTzx(filename string) error {
	tzx := &Tzx{Filename: filename}
	if err := tzx.create(); err != nil {
		return err
	}
	defer tzx.F.Close()

	fmt.Println("--- TZX WRITE ---")
	// header
	fmt.Printf("[%03d] ZXTape!\n", tzx.Written)
	tzx.write([]byte("ZXTape!")) // [7b] signature
	fmt.Printf("[%03d] 0x1A\n", tzx.Written)
	tzx.write([]byte{0x1A}) // [1b] 1A/26 string end
	fmt.Printf("[%03d] 1,20\n", tzx.Written)
	tzx.write([]byte{1, 20}) // [2b] 1,20 Major/minor revision number

	// tzx.write([]byte{30})
	// tzx.write([]byte(tzx.Filename))
	// tzx.write([]byte{0x1A}) // [1b] 1A/26 string end

	// blocks (all ID 15)
	for _, block := range w.Blocks {
		// block ID
		fmt.Printf("[%03d] BlockID: 15\n", tzx.Written)
		tzx.write([]byte{0x15})

		// [2b] T-States per bit
		tstates := []byte{0, 0}
		binary.LittleEndian.PutUint16(tstates, w.TStates)
		fmt.Printf("[%03d] TStates/Sample: %d LSB: %d MSB: %d\n",
			tzx.Written, w.TStates, tstates[0], tstates[1])
		tzx.write(tstates)

		// [2b]Pause after this block in ms
		pause := []byte{0, 0}
		binary.LittleEndian.PutUint16(pause, uint16(block.Pause))
		fmt.Printf("[%03d] Block pause: %d LSB: %d MSB: %d\n",
			tzx.Written, block.Pause, pause[0], pause[1])
		tzx.write(pause)

		// [1b] Used bits in the last byte
		var totalPulses uint
		for _, pulse := range block.Pulses {
			totalPulses += pulse.Duration
		}
		fmt.Printf("[%03d] Used bits in the last byte: %d\n", tzx.Written, (8 - (totalPulses % 8)))
		tzx.write([]byte{byte(8 - (totalPulses % 8))})

		// length 24bits LSB
		totalData := (totalPulses / 8) + (totalPulses % 8)
		u := make([]byte, 4)
		binary.LittleEndian.PutUint32(u, uint32(totalData))
		fmt.Printf("[%03d] Length: %d (%d|%X) LSB:%v %x%x%x\n", tzx.Written,
			totalData, uint32(totalData), uint32(totalData), u[0:3], u[0], u[1], u[2])
		tzx.write(u[0:3]) // [3b] Data len

		// write pulses in byte blocks, LSB
		var currByte byte
		var pos uint
		for _, pulse := range block.Pulses {
			for i := uint(1); i <= pulse.Duration; i++ {
				pos++
				if pulse.Positive {
					currByte |= 1 << (8 - pos)
				}
				if pos == 8 {
					tzx.write([]byte{currByte})
					currByte = 0
					pos = 0
				}
			}
		}
		if totalPulses%8 != 0 {
			tzx.write([]byte{currByte})
		}
	}
	fmt.Printf("TZX data %d bytes written to %s\n", tzx.Written, tzx.Filename)

	return nil
}
