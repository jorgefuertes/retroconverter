package audio

type Block struct {
	TStates uint16 // T-States per sample
	Pause   uint16 // Pause after this block
}

type TZX struct {
	Blocks []Block
}
