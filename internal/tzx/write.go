package tzx

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"git.martianoids.com/queru/retroconverter/internal/build"
	"git.martianoids.com/queru/retroconverter/internal/cfg"
	"github.com/dustin/go-humanize"
)

func (t *TZX) create() error {
	var err error
	t.OutFile, err = os.Create(t.OutFilename)
	return err
}

func (t *TZX) close() error {
	return t.OutFile.Close()
}

func (t *TZX) write(b []byte) error {
	c, err := t.OutFile.Write(b)
	t.Written += c
	return err
}

// SaveTzx - Save to TZX output file
func (t *TZX) SaveTzx() error {
	if len(t.OutFilename) == 0 {
		return errors.New("empty output filename?")
	}
	if err := t.create(); err != nil {
		return err
	}
	defer t.OutFile.Close()

	if cfg.Main.Verbose {
		fmt.Println("> TZX:")
	}
	// header
	if cfg.Main.Verbose {
		fmt.Printf("  [%03d] Header\n", t.Written)
	}
	t.write([]byte("ZXTape!")) // [7b] signature
	t.write([]byte{0x1A})      // [1b] 1A/26 string end
	t.write([]byte{1, 20})     // [2b] 1,20 Major/minor revision number

	// ID 30 - Text Description
	if cfg.Main.Verbose {
		fmt.Printf("  [%03d] Block ID 0x30\n", t.Written)
	}
	signature := "Rip by RetroConvert " + build.VersionShort()
	t.write([]byte{0x30})
	t.write([]byte{byte(len(signature))})
	t.write([]byte(signature))

	// ID 32 - Archive info
	if cfg.Main.Verbose {
		fmt.Printf("  [%03d] Block ID 0x32\n", t.Written)
	}
	title := cfg.Main.Title
	if title == "" {
		title = filepath.Base(cfg.Main.InFile)
	}
	textBlock := []byte{0x32, 0, 0}
	textBlock = append(textBlock, 1)
	textBlock = append(textBlock, 0)
	textBlock = append(textBlock, byte(len(title)))
	for _, s := range title {
		textBlock = append(textBlock, byte(s))
	}
	textLen := []byte{0, 0}
	binary.LittleEndian.PutUint16(textLen, uint16(len(textBlock)-3))
	textBlock[1] = textLen[0]
	textBlock[2] = textLen[1]
	t.write(textBlock)

	// blocks (all ID 15)
	for _, block := range t.Blocks {
		// create the bitstream at first
		var bitstream []byte
		var pos uint
		var currByte byte
		for _, pulse := range block.Pulses {
			for i := uint(1); i <= pulse.Duration; i++ {
				pos++
				if pulse.Positive {
					currByte |= 1 << (8 - pos)
				}
				if pos == 8 {
					bitstream = append(bitstream, currByte)
					currByte = 0
					pos = 0
				}
			}
		}
		if pos > 0 {
			bitstream = append(bitstream, currByte)
		}

		// block ID
		t.write([]byte{0x15})
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] Block ID: 0x15\n", t.Written)
		}
		if cfg.Main.Verbose {
			fmt.Printf("  [INF] SAMPLES: %s BITSTREAM: %s x8: %s\n",
				humanize.Comma(int64(block.SampleCount)),
				humanize.Comma(int64(len(bitstream))),
				humanize.Comma(int64(len(bitstream)*8)),
			)
		}

		// [2b] T-States per bit
		tstates := []byte{0, 0}
		binary.LittleEndian.PutUint16(tstates, t.TStates)
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] TStates/Sample: %d LSB: %v\n", t.Written, t.TStates, tstates)
		}
		t.write(tstates)

		// [2b]Pause after this block in ms
		pause := []byte{0, 0}
		binary.LittleEndian.PutUint16(pause, uint16(block.Pause))
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] Pause after: %d LSB: %d MSB: %d\n",
				t.Written, block.Pause, pause[0], pause[1])
		}
		t.write(pause)

		// [1b] Used bits in the last byte
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] Used bits in the last byte: %d\n", t.Written, pos)
		}
		t.write([]byte{byte(pos)})

		// length 24bits LSB
		u := make([]byte, 4)
		binary.LittleEndian.PutUint32(u, uint32(len(bitstream)))
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] datalen: %d (0x%X) LSB:%v (0x%x%x%x)\n", t.Written,
				len(bitstream), len(bitstream), u[0:3], u[0], u[1], u[2])
		}
		t.write(u[0:3]) // [3b] Data len

		// bitstream data
		t.write(bitstream)
	}

	t.close()

	if cfg.Main.Verbose {
		fmt.Printf("  [INF] %s written\n", humanize.Bytes(uint64(t.Written)))
	}

	return nil
}
