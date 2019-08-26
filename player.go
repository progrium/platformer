package main

import "github.com/EngoEngine/ecs"

type playerEntity struct {
	*ecs.BasicEntity
	*Actor
	*Intent
	*Tile
}

type PlayerSystem struct {
	JumpForce float64

	entities []playerEntity
}

func (s *PlayerSystem) New(*ecs.World) {
	s.JumpForce = 15
}

func (s *PlayerSystem) Update(dt float32) {
	for _, p := range s.entities {
		if p.MoveLeft {
			p.Speed.X = -3
		}
		if p.MoveRight {
			p.Speed.X = 3
		}
		if !p.MoveLeft && !p.MoveRight {
			p.Speed.X = 0
		}

		p.Tile.Pos.X = p.Collider.X
		p.Tile.Pos.Y = p.Collider.Y

		if p.onGround && p.Jump {
			p.Speed.Y = s.JumpForce
		}
	}
}

func (s *PlayerSystem) Add(basic *ecs.BasicEntity, actor *Actor, intent *Intent, tile *Tile) {
	s.entities = append(s.entities, playerEntity{basic, actor, intent, tile})
}

func (s *PlayerSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}
