package main

import (
	"fmt"
	"time"
)

func CLI_main(Is_GPU bool) {
	benchmark_running = true
	chan_data = make(chan chan_t, 4096)

	go benchmark(Is_GPU)
	var temp_chan_data chan_t
	for benchmark_running {
	L2:
		for {
			select {
			case temp_chan_data = <-chan_data:

			default:
				break L2
			}
		}
		fmt.Printf("\r Score:%f Area:%f", temp_chan_data.score, temp_chan_data.S)
		time.Sleep(time.Millisecond * 250)
	}
}
