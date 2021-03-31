package cfg

// MainConfig - main config type
type MainConfig struct {
	Version   bool   `flag optional name:"version" short:"v" help:"Display version and exit." type:"flag"`
	Normalize bool   `flag optional name:"normalize" short:"n" help:"Normalization on/off (on). " type:"flag" default:"true"`
	Verbose   bool   `flag optional name:"verbose" short:"b" help:"Verbosity (off)." type:"flag" default:"false"`
	InFile    string `arg optional name:"in" help:"Input .wav file." type:"path"`
	ReSample  int    `flag optional name:"resample" short:"r" help:"Down sample to 22050 or 11025 from greather multiple freq."`
	OutFile   string `arg optional name:"out" help:"Output .tzx file. Defaults to out.tzx." type:"path"`
	Title     string `arg optional name:"title" help:"Archite title. Defaults to out filename."`
}

// Main - Main configuration
var Main MainConfig
