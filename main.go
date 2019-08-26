package main

import (
	_ "image/png"
	"log"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	world     ecs.World
	renderer  *RenderSystem
	replay    *ReplaySystem
	lastTick  int64
	resetTick int64
)

type Player struct {
	ecs.BasicEntity
	*Actor
	*Intent
	*Tile
}

func init() {
	replay = &ReplaySystem{}
	reset()
}

func reset() {
	log.Printf("Replays: %d\n", len(replay.replays))
	resetTick = time.Now().UnixNano()
	world = ecs.World{}

	var players []Player
	players = append(players, Player{
		Actor: &Actor{
			Pos: Vec2i{100, 300},
		},
		Intent: &Intent{},
		Tile: &Tile{
			GID:   50,
			Layer: 0,
		},
	})

	replay.input = players[0].Intent

	renderer = &RenderSystem{}
	input := &InputSystem{}
	actor := &ActorSystem{}
	player := &PlayerSystem{}
	world.AddSystem(replay)
	world.AddSystem(renderer)
	world.AddSystem(&MapSystem{
		tiles: renderer,
	})
	world.AddSystem(input)
	world.AddSystem(actor)
	world.AddSystem(player)

	for _ = range replay.replays {
		intent := &Intent{}
		players = append(players, Player{
			Actor: &Actor{
				Pos: Vec2i{100, 300},
			},
			Intent: intent,
			Tile: &Tile{
				GID:   50,
				Layer: 0,
			},
		})
		replay.outputs = append(replay.outputs, intent)
	}

	for _, p := range players {
		renderer.Add(&p.BasicEntity, p.Tile)
		input.Add(&p.BasicEntity, p.Intent)
		actor.Add(&p.BasicEntity, p.Actor)
		player.Add(&p.BasicEntity, p.Actor, p.Intent, p.Tile)
	}

}

func main() {
	lastTick = time.Now().UnixNano()
	if err := ebiten.Run(func(screen *ebiten.Image) error {
		tickTime := time.Now().UnixNano()
		tickDelta := float32(tickTime-lastTick) / float32(time.Second)
		lastTick = tickTime

		renderer.screen = screen
		replay.frameTick = tickTime - resetTick

		world.Update(tickDelta)

		return nil
	}, screenWidth, screenHeight, 1, "Platformer"); err != nil {
		panic(err)
	}
}
