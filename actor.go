package main

import (
	"log"
	"math"

	"github.com/EngoEngine/ecs"
)

var (
	solids []AABB
	items  []AABB
)

type Vec2i struct {
	X, Y int
}

type Vec2f struct {
	X, Y float64
}

type AABB struct {
	X, Y, W, H int
}

type Actor struct {
	Speed    Vec2f
	Pos      Vec2i
	Collider AABB

	onGround    bool
	wasOnGround bool
	moveAcc     Vec2f
}

func (a *Actor) MoveH(amount float64) bool {
	a.moveAcc.X += amount
	pixels := int(math.Floor(a.moveAcc.X))
	if pixels != 0 {
		a.moveAcc.X -= float64(pixels)
		return a.MoveHPixels(pixels)
	}
	return false
}

func (a *Actor) MoveV(amount float64) bool {
	a.moveAcc.Y += amount
	pixels := int(math.Floor(a.moveAcc.Y))
	if pixels != 0 {
		a.moveAcc.Y -= float64(pixels)
		return a.MoveVPixels(pixels)
	}
	return false
}

func (a *Actor) MoveHPixels(pixels int) bool {
	var incr int
	if math.Signbit(float64(pixels)) {
		incr = -1
	} else {
		incr = 1
	}
	for pixels != 0 {
		solid := a.checkColAtPlace(Vec2i{incr, 0}, solids)
		if solid {
			a.moveAcc.X = 0
			return true
		}
		pixels -= incr
		a.Pos.X += incr
		a.updateCollider()
	}
	return false
}

func (a *Actor) MoveVPixels(pixels int) bool {
	var incr int
	if math.Signbit(float64(pixels)) {
		incr = -1
	} else {
		incr = 1
	}
	for pixels != 0 {
		var solid bool
		if incr > 0 {
			solid = a.checkColAtPlace(Vec2i{0, incr}, solids)
		} else {
			solid = a.OnGround()
		}
		if solid {
			a.moveAcc.Y = 0
			return true
		}
		pixels -= incr
		a.Pos.Y += incr
		a.updateCollider()
	}
	return false
}

func (a *Actor) updateCollider() {
	tempSize := 16 // todo: unharcode
	a.Collider.Y = a.Pos.Y + tempSize
	a.Collider.X = a.Pos.X - (a.Collider.W / 2)
	a.Collider.W = tempSize
	a.Collider.H = tempSize
}

func (a *Actor) OnGround() bool {
	return a.checkColAtPlace(Vec2i{0, -1}, solids) // !CollisionSelf
}

func (a *Actor) checkColAtPlace(extraPos Vec2i, colliders []AABB) bool {
	place := AABB{
		X: a.Collider.X + extraPos.X,
		Y: a.Collider.Y + extraPos.Y,
		W: a.Collider.W,
		H: a.Collider.H,
	}
	for _, c := range colliders {
		if place.X < c.X+c.W &&
			place.X+place.W > c.X &&
			place.Y < c.Y+c.H &&
			place.Y+place.H > c.Y {
			return true
		}
	}
	return false
}

type actorEntity struct {
	*ecs.BasicEntity
	*Actor
}

type ActorSystem struct {
	Gravity float64

	entities []actorEntity
}

func (s *ActorSystem) New(*ecs.World) {
	s.Gravity = -1
}

func (s *ActorSystem) Update(dt float32) {
	for _, p := range s.entities {
		p.wasOnGround = p.onGround
		p.onGround = p.OnGround()

		if !p.onGround {
			p.Speed.Y += s.Gravity
		}

		if p.MoveH(p.Speed.X) {
			p.Speed.X = 0
		}

		if p.MoveV(p.Speed.Y) {
			p.Speed.Y = 0
		}

		if p.checkColAtPlace(Vec2i{}, items) {
			log.Println("RESET")
			reset()
		}
	}
}

func (s *ActorSystem) Add(basic *ecs.BasicEntity, actor *Actor) {
	s.entities = append(s.entities, actorEntity{basic, actor})
}

func (s *ActorSystem) Remove(basic ecs.BasicEntity) {
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
