package assets

import (
	"bytes"
	_ "embed"
	"log"

	_ "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// public

func SfxPlayJump() error {
    return playSfx(jumpWavPlayer)
}

func SfxPlayDie() error {
    return playSfx(dieWavPlayer)
}

func SfxPlayNewHighscore() error {
    return playSfx(newHighscoreWavPlayer)
}

func SfxPlayPoint() error {
    return playSfx(pointWavPlayer)
}


// private

func playSfx(player *audio.Player) error {
    if err := player.Rewind(); err != nil {
        return err
    }
    player.Play()
    return nil
}


const (
    sampleRate = 22050
)

var (
    audioContext = audio.NewContext(sampleRate)

    //go:embed jump.wav
    jump_wav []byte

    //go:embed die.wav
    die_wav []byte

    //go:embed new_highscore.wav
    new_highscore_wav []byte

    //go:embed point.wav
    point_wav []byte

    jumpWavPlayer *audio.Player
    dieWavPlayer *audio.Player
    newHighscoreWavPlayer *audio.Player
    pointWavPlayer *audio.Player
)

func init() {
    jumpWavPlayer = createAudioPlayer(jump_wav)
    dieWavPlayer = createAudioPlayer(die_wav)
    newHighscoreWavPlayer = createAudioPlayer(new_highscore_wav)
    pointWavPlayer = createAudioPlayer(point_wav)
}

func createAudioPlayer(b []byte) *audio.Player {
    dec, err := wav.DecodeWithoutResampling(bytes.NewReader(b))
    if err != nil { log.Fatal(err) }
    player, err := audioContext.NewPlayer(dec)
    if err != nil { log.Fatal(err) }
    return player
}
