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
	w.Filters.PHi = cfg.Main.PulseHi
	w.Filters.PLo = cfg.Main.PulseLo
	w.Filters.SHi = cfg.Main.SilenceHi
	w.Filters.SLo = cfg.Main.SilenceLo
	err := w.Load(cfg.Main.InFile)
	check(err)

	if cfg.Main.Normalize {
		w.CalcLevels()
		fmt.Printf("Pulse HI: %+04d LO: %+04d Fix Factor: %.3f\n", w.Highest, w.Lowest, w.Factor)
		w.FixLevels()
		fmt.Printf("Pulse HI: %+04d LO: %+04d (Normalized)\n", w.Highest, w.Lowest)
	}

	if err := w.CalcPulse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w.BlockStats()
	w.SaveTzx(cfg.Main.OutFile)
}
