package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 800
	screenScale  = 5
	baseSize     = 3
	G            = 0.000001
)

type Body struct {
	Mass     float64
	Position [2]float64
	Velocity [2]float64
	Color    color.Color
}

func main() {
	// Periodic orbits:
	// https://arxiv.org/pdf/1805.07980.pdf
	bodies := []*Body{
		{
			Mass:     0.6,
			Position: [2]float64{-0.5, 0},
			Velocity: [2]float64{0, 0},
			Color:    color.RGBA{255, 0, 0, 0},
		},
		{
			Mass:     0.8,
			Position: [2]float64{0.5, 0},
			Velocity: [2]float64{0, 0},
			Color:    color.RGBA{0, 255, 0, 0},
		},
		{
			Mass:     1,
			Position: [2]float64{0.289925259, 0.4030770616},
			Velocity: [2]float64{0, 0},
			Color:    color.RGBA{0, 0, 255, 0},
		},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("The 3 Body Problem: a Simulation")

	fmt.Println("Starting simulation...")
	if err := ebiten.RunGame(&Game{bodies: bodies}); err != nil {
		fmt.Println(err)
	}
}

type Game struct {
	bodies []*Body
}

func (g *Game) Update() error {
	// Check for exit
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return fmt.Errorf("Simulation ended.")
	}

	// Calculate forces
	for i := range g.bodies {
		for j := range g.bodies {
			if i != j {
				dx := g.bodies[j].Position[0] - g.bodies[i].Position[0]
				dy := g.bodies[j].Position[1] - g.bodies[i].Position[1]
				distSq := dx*dx + dy*dy
				// dist := math.Sqrt(distSq)
				force := G * g.bodies[i].Mass * g.bodies[j].Mass / distSq
				angle := math.Atan2(dy, dx)
				g.bodies[i].Velocity[0] += force * math.Cos(angle) / g.bodies[i].Mass
				g.bodies[i].Velocity[1] += force * math.Sin(angle) / g.bodies[i].Mass
			}
		}
	}

	// Update positions
	for _, body := range g.bodies {
		body.Position[0] += body.Velocity[0]
		body.Position[1] += body.Velocity[1]
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	centerX := float64(screenWidth / 2)
	centerY := float64(screenHeight / 2)

	for _, body := range g.bodies {
		// Map the coordinates to screen coordinates
		screenX := centerX + body.Position[0]*screenScale
		screenY := centerY + body.Position[1]*screenScale

		vector.DrawFilledCircle(screen, float32(screenX), float32(screenY),
			float32(baseSize*body.Mass), body.Color, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
