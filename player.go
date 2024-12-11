package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/roodletoof/flappyturd/assets"
)


type player struct{
    position v2
    velocity v2
    collisionRect rect
}

func (p *player) init(position v2) {
    p.position = position
    p.velocity = v2{0,0}
}

func (p *player) move(dt float64) error {
    if (inpututil.IsKeyJustPressed(KeyJump)) {
        assets.SfxPlayJump()
        p.velocity.y = PlayerJumpVel
    }
    p.velocity.y += PlayerGravity * dt * 0.5
    p.position.y += p.velocity.y * dt
    p.velocity.y += PlayerGravity * dt * 0.5

    return nil
}
func (p *player) draw(screen *ebiten.Image) {
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(-8, -8)
    op.GeoM.Translate(p.position.x, p.position.y)
    screen.DrawImage(assets.SpritePlayer, op)
}
func (p *player) rect() *rect {
    p.collisionRect = rect{
        v2{p.position.x-6, p.position.y-7},
        v2{13,15},
    }
    return &p.collisionRect
}
