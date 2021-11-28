package seek

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
	thickness     = 3
	cartWidth     = 150
	cartHeight    = 20
	radius        = 25
	movementAcc   = 200
	frictionCoeff = 0.99
	cartM         = 5 // Kg
)

var (
	bounds            = pixel.R(0, 0, 1024, 768)
	initialCartPos    = pixel.V(200, 200)
	initialTargetPos  = pixel.V(900, initialCartPos.Y)
	otherTargetPos    = pixel.V(200, initialCartPos.Y)
	useOtherTargetPos = false
	cart              = new()
	log               = logging.Base()
)

type StateVec struct {
	cartPos pixel.Vec // Center of the cart
	cartVel pixel.Vec // Velocity of the cart
}

type SeekSim struct {
	cart  pixel.Rect
	goal  pixel.Circle
	state StateVec
}

func new() *SeekSim {
	cart := pixel.Rect{
		Min: pixel.V(initialCartPos.X-cartWidth, initialCartPos.Y-cartHeight),
		Max: pixel.V(initialCartPos.X+cartWidth, initialCartPos.Y),
	}
	goal := pixel.Circle{
		Center: initialTargetPos,
		Radius: radius,
	}
	state := StateVec{
		cartPos: initialCartPos,
		cartVel: pixel.ZV,
	}
	return &SeekSim{
		cart:  cart,
		state: state,
		goal:  goal,
	}
}

func (s *SeekSim) draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	// Drawing the cart
	imd.Color = color.White
	imd.Push(s.cart.Min, s.cart.Max)
	imd.Rectangle(thickness)

	// Drawing the target to seek
	imd.Color = colornames.Gold
	imd.Push(s.goal.Center)
	imd.Circle(radius, thickness)

	imd.Draw(win)
}

func (s *SeekSim) update(win *pixelgl.Window, dt float64) {
	if !s.cart.IntersectCircle(s.goal).Eq(pixel.ZV) {
		// Intersection!
		useOtherTargetPos = !useOtherTargetPos
		if useOtherTargetPos {
			s.goal.Center = otherTargetPos
		} else {
			s.goal.Center = initialTargetPos
		}
	}

	if win.Pressed(pixelgl.KeyA) || win.Pressed(pixelgl.KeyD) {
		appliedForceVec := pixel.V(cartM*movementAcc, 0)
		if win.Pressed(pixelgl.KeyA) {
			// Move cart left
			appliedForceVec = appliedForceVec.Scaled(-1)
		}
		accVec := appliedForceVec.Scaled(1.0 / cartM)
		s.state.cartVel = s.state.cartVel.Add(accVec.Scaled(dt))
	}
	// Friction
	s.state.cartVel = s.state.cartVel.Scaled(frictionCoeff)

	newMin := s.cart.Min.Add(s.state.cartVel.Scaled(dt))
	newMax := s.cart.Max.Add(s.state.cartVel.Scaled(dt))
	if bounds.Contains(newMin) && bounds.Contains(newMax) {
		s.cart.Min = newMin
		s.cart.Max = newMax
	}
}

func Run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Seek Game",
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
		cart.update(win, dt)

		win.Clear(colornames.Black)
		cart.draw(win)
		win.Update()
	}
}
