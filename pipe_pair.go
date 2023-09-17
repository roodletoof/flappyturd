package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/roodletoof/flappyturd/assets"
)


type pipePair struct{
    position v2
    collisionRects [2] rect
}

func (pp *pipePair) init(x float64) {
    pp.position.x = x
    pp.setRandomYPos()
}

func (pp *pipePair) setRandomYPos() {
    pp.position.y = pipeScreenMargin + pipeGapY * 0.5 + (screenHeight - pipeGapY - 2 * pipeScreenMargin) * rand.Float64()
}

func (pp *pipePair) move(dt float64) error {
    pp.position.x -= pipeSpeed * dt
    if (pp.position.x <= -8) {
        pp.position.x = pipeN * pipeGapX -8
        pp.setRandomYPos()
    }
    return nil
}

var (
    spriteLip = assets.SpritePipe.SubImage(image.Rect(0,0,16,6)).(*ebiten.Image)
    spriteShaft = assets.SpritePipe.SubImage(image.Rect(0,6,16,7)).(*ebiten.Image)
)
func (pp *pipePair) draw(screen *ebiten.Image) {

    lowerY := pp.position.y + pipeGapY / 2
    upperY := pp.position.y - pipeGapY / 2

    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(-8, 0)
    op.GeoM.Scale(1, screenHeight - lowerY)
    op.GeoM.Translate(pp.position.x, lowerY)
    screen.DrawImage(spriteShaft, op)

    op.GeoM.Reset()

    op.GeoM.Translate(-8, 0)
    op.GeoM.Translate(pp.position.x, lowerY)
    screen.DrawImage(spriteLip, op)

    op.GeoM.Reset()

    op.GeoM.Translate(-8, 0)
    op.GeoM.Scale(1, -upperY)
    op.GeoM.Translate(pp.position.x, upperY)
    screen.DrawImage(spriteShaft, op)

    op.GeoM.Reset()

    op.GeoM.Translate(-8, 0)
    op.GeoM.Scale(1,-1)
    op.GeoM.Translate(pp.position.x, upperY)
    screen.DrawImage(spriteLip, op)
}

func (pp *pipePair) rect() []rect {
    lowerRect := &pp.collisionRects[0]
    lowerRect.position = v2{pp.position.x-7,pp.position.y + pipeGapY / 2}
    lowerRect.size = v2{14, screenHeight - pp.position.y + pipeGapY / 2}

    upperRect := &pp.collisionRects[1]
    upperRect.position = v2{pp.position.x-7,0}
    upperRect.size = v2{14, pp.position.y - pipeGapY / 2}

    return pp.collisionRects[:]
}
