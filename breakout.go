package breakout

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	W_WIDTH  = 800
	W_HEIGHT = 600
)

type player struct {
	w     float64
	h     float64
	speed float64
	x     float64
	y     float64
	color color.RGBA
	image *ebiten.Image
}

type Game struct {
	p *player
}

func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s_w, _ := screen.Size()
		if g.p.x+g.p.w+g.p.speed <= float64(s_w) {
			g.p.x += g.p.speed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.p.x -= g.p.speed
		if g.p.x < 0 {
			g.p.x = 0
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	red := float64(g.p.color.R) / 0xff
	green := float64(g.p.color.G) / 0xff
	blue := float64(g.p.color.B) / 0xff
	op.ColorM.Translate(red, green, blue, 1)

	op.GeoM.Translate(g.p.x, g.p.y)

	if err := screen.DrawImage(g.p.image, op); err != nil {
		log.Fatal(err)
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return W_WIDTH, W_HEIGHT
}

func (g *Game) Init() {
	ebiten.SetWindowSize(W_WIDTH, W_HEIGHT)
	ebiten.SetWindowTitle("Breakout")
	i, err := ebiten.NewImage(100, 20, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	w, h := i.Size()
	g.p = &player{
		w:     float64(w),
		h:     float64(h),
		x:     float64(W_WIDTH/2 - w/2),
		y:     float64(W_HEIGHT - h - 20),
		speed: 15,
		color: color.RGBA{0xff, 0x00, 0x00, 0xff},
		image: i,
	}
}

func Start() {
	g := &Game{}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
