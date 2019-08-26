package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	tileSize = 16
	tileXNum = 7
)

var (
	emptyImage, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
)

func init() {
	emptyImage.Fill(color.White)
}

type Tile struct {
	GID   uint
	Layer int
	Pos   Vec2i
}

type tileEntity struct {
	*ecs.BasicEntity
	*Tile
}

type RenderSystem struct {
	screen     *ebiten.Image
	tilesImage *ebiten.Image
	entities   []tileEntity
}

func (s *RenderSystem) New(*ecs.World) {
	ts, err := os.Open("tileset.png")
	defer ts.Close()
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(ts)
	if err != nil {
		log.Fatal(err)
	}
	s.tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

func (s *RenderSystem) TileImage(gid uint) *ebiten.Image {
	x := (int(gid-1) % tileXNum) * tileSize
	y := (int(gid-1) / tileXNum) * tileSize
	return s.tilesImage.SubImage(image.Rect(x, y, x+tileSize, y+tileSize)).(*ebiten.Image)
}

func (s *RenderSystem) Update(dt float32) {
	if ebiten.IsDrawingSkipped() {
		return
	}

	layers := make(map[int][]*Tile)
	bottomLayer := 0
	for _, t := range s.entities {
		if t.Layer > bottomLayer {
			bottomLayer = t.Layer
		}
		layers[t.Layer] = append(layers[t.Layer], t.Tile)
	}

	for l := bottomLayer; l >= 0; l-- {
		tiles, ok := layers[l]
		if !ok {
			continue
		}
		for _, t := range tiles {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(t.Pos.X), float64(screenHeight-t.Pos.Y))
			s.screen.DrawImage(s.TileImage(t.GID), op)
		}
	}

	// TPS counter
	ebitenutil.DebugPrint(s.screen, fmt.Sprintf("   %.0f", math.Round(ebiten.CurrentTPS())))

}

func (s *RenderSystem) Add(basic *ecs.BasicEntity, tile *Tile) {
	s.entities = append(s.entities, tileEntity{basic, tile})
}

func (s *RenderSystem) Remove(basic ecs.BasicEntity) {
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

func line(x0, y0, x1, y1 float32, clr color.RGBA) ([]ebiten.Vertex, []uint16) {
	const width = 1

	theta := math.Atan2(float64(y1-y0), float64(x1-x0))
	theta += math.Pi / 2
	dx := float32(math.Cos(theta))
	dy := float32(math.Sin(theta))

	r := float32(clr.R) / 0xff
	g := float32(clr.G) / 0xff
	b := float32(clr.B) / 0xff
	a := float32(clr.A) / 0xff

	return []ebiten.Vertex{
		{
			DstX:   x0 - width*dx/2,
			DstY:   y0 - width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x0 + width*dx/2,
			DstY:   y0 + width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1 - width*dx/2,
			DstY:   y1 - width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1 + width*dx/2,
			DstY:   y1 + width*dy/2,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, []uint16{0, 1, 2, 1, 2, 3}
}

func rect(x, y, w, h float32, clr color.RGBA) ([]ebiten.Vertex, []uint16) {
	r := float32(clr.R) / 0xff
	g := float32(clr.G) / 0xff
	b := float32(clr.B) / 0xff
	a := float32(clr.A) / 0xff
	x0 := x
	y0 := y
	x1 := x + w
	y1 := y + h

	return []ebiten.Vertex{
		{
			DstX:   x0,
			DstY:   y0,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1,
			DstY:   y0,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x0,
			DstY:   y1,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1,
			DstY:   y1,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, []uint16{0, 1, 2, 1, 2, 3}
}
