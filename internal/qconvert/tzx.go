package qconvert

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"git.martianoids.com/queru/retroconverter/internal/build"
	"git.martianoids.com/queru/retroconverter/internal/cfg"
	"github.com/dustin/go-humanize"
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

	if cfg.Main.Verbose {
		fmt.Println("> TZX:")
	}
	// header
	if cfg.Main.Verbose {
		fmt.Printf("  [%03d] Header\n", tzx.Written)
	}
	tzx.write([]byte("ZXTape!")) // [7b] signature
	tzx.write([]byte{0x1A})      // [1b] 1A/26 string end
	tzx.write([]byte{1, 20})     // [2b] 1,20 Major/minor revision number

	// ID 30 - Text Description
	if cfg.Main.Verbose {
		fmt.Printf("  [%03d] Block ID 0x30\n", tzx.Written)
	}
	signature := "Rip by RetroConvert " + build.VersionShort()
	tzx.write([]byte{0x30})
	tzx.write([]byte{byte(len(signature))})
	tzx.write([]byte(signature))

	// ID 32 - Archive info
	if cfg.Main.Verbose {
		fmt.Printf("  [%03d] Block ID 0x32\n", tzx.Written)
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
	tzx.write(textBlock)
	fmt.Printf("  TEXTBLOCK %d:\n", uint16(len(textBlock)-3))
	for _, v := range textBlock {
		if v >= 32 {
			fmt.Printf("    0x%02x %03d [%s]\n", v, v, string(v))
			continue
		}
		fmt.Printf("    0x%02x %03d [.]\n", v, v)
	}
	fmt.Println()

	// blocks (all ID 15)
	for _, block := range w.Blocks {
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
		tzx.write([]byte{0x15})
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] Block ID: 0x15\n", tzx.Written)
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
		binary.LittleEndian.PutUint16(tstates, w.TStates)
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] TStates/Sample: %d LSB: %v\n", tzx.Written, w.TStates, tstates)
		}
		tzx.write(tstates)

		// [2b]Pause after this block in ms
		pause := []byte{0, 0}
		binary.LittleEndian.PutUint16(pause, uint16(block.Pause))
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] Pause after: %d LSB: %d MSB: %d\n",
				tzx.Written, block.Pause, pause[0], pause[1])
		}
		tzx.write(pause)

		// [1b] Used bits in the last byte
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] Used bits in the last byte: %d\n", tzx.Written, pos)
		}
		tzx.write([]byte{byte(pos)})

		// length 24bits LSB
		u := make([]byte, 4)
		binary.LittleEndian.PutUint32(u, uint32(len(bitstream)))
		if cfg.Main.Verbose {
			fmt.Printf("  [%03d] datalen: %d (0x%X) LSB:%v (0x%x%x%x)\n", tzx.Written,
				len(bitstream), len(bitstream), u[0:3], u[0], u[1], u[2])
		}
		tzx.write(u[0:3]) // [3b] Data len

		// bitstream data
		tzx.write(bitstream)
	}

	// block ID 20 - Stop the tape
	if cfg.Main.Verbose {
		fmt.Printf("  [%03d] ID 0x20 Stop tape\n", tzx.Written)
	}
	tzx.write([]byte{0x20, 0, 0})

	if cfg.Main.Verbose {
		fmt.Printf("  [INF] %s written\n", humanize.Bytes(uint64(tzx.Written)))
	}

	return nil
}
