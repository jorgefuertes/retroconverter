package cfg

// MainConfig - main config type
type MainConfig struct {
	Version  bool   `flag optional name:"version" short:"v" help:"Display version and exit." type:"flag"`
	Verbose  bool   `flag optional name:"verbose" short:"b" help:"Verbosity (off)." type:"flag" default:"false"`
	Progress bool   `flag optional name:"progress" short:"p" help:"Show progress (true)." type:"flag" default:"true"`
	Title    string `flag optional name:"title" short:"t" help:"Archite title. Defaults to out filename."`
	ReSample uint32 `flag optional name:"resample" short:"r" help:"Down sample to 22050 or 11025 from greater multiple freq." default:0`
	InFile   string `arg optional name:"in" help:"Input .wav file." type:"path"`
	OutFile  string `arg optional name:"out" help:"Output .tzx file. Defaults to out.tzx." type:"path"`
}

// Main - Main configuration
var Main MainConfig
