package cfg

// MainConfig - main config type
type MainConfig struct {
	Version   bool   `flag optional name:"version" short:"v" help:"Display version and exit." type:"flag"`
	Normalize bool   `flag optional name:"normalize" short:"n" help:"Normalization on/off (on). " type:"flag" default:"true"`
	Stats     bool   `flag optional name:"stats" short:"s" help:"Display statistics (off)." type:"flag" default:"false"`
	PulseHi   int    `flag name:"pulse-high" help:"Pulse filter high (64)." default:"64"`
	PulseLo   int    `flag name:"pulse-low" help:"Pusle filter low (-64)." default:"-64"`
	SilenceHi int    `flag name:"silence-high" help:"Silence filter high (10)." default:"10"`
	SilenceLo int    `flag name:"silence-low" help:"Silence filter low (-10." default:"-10"`
	InFile    string `arg optional name:"in" help:"Input .wav file." type:"path"`
	OutFile   string `arg optional name:"out" help:"Output .tzx file. Defaults to out.tzx." type:"path"`
}

// Main - Main configuration
var Main MainConfig
