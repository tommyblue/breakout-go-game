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
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.p.x+g.p.w+g.p.speed <= float64(screenW) {
			g.p.x += g.p.speed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
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
		collided, ratio := g.collision(b)
		if newY < 0 || collided {
			b.direction[1] = -b.direction[1]
			if collided {
				b.direction[0] += ratio
			}
		} else if newY+2*b.radius > float64(screenH) {
			g.initElements()
		}
		b.y += b.direction[1] * b.speed
	}
	return nil
}

/*

 .ooxiixoo.
 ----------
|          |
 ----------
*/
func (g *Game) collision(b *ball) (bool, float64) {
	if b.y+2*b.radius == g.p.y && (b.x <= g.p.x+g.p.w && b.x+2*b.radius >= g.p.x) {
		var ratio float64
		r := b.x + 2*b.radius
		l := b.x
		if l >= g.p.x && r <= g.p.x+g.p.w*10/100 {
			// fmt.Println("<=10%")
			ratio = -0.8
		} else if r > g.p.x+g.p.w*10/100 && r <= g.p.x+g.p.w*30/100 {
			// fmt.Println(">10% <=30%")
			ratio = -0.5
		} else if r > g.p.x+g.p.w*30/100 && r <= g.p.x+g.p.w*45/100 {
			// fmt.Println(">30% <= 45%")
			ratio = -0.2
		} else if l >= g.p.x+g.p.w*55/100 && l < g.p.x+g.p.w*70/100 {
			// fmt.Println(">=55% <70%")
			ratio = 0.2
		} else if l >= g.p.x+g.p.w*70/100 && l < g.p.x+g.p.w*90/100 {
			// fmt.Println(">=70% <90%")
			ratio = 0.5
		} else if l <= g.p.x+g.p.w && l >= g.p.x+g.p.w*90/100 {
			// fmt.Println(">=90%")
			ratio = 0.8
		}
		return true, ratio
	}
	return false, 0
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
		speed: 10,
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
