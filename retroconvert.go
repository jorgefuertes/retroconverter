package main

import (
	"errors"
	"fmt"
	"os"

	"git.martianoids.com/queru/retroconverter/internal/banner"
	"git.martianoids.com/queru/retroconverter/internal/build"
	"git.martianoids.com/queru/retroconverter/internal/cfg"
	"git.martianoids.com/queru/retroconverter/internal/tzx"
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

	t := new(tzx.TZX)
	t.OutFilename = cfg.Main.OutFile
	err := t.Load(cfg.Main.InFile)
	check(err)

	t.ToPulses()

	if cfg.Main.ReSample > 0 {
		if cfg.Main.Verbose {
			fmt.Printf("> Downsampling from %d to %d\n", t.Format.SampleRate, cfg.Main.ReSample)
		}
		t.DownSample(cfg.Main.ReSample)
	}

	if cfg.Main.Verbose {
		t.BlockStats()
	}

	err = t.SaveTzx()
	check(err)

	if cfg.Main.Verbose {
		fmt.Println("> EOF")
	}
}
