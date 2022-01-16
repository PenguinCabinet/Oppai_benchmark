package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func Oppai_func(y float64, t float64) float64 {
	y = 0.02 * (y - 100)

	a1 := (1.5 * math.Exp((0.12*math.Sin(t)-0.5)*math.Pow((y+0.16*math.Sin(t)), 2))) / (1 + math.Exp(-20*(5*y+math.Sin(t))))
	a2 := ((1.5 + 0.8*math.Pow((y+0.2*math.Sin(t)), 3)) * math.Pow(1+math.Exp(20*(5*y+math.Sin(t))), -1)) / (1 + math.Exp(-(100*(y+1) + 16*math.Sin(t))))
	a3 := (0.2 * (math.Exp(-math.Pow(y+1, 2)) + 1)) / (1 + math.Exp(100*(y+1)+16*math.Sin(t)))
	a4 := 0.1 / math.Exp(2*math.Pow((10*y+1.2*(2+math.Sin(t))*math.Sin(t)), 4))

	return 65 * (a1 + a2 + a3 + a4)

}

type chan_t struct {
	t     float64
	S     float64
	score float64
}

var chan_data chan chan_t

var benchmark_running bool

type Game struct {
	x              []float64
	y              []float64
	temp_chan_data chan_t
}

var font_face font.Face

func (g *Game) Init() {
	g.y = make([]float64, 600)
	g.x = make([]float64, 600)
	for i := 0; i < len(g.y); i += 1 {
		delta := float64(i)
		//fmt.Printf("%f\n", g.t)
		//fmt.Printf("%f\n", delta)
		g.y[i] = delta
		g.x[i] = Oppai_func(g.y[i], g.temp_chan_data.t)
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if !benchmark_running {
			benchmark_running = true
			go benchmark()
		}
	}
L1:
	for {
		select {
		case g.temp_chan_data = <-chan_data:

		default:
			break L1
		}
	}
	g.Init()
	screen.Fill(color.RGBA{255, 255, 255, 0xff})
	if !benchmark_running {
		//fmt.Println("TRUE")
		text.Draw(screen, "Sキーでベンチマークスタート", font_face, 120, 50, color.Black)
		text.Draw(screen, "Putting the key of S,\nstart a benchmark ", font_face, 120, 70, color.Black)
	}
	text.Draw(screen, "おっぱい関数積分ベンチマーク", font_face, 120, 20, color.Black)
	text.Draw(screen, fmt.Sprintf("Score:%f", g.temp_chan_data.score), font_face, 120, 150, color.Black)
	text.Draw(screen, fmt.Sprintf("面積(Area):%f", g.temp_chan_data.S), font_face, 120, 200, color.Black)

	for i := 5; i < len(g.y); i += 5 {
		//fmt.Printf("%f,%f\n", g.x[i], g.y[i])
		//ebitenutil.DebugPrint(screen, "Hello, World!")
		ebitenutil.DrawLine(screen, g.x[i-5], g.y[i-5], g.x[i], g.y[i], color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func GUI_main() {
	benchmark_running = false
	temp, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	font_face, _ = opentype.NewFace(temp, &opentype.FaceOptions{
		Size:    8,
		DPI:     128,
		Hinting: font.HintingFull,
	})

	chan_data = make(chan chan_t, 4096)

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("OPPAI Benchmark")
	G := &Game{}
	G.Init()
	//go benchmark()
	if err := ebiten.RunGame(G); err != nil {
		log.Fatal(err)
	}
}
