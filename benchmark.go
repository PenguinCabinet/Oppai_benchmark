package main

import (
	"runtime"
	"sync"
	"time"
)

var thread_n int

func integral_f_p(alpha, beta float64, f func(float64) float64) float64 {
	wg := &sync.WaitGroup{}

	N := 1000000

	delta := (beta - alpha) / float64(N)

	Data := make([]float64, N+1)

	for i := 0; i < thread_n; i++ {
		wg.Add(1)
		go func(i2 int) {
			temp_s := (beta - alpha) / float64(thread_n) * float64(i2)
			temp_e := (beta - alpha) / float64(thread_n) * float64(i2+1)
			j := 0
			for j_d := temp_s; j_d < temp_e; j_d += delta {
				f_a := f(j_d)
				f_b := f(j_d + delta)
				f_m := f((j_d + j_d + delta) / 2)

				Data[j+i2*(N/thread_n)] = delta / 6 * (f_a + 4*f_m + f_b)
				j++
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	A := 0.0
	for _, e := range Data {
		A += e
	}
	return A
}

func Get_score(v float64) float64 {
	return 1 / (v / 1000000.0) * 1000000
}

func benchmark() {
	thread_n = runtime.NumCPU()
	N := 32.0
	N_sec := 30
	delta_time := 0.5
	chan_data <- chan_t{
		t:     0,
		S:     -1,
		score: -1,
	}
	scores := []float64{}
	start_time := time.Now()
	for t := 0.0; time.Now().Unix() <= start_time.Add(time.Duration(N_sec)*time.Second).Unix() || t < N; t += delta_time {
		t2 := t
		start_time := time.Now()
		S := integral_f_p(-1000, 1000, func(v float64) float64 { return Oppai_func(v, t2) })
		end_time := time.Now()
		scores = append(scores, Get_score(float64(end_time.Sub(start_time))))
		mean_score := 0.0
		for i := 0; i < len(scores); i++ {
			mean_score += scores[i]
		}
		mean_score /= float64(len(scores))
		chan_data <- chan_t{
			t:     t,
			S:     S,
			score: mean_score,
		}
	}
	benchmark_running = false
}
