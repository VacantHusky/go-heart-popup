package main

import (
	"container/list"
	"image/color"
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Particle struct {
	X, Y   float64
	Vx, Vy float64
	Life   float64
	Size   float32
	Color  color.NRGBA
}

type ParticleManager struct {
	Particles list.List
	mu        sync.Mutex
}

func (pm *ParticleManager) New(x, y float64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	count := 220 + rand.Intn(400)
	baseHue := rand.Float64() * 360
	for range count {
		theta := rand.Float64() * 2 * math.Pi
		size := 1 + rand.Float32()*4
		speed := 1 + rand.Float64()*1.5*float64(count)/float64(size)
		p := &Particle{
			X:    x,
			Y:    y,
			Vx:   math.Cos(theta) * speed,
			Vy:   math.Sin(theta) * speed,
			Life: 3 + rand.Float64()*3,
			Size: size,
		}
		p.Color.R, p.Color.G, p.Color.B = hsvToRgb(baseHue+rand.Float64()*40-20, 0.9, 1.0)
		p.Color.A = 200
		pm.Particles.PushFront(p)
	}
}

func (pm *ParticleManager) Update(dt float64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for element := pm.Particles.Front(); element != nil; {
		p := element.Value.(*Particle)
		p.X += p.Vx * dt
		p.Y += p.Vy * dt
		p.Vy += 30 * dt
		p.Life -= dt
		p.Vx *= 0.998
		p.Vy *= 0.998
		p.Size *= 0.995

		if p.Life <= 0 || p.X < 0 || p.X > windowWidth || p.Y < 0 || p.Y > windowHeight {
			next := element.Next()
			pm.Particles.Remove(element)
			element = next
		} else {
			element = element.Next()
		}
	}
}

func (pm *ParticleManager) Draw(screen *ebiten.Image) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for element := pm.Particles.Front(); element != nil; element = element.Next() {
		p := element.Value.(*Particle)
		vector.FillCircle(screen, float32(p.X), float32(p.Y), p.Size, p.Color, false)
	}
}
