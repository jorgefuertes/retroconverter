package cfg

// MainConfig - main config type
type MainConfig struct {
	Debug   bool   `flag optional name:"debug" help:"Debug level verbosity on." type:"flag"`
	Version bool   `flag optional name:"v" help:"Display version and exit." type:"flag"`
	InFile  string `arg optional name:"in" help:"Input .wav file." type:"path"`
	OutFile string `arg optional name:"out" help:"Output .tzx file. Defaults to out.tzx." type:"path"`
}

// Main - Main configuration
var Main MainConfig
