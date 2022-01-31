package main

import (
	"math"
	"runtime"
	"sync"
	"time"

	"gorgonia.org/gorgonia"
	. "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

func Oppai_func(y float64, t float64) float64 {
	y = 0.02 * (y - 100)

	a1 := (1.5 * math.Exp((0.12*math.Sin(t)-0.5)*math.Pow((y+0.16*math.Sin(t)), 2))) / (1 + math.Exp(-20*(5*y+math.Sin(t))))
	a2 := ((1.5 + 0.8*math.Pow((y+0.2*math.Sin(t)), 3)) * math.Pow(1+math.Exp(20*(5*y+math.Sin(t))), -1)) / (1 + math.Exp(-(100*(y+1) + 16*math.Sin(t))))
	a3 := (0.2 * (math.Exp(-math.Pow(y+1, 2)) + 1)) / (1 + math.Exp(100*(y+1)+16*math.Sin(t)))
	a4 := 0.1 / math.Exp(2*math.Pow((10*y+1.2*(2+math.Sin(t))*math.Sin(t)), 4))

	return 65 * (a1 + a2 + a3 + a4)

}

var Float tensor.Dtype

func Make_scalar_gpu(g *ExprGraph, v float64) *gorgonia.Node {
	//return gorgonia.NewScalar(g, gorgonia.Float64, gorgonia.WithName("y_temp1"), gorgonia.WithValue(v))
	return gorgonia.NodeFromAny(g, v, gorgonia.WithName("x"))

}

func Mul_Must(v1 *Node, v2 *Node) *Node {
	return Must(Mul(v1, v2))
}

func Div_Must(v1 *Node, v2 *Node) *Node {
	return Must(Div(v1, v2))
}

func Add_Must(v1 *Node, v2 *Node) *Node {
	return Must(Add(v1, v2))
}

func Sin_Must(v1 *Node) *Node {
	return Must(Sin(v1))
}

func Exp_Must(v1 *Node) *Node {
	return Must(Exp(v1))
}

func Pow_Must(v1 *Node, v2 *Node) *Node {
	return Must(Pow(v1, v2))
}

func Oppai_func_GPU(y2 float64, t2 float64) float64 {

	g := gorgonia.NewGraph()

	my_Make_scalar_gpu := func(v float64) *gorgonia.Node {
		return Make_scalar_gpu(g, v)
	}

	y := my_Make_scalar_gpu(y2)
	t := my_Make_scalar_gpu(t2)

	//y = gorgonia.NewScalar(g, Float, gorgonia.WithName("y"), gorgonia.WithValue(0.02*(y-100)))
	y_temp1 := my_Make_scalar_gpu(0.01)
	y_temp2 := my_Make_scalar_gpu(100)
	y = gorgonia.Must(gorgonia.Mul(gorgonia.Must(gorgonia.Sub(y, y_temp2)), y_temp1))
	t_sin_temp, _ := gorgonia.Sin(t)

	a1 := Div_Must(Mul_Must(my_Make_scalar_gpu(1.5), Exp_Must(Mul_Must(Add_Must(Mul_Must(my_Make_scalar_gpu(0.12), t_sin_temp), my_Make_scalar_gpu(-0.5)), Pow_Must(Mul_Must(Add_Must(y, my_Make_scalar_gpu(0.16)), t_sin_temp), my_Make_scalar_gpu(2.0))))), Add_Must(my_Make_scalar_gpu(1), Exp_Must(Mul_Must(my_Make_scalar_gpu(-20), Add_Must(Mul_Must(my_Make_scalar_gpu(5), y), Sin_Must(t))))))
	a2 := Mul_Must(Add_Must(my_Make_scalar_gpu(1.5), Mul_Must(my_Make_scalar_gpu(0.8), Pow_Must(Add_Must(y, Mul_Must(my_Make_scalar_gpu(0.2), t_sin_temp)), my_Make_scalar_gpu(3)))), Div_Must(Pow_Must(Add_Must(my_Make_scalar_gpu(1), Exp_Must(Mul_Must(my_Make_scalar_gpu(20), Add_Must(Mul_Must(my_Make_scalar_gpu(5), y), t_sin_temp)))), my_Make_scalar_gpu(-1)), Add_Must(my_Make_scalar_gpu(1), Exp_Must(Mul_Must(my_Make_scalar_gpu(-1), Add_Must(Mul_Must(my_Make_scalar_gpu(100), Add_Must(y, my_Make_scalar_gpu(1))), Mul_Must(my_Make_scalar_gpu(16), t_sin_temp)))))))
	a3 := Div_Must(Mul_Must(my_Make_scalar_gpu(0.2), Add_Must(Exp_Must(Mul_Must(my_Make_scalar_gpu(-1), Pow_Must(Add_Must(y, my_Make_scalar_gpu(1)), my_Make_scalar_gpu(2)))), my_Make_scalar_gpu(1))), Add_Must(my_Make_scalar_gpu(1), Exp_Must(Add_Must(Mul_Must(my_Make_scalar_gpu(100), Add_Must(y, my_Make_scalar_gpu(1))), Mul_Must(my_Make_scalar_gpu(16), t_sin_temp)))))
	a4 := Div_Must(my_Make_scalar_gpu(0.1), Exp_Must(Mul_Must(my_Make_scalar_gpu(2), Pow_Must(Add_Must(Mul_Must(my_Make_scalar_gpu(10), y), Mul_Must(my_Make_scalar_gpu(1.2), Mul_Must(Add_Must(my_Make_scalar_gpu(2), t_sin_temp), t_sin_temp))), my_Make_scalar_gpu(4)))))

	A := Add_Must(Add_Must(a1, a2), Add_Must(a3, a4))
	vm := gorgonia.NewTapeMachine(g)
	vm.RunAll()

	return float64(A.Value().Data().(float64)) //* (a1 + a2 + a3 + a4)

}

var thread_n int

func integral_f_p(alpha, beta float64, f func(float64) float64, Is_GPU bool) float64 {
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

func benchmark(Is_GPU bool) {
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
		S := 0.0
		if Is_GPU {
			S = integral_f_p(-1000, 1000, func(v float64) float64 { return Oppai_func(v, t2) }, Is_GPU)
		} else {
			S = integral_f_p(-1000, 1000, func(v float64) float64 { return Oppai_func_GPU(v, t2) }, Is_GPU)
		}
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
