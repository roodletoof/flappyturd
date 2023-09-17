package assets

import (
	_ "embed"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var ( // public
    FontMonoJewesbury font.Face
)


// private

//go:embed MonoJewesburyZX.ttf
var	monoJewesbury_ttf []byte

func init() {
    tt, err := opentype.Parse(monoJewesbury_ttf)
    if err != nil { log.Fatal(err) }

    FontMonoJewesbury, err = opentype.NewFace(
        tt,
        &opentype.FaceOptions{
            DPI: 72,
            Size: 16,
            Hinting: font.HintingNone,

        },
    )
    if err != nil { log.Fatal(err) }
}
