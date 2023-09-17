package assets

import (
    _ "embed"
    "github.com/hajimehoshi/ebiten/v2"
    "bytes"
    "image"
    "log"
)

var ( // public
    SpritePlayer *ebiten.Image
    SpritePipe *ebiten.Image
)


// private


//go:embed turd.png
var turd_png []byte

//go:embed pipe.png
var pipe_png []byte

func init() {
    player_img, _, err := image.Decode(bytes.NewReader(turd_png))
    if err != nil { log.Fatal(err) }
    SpritePlayer = ebiten.NewImageFromImage(player_img)

    pipe_img, _, err := image.Decode(bytes.NewReader(pipe_png))
    if err != nil { log.Fatal(err) }
    SpritePipe = ebiten.NewImageFromImage(pipe_img)
}
