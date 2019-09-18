package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"

	"github.com/EngoEngine/ecs"
)

type ReplaySystem struct {
	input     *Intent
	lastInput *Intent
	buffer    []ReplayFrame
	outputs   []*Intent
	replays   [][]ReplayFrame

	conn net.Conn
}

func (s *ReplaySystem) New(*ecs.World) {
	if len(s.buffer) > 1 {
		s.replays = append(s.replays, s.buffer)
	}
	s.buffer = []ReplayFrame{
		ReplayFrame{
			Tick:   0,
			Intent: Intent{},
		},
	}
	s.outputs = []*Intent{}
	s.lastInput = nil
	var err error
	s.conn, err = net.Dial("udp", ":3333")
	if err != nil {
		panic(err)
	}
}

func (s *ReplaySystem) Update(dt float32) {
	if len(s.outputs) == 0 && len(s.replays) > 0 {
		for _ = range s.replays {
			intent := &Intent{}
			game.AddPlayer(Player{
				Actor: &Actor{
					Pos: Vec2i{100, 300},
				},
				Intent: intent,
				Tile: &Tile{
					GID:   50,
					Layer: 0,
				},
			})
			s.outputs = append(s.outputs, intent)
		}
	}

	for idx, frames := range s.replays {
		for _, frame := range frames {
			if frame.Tick == game.Tick {
				//log.Println(game.Tick, frame)
				s.outputs[idx].Jump = frame.Intent.Jump
				s.outputs[idx].MoveLeft = frame.Intent.MoveLeft
				s.outputs[idx].MoveRight = frame.Intent.MoveRight
				break
			}
		}
	}

	frame := ReplayFrame{
		Tick:   game.Tick,
		Intent: *s.input,
	}

	if s.lastInput != nil && EqualIntents(*s.lastInput, frame.Intent) {
		return
	}

	s.lastInput = &frame.Intent
	s.buffer = append(s.buffer, frame)

	var netbuf bytes.Buffer
	enc := gob.NewEncoder(&netbuf)
	err := enc.Encode(frame)
	if err != nil {
		log.Fatal(err)
	}
	s.conn.Write(netbuf.Bytes())

}

func (s *ReplaySystem) Remove(e ecs.BasicEntity) {

}

type ReplayFrame struct {
	Tick   int64
	Intent Intent
}
