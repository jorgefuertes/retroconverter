package main

import (
	"fmt"

	"git.martianoids.com/queru/retroconverter/internal/build"
	"git.martianoids.com/queru/retroconverter/internal/cfg"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	// command line flags and params
	cfg.Main.InFile = kingpin.Flag("in", "Input file (wav)").Short('i').Required().String()
	cfg.Main.OutFile = kingpin.Flag("out", "Output file: Defaults to input_file_name.tzx").
		Short('o').String()
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(build.VersionShort()).
		Author(cfg.Author)
	kingpin.CommandLine.Help = "RetroWiki Tape Converter"
	kingpin.Parse()

	fmt.Println("EOF")
}
