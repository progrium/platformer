package main

import (
	"log"
	"os"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	joystickThreshold = 0.15

	buttonA       = 11
	leftShoulder  = 8
	rightShoulder = 9
)

type InputSystem struct {
	gamepadIDs    map[int]struct{}
	intents       []*Intent
	lastAxis      float64
	currentPlayer int
}

func (s *InputSystem) Add(basic *ecs.BasicEntity, intent *Intent) {
	s.intents = append(s.intents, intent)
}

func (s *InputSystem) New(*ecs.World) {
	s.gamepadIDs = make(map[int]struct{})
}

func (s *InputSystem) Update(dt float32) {
	// Log the gamepad connection events.
	for _, id := range inpututil.JustConnectedGamepadIDs() {
		log.Printf("gamepad connected: id: %d", id)
		s.gamepadIDs[id] = struct{}{}
	}
	for id := range s.gamepadIDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			log.Printf("gamepad disconnected: id: %d", id)
			delete(s.gamepadIDs, id)
		}
	}

	for _, id := range ebiten.GamepadIDs() {
		// maxAxis := ebiten.GamepadAxisNum(id)
		// for a := 0; a < maxAxis; a++ {
		// 	v :=
		// 	axes[id] = append(axes[id], fmt.Sprintf("%d:%0.2f", a, v))
		// }
		maxButton := ebiten.GamepadButton(ebiten.GamepadButtonNum(id))
		for b := ebiten.GamepadButton(id); b < maxButton; b++ {
			// Log button events.
			// if inpututil.IsGamepadButtonJustPressed(id, b) {
			// 	log.Printf("button pressed: id: %d, button: %d", id, b)
			// }
			// if inpututil.IsGamepadButtonJustReleased(id, b) {
			// 	log.Printf("button released: id: %d, button: %d", id, b)
			// }
		}
	}

	// Controls
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.GamepadAxis(0, 0) < -joystickThreshold {
		s.intents[s.currentPlayer].MoveLeft = true
		s.intents[s.currentPlayer].MoveRight = false
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.GamepadAxis(0, 0) > joystickThreshold {
		s.intents[s.currentPlayer].MoveLeft = false
		s.intents[s.currentPlayer].MoveRight = true
	} else {
		s.intents[s.currentPlayer].MoveLeft = false
		s.intents[s.currentPlayer].MoveRight = false
	}

	// jump
	s.intents[s.currentPlayer].Jump = inpututil.IsGamepadButtonJustPressed(0, buttonA)

	// switch player
	if inpututil.IsGamepadButtonJustPressed(0, leftShoulder) {
		s.currentPlayer--
		if s.currentPlayer < 0 {
			s.currentPlayer = len(s.intents) - 1
		}
	}
	if inpututil.IsGamepadButtonJustPressed(0, rightShoulder) {
		s.currentPlayer++
		if s.currentPlayer == len(s.intents) {
			s.currentPlayer = 0
		}
	}

	// exit
	if inpututil.IsGamepadButtonJustPressed(0, 4) {
		os.Exit(0)
	}
}

func (s *InputSystem) Remove(e ecs.BasicEntity) {

}
