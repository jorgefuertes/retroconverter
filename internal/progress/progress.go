package progress

import (
	"fmt"
	"sync"
	"time"
)

const END = "#{CMD:END}"

// Fake progress msg channel listener
func Fake(messages chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg := <-messages
		if msg == END {
			return
		}
	}
}

// Run the progress indicator
func Run(messages chan string, wg *sync.WaitGroup, initMsg string) {
	var lastMsg string
	var count int
	bar := "\\|/-"
	defer wg.Done()

	fmt.Printf("> %s....", initMsg)
	lastMsg = initMsg

	for {
		select {
		case msg := <-messages:
			if msg == END {
				fmt.Println("\bOK")
				return
			}
			if msg != lastMsg {
				lastMsg = msg
				fmt.Println("\bOK")
				fmt.Printf("> %s....", msg)
			}
		default:
			fmt.Printf("\b%c", bar[count])
			count++
			if count >= len(bar) {
				count = 0
			}
			time.Sleep(30 * time.Millisecond)
		}
	}
}
