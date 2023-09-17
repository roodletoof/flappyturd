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
    screenWidth = 160
    screenHeight = 144
    playerGravity = 500
    playerJumpVel = -210
    pipeSpeed = 60
    pipeGapX = 100
    pipeGapY = 60
    pipeScreenMargin = 7
    pipeN = screenWidth / pipeGapX + 1
    keyJump = ebiten.KeySpace
    keyStartGame = ebiten.KeySpace
    keyRestartGame = ebiten.KeyR
    stateTitle = 0
    statePlaying = 1
    stateGameOver = 2
    stateNewHighScore = 3
    saveFilePath = "turd_highscore"
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
        state : stateTitle,
        player : player{},
        pipePairs: [pipeN]pipePair{},
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
    pipePairs [pipeN]pipePair
    score uint64
    highScore uint64
    distFromPipe float64
}

func (g *Game ) prepareNewGame() {
    rand.Seed(time.Now().UnixNano())
    g.score = 0
    g.player.init(v2{16,16})
    for i := range g.pipePairs {
        g.pipePairs[i].init( screenWidth + pipeGapX * float64(i) )
    }
    g.distFromPipe = g.pipePairs[0].position.x - g.player.position.x
}

var lastTime time.Time = time.Now()
func (g *Game) Update() error {

    currentTime := time.Now()
    dt := currentTime.Sub(lastTime).Seconds()
    lastTime = currentTime

    switch g.state {
    case stateTitle:
        if (inpututil.IsKeyJustPressed(keyStartGame)) {
            g.state = statePlaying
        }
    case statePlaying:


        err := g.player.move(dt)
        if err != nil { return err }

        player_rect := g.player.rect()

        var newDistFromPipe float64 = screenWidth

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

        screenRect := rect{ v2{0, 0}, v2{screenWidth, screenHeight} }
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

    case stateGameOver, stateNewHighScore:
        if (inpututil.IsKeyJustPressed(keyRestartGame)) {
            g.prepareNewGame()
            g.state = statePlaying
        }
    }

	return nil
}

func (g *Game) handleDefeat() error {
    if g.highScore < g.score {
        assets.SfxPlayNewHighscore()
        g.highScore = g.score
        g.state = stateNewHighScore
        if err := setHighScore(g.highScore); err != nil {
            return err
        }
    } else {
        assets.SfxPlayDie()
        g.state = stateGameOver
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
    case stateTitle:
        g.drawEntities(screen)
        g.printToScreen(
            screen,
            fmt.Sprint(
                "Best score: ", g.highScore, "\n",
                "Press ", keyStartGame.String(), " to start",
            ),
            8, 20,
        )
    case statePlaying:
        g.drawEntities(screen)
        g.printToScreen(screen, fmt.Sprint(g.score), 2, 9)
    case stateGameOver:
        g.drawEntities(screen)
        g.printToScreen(
            screen,
            fmt.Sprint(
                "Score: ", g.score, "\n",
                "Highscore: ", g.highScore, "\n",
                "Press ", keyRestartGame.String(), " to restart",
            ),
            2, 9,
        )
    case stateNewHighScore:
        g.drawEntities(screen)
        g.printToScreen(
            screen,
            fmt.Sprint(
                "New highscore: ", g.highScore, "\n",
                "Press ", keyRestartGame.String(), " to restart",
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
	return screenWidth, screenHeight
}



