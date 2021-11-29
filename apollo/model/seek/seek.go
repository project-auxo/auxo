package seek

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	zmq "github.com/pebbe/zmq4"
	"golang.org/x/image/colornames"
	"google.golang.org/protobuf/proto"

	pb "github.com/project-auxo/auxo/apollo/model/seek/proto"
	"github.com/project-auxo/auxo/olympus/logging"
)

const (
	Hostname      = "*"
	Port          = 5559
	CommandPort   = 5560
	PublishRate   = time.Second / 120
	StateTopic    = "seek/state"
	thickness     = 3
	cartWidth     = 150
	cartHeight    = 20
	radius        = 25
	movementAcc   = 300
	frictionCoeff = 0.99
	cartM         = 5 // Kg
)

var (
	bounds            = pixel.R(0, 0, 1024, 768)
	initialCartPos    = pixel.V(200, 200)
	initialTargetPos  = pixel.V(900, initialCartPos.Y)
	otherTargetPos    = pixel.V(200, initialCartPos.Y)
	useOtherTargetPos = false
	sim               = new()
	log               = logging.Base()
)

type Cart struct {
	cart    pixel.Rect
	cartVel pixel.Vec
}

type SeekSim struct {
	cart Cart
	goal pixel.Circle
}

func new() *SeekSim {
	cart := Cart{
		cart: pixel.Rect{
			Min: pixel.V(initialCartPos.X-cartWidth, initialCartPos.Y-cartHeight),
			Max: pixel.V(initialCartPos.X+cartWidth, initialCartPos.Y),
		},
		cartVel: pixel.ZV,
	}
	goal := pixel.Circle{
		Center: initialTargetPos,
		Radius: radius,
	}
	return &SeekSim{
		cart: cart,
		goal: goal,
	}
}

func (s *SeekSim) draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	cart := s.cart.cart

	// Drawing the cart
	imd.Color = color.White
	imd.Push(cart.Min, cart.Max)
	imd.Rectangle(thickness)

	// Drawing the target to seek
	imd.Color = colornames.Gold
	imd.Push(s.goal.Center)
	imd.Circle(radius, thickness)

	imd.Draw(win)
}

func (s *SeekSim) update(win *pixelgl.Window, dt float64, sock *zmq.Socket) {
	if !s.cart.cart.IntersectCircle(s.goal).Eq(pixel.ZV) {
		// Intersection!
		useOtherTargetPos = !useOtherTargetPos
		if useOtherTargetPos {
			s.goal.Center = otherTargetPos
		} else {
			s.goal.Center = initialTargetPos
		}
	}

	msg, _ := sock.RecvBytes(zmq.DONTWAIT)
	if len(msg) > 0 {
		command := &pb.Command{}
		if err := proto.Unmarshal(msg, command); err != nil {
			return
		}
		// Just send an OK bit...
		sock.SendBytes([]byte{1}, zmq.DONTWAIT)

		appliedForceVec := pixel.V(cartM*movementAcc, 0)
		if command.GetDirection() == pb.Direction_LEFT {
			appliedForceVec = appliedForceVec.Scaled(-1)
		}
		accVec := appliedForceVec.Scaled(1.0 / cartM)
		s.cart.cartVel = s.cart.cartVel.Add(accVec.Scaled(dt))
	}
	s.cart.cartVel = s.cart.cartVel.Scaled(frictionCoeff)

	newMin := s.cart.cart.Min.Add(s.cart.cartVel.Scaled(dt))
	newMax := s.cart.cart.Max.Add(s.cart.cartVel.Scaled(dt))
	if bounds.Contains(newMin) && bounds.Contains(newMax) {
		s.cart.cart.Min = newMin
		s.cart.cart.Max = newMax
	}
}

func (s *SeekSim) shareState() {
	publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatalln("failed to make publisher in order to share state.")
	}

	defer publisher.Close()
	publisher.Bind(fmt.Sprintf("tcp://%s:%d", Hostname, Port))

	stateMsg := &pb.SimState{}
	rate := time.NewTicker(PublishRate)
	for {
		stateMsg.Cart = &pb.Cart{
			CartPos: &pb.Vec{X: s.cart.cart.Min.X, Y: s.cart.cart.Max.X},
			CartVel: &pb.Vec{X: s.cart.cartVel.X, Y: s.cart.cartVel.Y},
		}
		stateMsg.GoalPos = &pb.Vec{X: s.goal.Center.X, Y: 0}
		stateBytes, err := proto.Marshal(stateMsg)
		if err != nil {
			continue
		}
		publisher.SendMessageDontwait(StateTopic, stateBytes)
		<-rate.C
	}
}

func Run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Seek Game",
		Bounds: bounds,
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatalf("failed to make simulation window: %v", err)
	}

	// Share the sim's state to any subscriber wanting to listen
	go sim.shareState()

	commandSocket, _ := zmq.NewSocket(zmq.REP)
	defer commandSocket.Close()
	commandSocket.Bind(fmt.Sprintf("tcp://%s:%d", Hostname, CommandPort))

	last := time.Now()
	fps := time.NewTicker(time.Second / 120)
	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			win.SetClosed(true)
			break
		}
		dt := time.Since(last).Seconds()
		last = time.Now()

		sim.update(win, dt, commandSocket)

		win.Clear(colornames.Black)
		sim.draw(win)
		win.Update()

		<-fps.C
	}
}
