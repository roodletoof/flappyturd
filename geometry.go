package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)



type v2 struct{
    x, y float64
}

type rect struct{
    position, size v2
}

func (r *rect) draw(screen *ebiten.Image, clr color.Color) {
    vector.DrawFilledRect(
        screen,
        float32(r.position.x),
        float32(r.position.y),
        float32(r.size.x),
        float32(r.size.y),
        clr,
        false,
    )
}

func (r1 *rect) overlaps(r2 *rect) bool {
    return  r1.position.x <= r2.position.x + r2.size.x -1 &&
            r1.position.y <= r2.position.y + r2.size.y -1 &&
            r2.position.x <= r1.position.x + r1.size.x -1 &&
            r2.position.y <= r1.position.y + r1.size.y -1
}
