package main

import "fmt"

func CLI_main() {
	benchmark_running = true
	chan_data = make(chan chan_t, 4096)

	go benchmark()
	var temp_chan_data chan_t
	for benchmark_running {
		select {
		case temp_chan_data = <-chan_data:
			fmt.Printf("\r Score:%f Area:%f", temp_chan_data.score, temp_chan_data.S)
		default:

		}
	}
}
