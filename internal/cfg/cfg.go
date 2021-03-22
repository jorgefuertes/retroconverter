package cfg

// AppName - Applicaton name
const AppName = "RetroWiki Tape Converter"

// Author - Application author
const Author = "©2021 Jorge Fuertes & Ramón Martinez"

// MainConfig - main config type
type MainConfig struct {
	InFile  *string
	OutFile *string
}

// Main - Main configuration
var Main MainConfig
