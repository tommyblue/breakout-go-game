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
	w, h  float64
	speed float64
	x, y  float64
	color color.RGBA
	image *ebiten.Image
}

type ball struct {
	radius    float64
	x, y      float64
	speed     float64
	direction [2]float64
	color     color.RGBA
	image     *ebiten.Image
}

type Game struct {
	p     *player
	balls []*ball
}

func (g *Game) Update(screen *ebiten.Image) error {
	s_w, s_h := screen.Size()
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.p.x+g.p.w+g.p.speed <= float64(s_w) {
			g.p.x += g.p.speed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.p.x -= g.p.speed
		if g.p.x < 0 {
			g.p.x = 0
		}
	}
	for _, b := range g.balls {
		new_x := b.x + b.direction[0]*b.speed

		if new_x < 0 || new_x+2*b.radius > float64(s_w) {
			b.direction[0] = -b.direction[0]
		}
		b.x += b.direction[0] * b.speed

		new_y := b.y + b.direction[1]*b.speed
		if new_y < 0 {
			b.direction[1] = -b.direction[1]
		} else if new_y+2*b.radius > float64(s_h) {
			// TODO: this is the lose condition

			// inverting the direction for debugging purposes
			b.direction[1] = -b.direction[1]
		}
		b.y += b.direction[1] * b.speed
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawPlayer(screen)
	g.drawBalls(screen)
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
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

func (g *Game) drawBalls(screen *ebiten.Image) {
	for _, b := range g.balls {
		g.drawBall(screen, b)
	}
}

func (g *Game) drawBall(screen *ebiten.Image, b *ball) {
	op := &ebiten.DrawImageOptions{}

	red := float64(b.color.R) / 0xff
	green := float64(b.color.G) / 0xff
	blue := float64(b.color.B) / 0xff
	op.ColorM.Translate(red, green, blue, 1)

	op.GeoM.Translate(b.x, b.y)

	if err := screen.DrawImage(b.image, op); err != nil {
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
	b_i, err := ebiten.NewImage(20, 20, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	g.balls = append(g.balls, &ball{
		radius:    10,
		x:         float64(W_WIDTH/2 - 10*2/2),
		y:         10,
		speed:     10,
		direction: [2]float64{0, 1},
		color:     color.RGBA{0xff, 0x00, 0xff, 0xff},
		image:     b_i,
	})
}

func Start() {
	g := &Game{}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
