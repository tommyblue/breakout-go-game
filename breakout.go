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
	screenW, screenH := screen.Size()
	var playerDirection float64
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		playerDirection = 1
		if g.p.x+g.p.w+g.p.speed <= float64(screenW) {
			g.p.x += g.p.speed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		playerDirection = -1
		g.p.x -= g.p.speed
		if g.p.x < 0 {
			g.p.x = 0
		}
	}
	for _, b := range g.balls {
		newX := b.x + b.direction[0]*b.speed

		if newX < 0 || newX+2*b.radius > float64(screenW) {
			b.direction[0] = -b.direction[0]
		}
		b.x += b.direction[0] * b.speed

		newY := b.y + b.direction[1]*b.speed
		collided := g.collision(b)
		if newY < 0 || collided {
			b.direction[1] = -b.direction[1]
			if collided {
				b.direction[0] += 0.2 * playerDirection
			}
		} else if newY+2*b.radius > float64(screenH) {
			g.initElements()
		}
		b.y += b.direction[1] * b.speed
	}
	return nil
}

func (g *Game) collision(b *ball) bool {
	if b.y+2*b.radius == g.p.y && (b.x <= g.p.x+g.p.w && b.x+2*b.radius >= g.p.x) {
		return true
	}
	return false
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
	g.initElements()
}

func (g *Game) initElements() {
	playerImg, err := ebiten.NewImage(100, 20, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	w, h := playerImg.Size()
	g.p = &player{
		w:     float64(w),
		h:     float64(h),
		x:     float64(W_WIDTH/2 - w/2),
		y:     float64(W_HEIGHT - h - 20),
		speed: 15,
		color: color.RGBA{0xff, 0x00, 0x00, 0xff},
		image: playerImg,
	}
	ballImg, err := ebiten.NewImage(20, 20, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	g.balls = []*ball{}
	g.balls = append(g.balls, &ball{
		radius:    10,
		x:         float64(W_WIDTH/2 - 10*2/2),
		y:         10,
		speed:     10,
		direction: [2]float64{0, 1},
		color:     color.RGBA{0xff, 0x00, 0xff, 0xff},
		image:     ballImg,
	})
}

func Start() {
	g := &Game{}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
