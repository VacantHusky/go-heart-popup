package main

import (
	_ "embed"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed simkai.ttf
var notoSansSC []byte
var windowWidth, windowHeight float64

type Game struct {
	bg        *ebiten.Image
	particles ParticleManager
	popups    PopupManager
	nextSpawn time.Time
}

func NewGame() *Game {
	g := &Game{}
	g.nextSpawn = time.Now().Add(500 * time.Millisecond)
	g.popups.InitPopup()
	g.popups.InitTexts()
	return g
}

func (g *Game) Update() error {
	// ESC 退出
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// 鼠标左键
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.particles.New(float64(x), float64(y))
	}

	if time.Now().After(g.nextSpawn) {
		x := rand.Float64() * windowWidth
		y := rand.Float64()*windowHeight/2 + windowHeight/4
		g.particles.New(float64(x), float64(y))
		g.nextSpawn = time.Now().Add(time.Duration(200+rand.Intn(600)) * time.Millisecond)
	}

	fps := ebiten.ActualTPS()
	dt := 1.0 / 60.0
	if fps > 0 {
		dt = 1.0 / fps
	}

	g.particles.Update(dt)
	g.popups.Update(1.0 / 60.0)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.bg != nil {
		sw, sh := g.bg.Bounds().Dx(), g.bg.Bounds().Dy()
		op := &ebiten.DrawImageOptions{}
		scaleX := windowWidth / float64(sw)
		scaleY := windowHeight / float64(sh)
		op.GeoM.Scale(scaleX, scaleY)
		op.ColorScale.Scale(0.6, 0.6, 0.6, 1)
		screen.DrawImage(g.bg, op)
	} else {
		screen.Fill(color.Black)
	}

	g.popups.Draw(screen)
	g.particles.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(windowWidth), int(windowHeight)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := NewGame()
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("WH LOVE SJH")

	// time.Sleep(time.Millisecond * 100)
	img, w, h, err := captureFullScreenImage()
	if err != nil {
		log.Fatalf("无法截取屏幕: %v", err)
	}
	windowWidth, windowHeight = float64(w), float64(h)
	game.bg = ebiten.NewImageFromImage(img)
	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
