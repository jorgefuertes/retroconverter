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
	if cfg.Main.Verbose {
		fmt.Println(banner.Title)
	}

	w := new(qconvert.Wav)
	err := w.Load(cfg.Main.InFile)
	check(err)

	w.ToPulses()
	if cfg.Main.Verbose {
		w.BlockStats()
	}
	w.SaveTzx(cfg.Main.OutFile)
	if cfg.Main.Verbose {
		fmt.Println("> EOF")
	}
}
