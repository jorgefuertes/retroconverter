package qconvert

import (
	"os"
)

type Tzx struct {
	Filename string
	F        *os.File
	Written  int
}
