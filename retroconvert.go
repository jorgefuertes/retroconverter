package main

import (
	"errors"
	"fmt"
	"os"

	"git.martianoids.com/queru/retroconverter/internal/banner"
	"git.martianoids.com/queru/retroconverter/internal/build"
	"git.martianoids.com/queru/retroconverter/internal/cfg"
	"git.martianoids.com/queru/retroconverter/internal/qconvert"
	"github.com/alecthomas/kong"
)

func check(e error) {
	if e == nil {
		return
	}
	fmt.Println("ERROR:", e)
	fmt.Println("Please try --help")
	os.Exit(1)
}

func main() {
	// command line flags and params
	kong.Parse(&cfg.Main)

	// version
	if cfg.Main.Version {
		fmt.Println(build.Version())
		os.Exit(0)
	}

	// input
	if len(cfg.Main.InFile) == 0 {
		check(errors.New("no input file?"))
	}

	// output
	if len(cfg.Main.OutFile) == 0 {
		cfg.Main.OutFile = "out.tzx"
	}

	// banner
	fmt.Println(banner.Title)

	w := new(qconvert.Wav)
	err := w.Load(cfg.Main.InFile)
	check(err)

	w.Stats()

	// out, err := os.Create("work/out.wav")
	// check(err)
	// e := wav.NewEncoder(out,
	// 	int(w.Decoder.SampleRate),
	// 	int(w.Decoder.BitDepth),
	// 	int(w.Decoder.NumChans),
	// 	int(w.Decoder.WavAudioFormat),
	// )

	// for {
	// 	chunk, err := w.Decoder.NextChunk()
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Printf("CHUNK %v: ", chunk.Pos)
	// 	count := 0
	// 	for {
	// 		_, err := chunk.ReadByte()
	// 		if err != nil {
	// 			fmt.Printf("%d bytes written\n", count)
	// 			break
	// 		}
	// 		count++
	// 	}
	// }
}
