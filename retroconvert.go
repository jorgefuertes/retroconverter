package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"git.martianoids.com/queru/retroconverter/internal/banner"
	"git.martianoids.com/queru/retroconverter/internal/build"
	"git.martianoids.com/queru/retroconverter/internal/cfg"
	"git.martianoids.com/queru/retroconverter/internal/progress"
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
	beginTime := time.Now()
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
	if cfg.Main.Verbose || cfg.Main.Progress {
		fmt.Println(banner.Title)
	}

	msg := make(chan string)
	defer close(msg)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	if cfg.Main.Progress && !cfg.Main.Verbose {
		go progress.Run(msg, wg, "Loading wav")
	} else {
		go progress.Fake(msg, wg)
	}

	t := new(tzx.TZX)
	t.OutFilename = cfg.Main.OutFile
	err := t.Load(cfg.Main.InFile)
	check(err)

	msg <- "Converting to pulses"
	t.ToPulses()

	if cfg.Main.ReSample > 0 {
		msg <- "Downsampling"
		if cfg.Main.Verbose {
			fmt.Printf("> Downsampling from %d to %d\n", t.Format.SampleRate, cfg.Main.ReSample)
		}
		t.DownSample(cfg.Main.ReSample)
	}

	if cfg.Main.Verbose {
		t.BlockStats()
	}

	msg <- fmt.Sprintf("Saving %s", cfg.Main.OutFile)
	err = t.SaveTzx()
	check(err)

	msg <- progress.END
	wg.Wait()
	if cfg.Main.Verbose || cfg.Main.Progress {
		fmt.Printf("> Converted in %.2f seconds\n", time.Since(beginTime).Seconds())
	}
}
