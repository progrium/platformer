package main

import (
	"github.com/EngoEngine/ecs"
)

var replays map[*Intent][]ReplayFrame

type ReplaySystem struct {
	frameTick int64
	input     *Intent
	buffer    []ReplayFrame
	outputs   []*Intent
	replays   [][]ReplayFrame
}

func (s *ReplaySystem) New(*ecs.World) {
	if s.buffer != nil {
		s.replays = append(s.replays, s.buffer)
	}
	s.buffer = []ReplayFrame{}
	s.outputs = []*Intent{}
}

func (s *ReplaySystem) Update(dt float32) {
	s.buffer = append(s.buffer, ReplayFrame{
		Tick:   s.frameTick,
		Intent: *s.input,
	})
	for idx, frames := range s.replays {
		for _, frame := range frames {
			if frame.Tick >= s.frameTick {
				s.outputs[idx].Jump = frame.Intent.Jump
				s.outputs[idx].MoveLeft = frame.Intent.MoveLeft
				s.outputs[idx].MoveRight = frame.Intent.MoveRight
				break
			}
		}
	}
}

func (s *ReplaySystem) Remove(e ecs.BasicEntity) {

}

type ReplayFrame struct {
	Tick   int64
	Intent Intent
}
