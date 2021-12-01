package pendulum

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/project-auxo/auxo/olympus/logging"
)

const (
	thickness      = 3
	pendulumLength = 300
	radius         = 25
	cartWidth      = 150
	cartHeight     = 20
	movementAcc    = 200
	g              = 9.81
	frictionCoeff  = 0.992
	cartM          = 10 // Kg
	knobM          = 3  // Kg
)

var (
	bounds        = pixel.R(0, 0, 1024, 768)
	centerMassPos = pixel.V(512, 200)
	p             = new()
	log           = logging.Base()
)

type StateVec struct {
	cartPos, cartVel pixel.Vec
	penAngle, penVel float64
}

type Pendulum struct {
	line  pixel.Line
	knob  pixel.Circle
	cart  pixel.Rect
	state StateVec
}

func new() (p *Pendulum) {
	line := pixel.Line{
		A: centerMassPos,
		B: centerMassPos.Add(pixel.V(0, pendulumLength)),
	}
	knob := pixel.Circle{
		Center: line.B.Add(pixel.V(0, radius)),
		Radius: radius,
	}
	cart := pixel.Rect{
		Min: pixel.V(line.A.X-cartWidth, line.A.Y-cartHeight),
		Max: pixel.V(line.A.X+cartWidth, line.A.Y),
	}
	state := StateVec{
		cartPos:  centerMassPos,
		cartVel:  pixel.ZV,
		penAngle: 0,
		penVel:   0,
	}
	return &Pendulum{
		line:  line,
		knob:  knob,
		cart:  cart,
		state: state,
	}
}

func (p *Pendulum) draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = color.White

	// Drawing the pendulum stick
	imd.Push(p.line.A, p.line.B)
	imd.Line(thickness)

	// Drawing the knob at the top of the pendulum
	imd.Push(p.knob.Center)
	imd.Circle(p.knob.Radius, thickness)

	// Drawing the cart which is at the base of the pendulum
	imd.Push(p.cart.Min, p.cart.Max)
	imd.Rectangle(thickness)

	imd.Draw(win)
}

func (p *Pendulum) update(win *pixelgl.Window, dt float64) {
	// Should update the individual parts using time and the StateVec
	if win.Pressed(pixelgl.KeyA) || win.Pressed(pixelgl.KeyD) {
		appliedForceVec := pixel.V(cartM*movementAcc, 0)
		if win.Pressed(pixelgl.KeyA) {
			// Move cart left
			appliedForceVec = appliedForceVec.Scaled(-1)
		}
		accVec := appliedForceVec.Scaled(1.0 / cartM)
		p.state.cartVel = p.state.cartVel.Add(accVec.Scaled(dt))
	}
	// Friction
	p.state.cartVel = p.state.cartVel.Scaled(frictionCoeff)

	newMin := p.cart.Min.Add(p.state.cartVel.Scaled(dt))
	newMax := p.cart.Max.Add(p.state.cartVel.Scaled(dt))
	if bounds.Contains(newMin) && bounds.Contains(newMax) {
		p.cart.Min = newMin
		p.cart.Max = newMax
		p.line.A = p.line.A.Add(p.state.cartVel.Scaled(dt))
	}
}

func Run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Inverted Pendulum Simulation",
		Bounds: bounds,
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatalf("failed to make simulation window: %v", err)
	}

	last := time.Now()
	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			win.SetClosed(true)
			break
		}
		dt := time.Since(last).Seconds()
		last = time.Now()
		p.update(win, dt)

		win.Clear(colornames.Black)
		p.draw(win)
		win.Update()
	}
}
