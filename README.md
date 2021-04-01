# RetroConverter

Wav to TZX (ID 15, Direct recording) converter. It works with wav files: 11025/22050/44100 Khz.

## Author

©2021 Jorge Fuertes & Ramón Martinez

## License

This software is licensed under GPL3.

## Usage
```
Usage: retroconvert [<in>] [<out>]

Arguments:
  [<in>]     Input .wav file.
  [<out>]    Output .tzx file. Defaults to out.tzx.

Flags:
  -h, --help            Show context-sensitive help.
  -v, --version         Display version and exit.
  -b, --verbose         Verbosity (off).
  -r, --resample=0      Down sample to 22050 or 11025 from greater multiple freq.
  -t, --title=STRING    Archite title. Defaults to out filename.
```

## Download and installation

Go to the [latest release](https://github.com/jorgefuertes/retroconverter/releases/latest) page.
