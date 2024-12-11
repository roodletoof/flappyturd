package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/roodletoof/flappyturd/assets"
)

const (
    ScreenWidth = 160
    ScreenHeight = 144
    PlayerGravity = 500
    PlayerJumpVel = -210
    PipeSpeed = 60
    PipeGapX = 100
    PipeGapY = 60
    PipeScreenMargin = 7
    PipeN = ScreenWidth / PipeGapX + 1
    KeyJump = ebiten.KeySpace
    KeyStartGame = ebiten.KeySpace
    KeyRestartGame = ebiten.KeyR
    StateTitle = 0
    StatePlaying = 1
    StateGameOver = 2
    StateNewHighScore = 3
)

func init() {
    ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
    ebiten.SetTPS(ebiten.SyncWithFPS)
}

func main() {

    highScore, err := getHighScore()
	if err != nil { log.Fatal(err) }

    var game Game = Game {
        palette : color.Palette{
            color.RGBA{155,188,15,255},
            color.RGBA{139,172,15,255},
            color.RGBA{48,98,48,255},
            color.RGBA{15,56,15,255},
        },
        state : StateTitle,
        player : player{},
        pipePairs: [PipeN]pipePair{},
        highScore: highScore,
    }
    game.prepareNewGame()

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

type Game struct{
    player player
    palette color.Palette
    state int
    pipePairs [PipeN]pipePair
    score uint64
    highScore uint64
    distFromPipe float64
}

func (g *Game ) prepareNewGame() {
    rand.Seed(time.Now().UnixNano())
    g.score = 0
    g.player.init(v2{16,16})
    for i := range g.pipePairs {
        g.pipePairs[i].init( ScreenWidth + PipeGapX * float64(i) )
    }
    g.distFromPipe = g.pipePairs[0].position.x - g.player.position.x
}

var lastTime time.Time = time.Now()
func (g *Game) Update() error {

    currentTime := time.Now()
    dt := currentTime.Sub(lastTime).Seconds()
    lastTime = currentTime

    switch g.state {
    case StateTitle:
        if (inpututil.IsKeyJustPressed(KeyStartGame)) {
            g.state = StatePlaying
        }
    case StatePlaying:


        err := g.player.move(dt)
        if err != nil { return err }

        player_rect := g.player.rect()

        var newDistFromPipe float64 = ScreenWidth

        for i := range g.pipePairs {
            pp := &g.pipePairs[i]
            err = pp.move(dt)
            if err != nil { return err }

            dist := pp.position.x - g.player.position.x
            if 0 < dist && dist < newDistFromPipe {
                newDistFromPipe = dist
            }

            for _, ppr := range pp.rect() {
                if player_rect.overlaps( &ppr ) {
                    if err := g.handleDefeat(); err != nil {
                        return err
                    }
                }
            }

        }

        screenRect := rect{ v2{0, 0}, v2{ScreenWidth, ScreenHeight} }
        if !player_rect.overlaps(&screenRect) {
            if err := g.handleDefeat(); err != nil {
                return err
            }
        }
        
        if newDistFromPipe > g.distFromPipe {
            g.score += 1
            assets.SfxPlayPoint()
        }
        g.distFromPipe = newDistFromPipe

    case StateGameOver, StateNewHighScore:
        if (inpututil.IsKeyJustPressed(KeyRestartGame)) {
            g.prepareNewGame()
            g.state = StatePlaying
        }
    }

	return nil
}

func (g *Game) handleDefeat() error {
    if g.highScore < g.score {
        assets.SfxPlayNewHighscore()
        g.highScore = g.score
        g.state = StateNewHighScore
        if err := setHighScore(g.highScore); err != nil {
            return err
        }
    } else {
        assets.SfxPlayDie()
        g.state = StateGameOver
    }
    return nil
}

func (g *Game) printToScreen(screen *ebiten.Image, msg string, x int, y int) {
    for xOff := -1; xOff <= 1; xOff++ {
        for yOff := -1; yOff <= 1; yOff++ {
            // draws outline
            text.Draw( screen, msg, assets.FontMonoJewesbury, x+xOff, y+yOff, g.palette[0])
        }
    }
    text.Draw( screen, msg, assets.FontMonoJewesbury, x, y, g.palette[3])
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(g.palette[0])

    switch g.state {
    case StateTitle:
        g.drawEntities(screen)
        g.printToScreen(
            screen,
            fmt.Sprint(
                "Best score: ", g.highScore, "\n",
                "Press ", KeyStartGame.String(), " to start",
            ),
            8, 20,
        )
    case StatePlaying:
        g.drawEntities(screen)
        g.printToScreen(screen, fmt.Sprint(g.score), 2, 9)
    case StateGameOver:
        g.drawEntities(screen)
        g.printToScreen(
            screen,
            fmt.Sprint(
                "Score: ", g.score, "\n",
                "Highscore: ", g.highScore, "\n",
                "Press ", KeyRestartGame.String(), " to restart",
            ),
            2, 9,
        )
    case StateNewHighScore:
        g.drawEntities(screen)
        g.printToScreen(
            screen,
            fmt.Sprint(
                "New highscore: ", g.highScore, "\n",
                "Press ", KeyRestartGame.String(), " to restart",
            ),
            2, 9,
        )
    }
}

func (g *Game) drawEntities(screen *ebiten.Image) {
    g.player.draw(screen)
    for i := range g.pipePairs {
        g.pipePairs[i].draw(screen)
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}



