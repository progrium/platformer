package main

import (
	"log"
	"os"

	"github.com/EngoEngine/ecs"
	"github.com/Noofbiz/tmx"
)

type MapSystem struct {
	tiles   *RenderSystem
	tileMap tmx.Map
}

func (s *MapSystem) New(*ecs.World) {
	f, err := os.Open("debug.tmx")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	s.tileMap, err = tmx.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	const xNum = screenWidth / tileSize
	for idx, l := range s.tileMap.Layers {
		for i, t := range l.Data[0].Tiles {
			x := int((i % xNum) * tileSize)
			y := int((i / xNum) * tileSize)
			for _, ii := range []uint32{1, 2, 3, 4, 5, 6, 7, 14, 18, 19, 20, 21, 25, 26, 27, 28, 32, 33, 34, 35, 42} {
				if t.GID == ii {
					solids = append(solids, AABB{X: x, Y: screenHeight - y, W: tileSize, H: tileSize})
					break
				}
			}
			if t.GID == 56 {
				items = append(items, AABB{X: x, Y: screenHeight - y, W: tileSize, H: tileSize})
			}
			e := tileEntity{
				Tile: &Tile{
					GID:   uint(t.GID),
					Layer: len(s.tileMap.Layers) - idx,
					Pos:   Vec2i{x, screenHeight - y},
				},
			}
			s.tiles.Add(e.BasicEntity, e.Tile)
		}
	}
}

func (s *MapSystem) Update(dt float32) {

}

// func (s *PlayerSystem) Add(basic *ecs.BasicEntity, actor *Actor, intent *Intent, tile *Tile) {

// }

func (s *MapSystem) Remove(basic ecs.BasicEntity) {
}
