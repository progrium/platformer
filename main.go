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
	game   *Game
	replay *ReplaySystem

	lastTime  int64
	resetTime int64
)

type Player struct {
	ecs.BasicEntity
	*Actor
	*Intent
	*Tile
}

type Game struct {
	Tick  int64
	World ecs.World

	Render *RenderSystem
	Input  *InputSystem
	Actor  *ActorSystem
	Player *PlayerSystem
	Replay *ReplaySystem

	Players []Player
}

func NewGame() *Game {
	world := ecs.World{}
	renderer := &RenderSystem{}
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

	return &Game{
		World:  world,
		Render: renderer,
		Input:  input,
		Actor:  actor,
		Player: player,
		Replay: replay,
	}
}

func (g *Game) AddPlayer(p Player) {
	log.Println("Player added")
	g.Players = append(g.Players, p)
	g.Render.Add(&p.BasicEntity, p.Tile)
	g.Input.Add(&p.BasicEntity, p.Intent)
	g.Actor.Add(&p.BasicEntity, p.Actor)
	g.Player.Add(&p.BasicEntity, p.Actor, p.Intent, p.Tile)

	// always for now
	g.Replay.input = g.Players[0].Intent
}

func init() {
	replay = &ReplaySystem{}
	reset()
}

func reset() {
	resetTime = time.Now().UnixNano()
	game = NewGame()
	game.Tick = 0
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

}

func main() {
	go ListenAndServe()
	lastTime = time.Now().UnixNano()
	if err := ebiten.Run(func(screen *ebiten.Image) error {
		game.Tick++

		tickTime := time.Now().UnixNano()
		tickDelta := float32(tickTime-lastTime) / float32(time.Second)
		lastTime = tickTime

		game.Render.screen = screen

		game.World.Update(tickDelta)

		return nil
	}, screenWidth, screenHeight, 1, "Platformer"); err != nil {
		panic(err)
	}
}
