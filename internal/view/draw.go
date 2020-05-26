package view

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/ramonmacias/go-spaceship-shooter/internal/game"
	"github.com/rivo/tview"
)

const (
	backgroundColor = tcell.Color234
	wallColor       = tcell.Color24
	playerColor     = tcell.ColorBlue
	laserColor      = tcell.ColorRed
	textColor       = tcell.ColorWhite
	botColor        = tcell.ColorMediumAquamarine
)

type drawFunc func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int)
type drawCallback func()

// draw method will receive a variadric of draw funcs and will apply each one
// on the viewPort
func (ui *UserInterface) draw(drawFuncs ...drawFunc) {
	ui.viewPort.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for _, f := range drawFuncs {
			f(screen, x, y, width, height)
		}
		return 0, 0, 0, 0
	})
}

// setupDrawCallbacks will receive a variadric of drawcallbacks function and will
// add them into the user interface structure
func (ui *UserInterface) setupDrawCallbacks(callbacks ...drawCallback) {
	for _, f := range callbacks {
		ui.drawCallbacks = append(ui.drawCallbacks, f)
	}
}

// drawMap will render the game map into your terminal
func (ui *UserInterface) drawMap() drawFunc {
	return drawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		ui.Engine.Mu.RLock()
		defer ui.Engine.Mu.RUnlock()
		style := tcell.StyleDefault.Background(backgroundColor)
		// Re visit this center stuff
		centerX := width / 2
		centerY := height / 2
		for _, wall := range ui.Engine.GameMap.GetMapElements()[game.MapElementWall] {
			x := centerX + wall.X
			y := centerY + wall.Y
			screen.SetContent(x, y, '█', nil, style.Foreground(wallColor))
		}
		return 0, 0, 0, 0
	})
}

// drawActors will render all the actors involved on the game
func (ui *UserInterface) drawActors() drawFunc {
	return drawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		ui.Engine.Mu.RLock()
		defer ui.Engine.Mu.RUnlock()
		style := tcell.StyleDefault.Background(backgroundColor)
		// Re visit this center stuff
		centerX := width / 2
		centerY := height / 2
		for _, actor := range ui.Engine.Actors {
			x := centerX + actor.Position.X
			y := centerY + actor.Position.Y

			screen.SetContent(x, y, 'A', nil, style.Foreground(playerColor))
		}
		return 0, 0, 0, 0
	})
}

// drawLasers will render all the active lasers
func (ui *UserInterface) drawLasers() drawFunc {
	return drawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		ui.Engine.Mu.RLock()
		defer ui.Engine.Mu.RUnlock()
		style := tcell.StyleDefault.Background(backgroundColor)
		// Re visit this center stuff
		centerX := width / 2
		centerY := height / 2
		ui.Engine.Lasers.Range(func(laserID interface{}, la interface{}) bool {
			laser := la.(game.Laser)
			x := centerX + laser.Position.X
			y := centerY + laser.Position.Y

			screen.SetContent(x, y, 'X', nil, style.Foreground(laserColor))
			return true
		})
		return 0, 0, 0, 0
	})
}

// drawBots will render all the active bots
func (ui *UserInterface) drawBots() drawFunc {
	return drawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		ui.Engine.Mu.RLock()
		defer ui.Engine.Mu.RUnlock()
		style := tcell.StyleDefault.Background(backgroundColor)
		// Re visit this center stuff
		centerX := width / 2
		centerY := height / 2
		ui.Engine.Bots.Range(func(botID interface{}, value interface{}) bool {
			bot := value.(game.Bot)
			x := centerX + bot.Position.X
			y := centerY + bot.Position.Y

			screen.SetContent(x, y, 'Y', nil, style.Foreground(botColor))
			return true
		})
		return 0, 0, 0, 0
	})
}

// setupScore will render a modal with a ranked players and their scores
func (ui *UserInterface) setupScore() drawCallback {
	tv := tview.NewTextView()
	tv.SetBorder(true).SetTitle("Score").SetBackgroundColor(backgroundColor)
	modal := centeredModal(tv)
	ui.pages.AddPage("score", modal, true, false)
	return func() {
		ui.Engine.Mu.RLock()
		defer ui.Engine.Mu.RUnlock()
		var text string
		for actorID, actor := range ui.Engine.Actors {
			score := ui.Engine.Score[actorID]
			text += fmt.Sprintf("%s - %d\n", actor.Name, score)
		}
		tv.SetText(text)
	}
}

// setupLevelComplete will render a final modal showing the name of the winning
// player
func (ui *UserInterface) setupLevelComplete() drawCallback {
	tv := tview.NewTextView()
	tv.SetTextAlign(tview.AlignCenter).
		SetScrollable(true).
		SetBorder(true).
		SetBackgroundColor(backgroundColor).
		SetTitle("Level complete")
	modal := centeredModal(tv)
	ui.pages.AddPage("levelComplete", modal, true, false)
	return func() {
		ui.Engine.Mu.RLock()
		defer ui.Engine.Mu.RUnlock()
		if ui.Engine.LevelComplete {
			ui.pages.ShowPage("levelComplete")
			player := ui.Engine.Actors[ui.Engine.RoundWinner]
			text := fmt.Sprintf("\nCongratulations %s you are the winner!!\n\n", player.Name)
			tv.SetText(text)
		}
	}
}

// setupGameOver will render a final modal when the main player dies
func (ui *UserInterface) setupGameOver() drawCallback {
	tv := tview.NewTextView()
	tv.SetTextAlign(tview.AlignCenter).
		SetScrollable(true).
		SetBorder(true).
		SetBackgroundColor(backgroundColor).
		SetTitle("GAME OVER")
	modal := centeredModal(tv)
	ui.pages.AddPage("gameOver", modal, true, false)
	return func() {
		ui.Engine.Mu.RLock()
		defer ui.Engine.Mu.RUnlock()
		if ui.Engine.GameOver {
			ui.pages.ShowPage("gameOver")
			text := "\nThis is the end of your adventure, try again\n\n"
			tv.SetText(text)
		}
	}
}

func centeredModal(p tview.Primitive) tview.Primitive {
	width := 0
	height := 0
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, false).
			AddItem(nil, 0, 1, false), width, 1, false).
		AddItem(nil, 0, 1, false)
}
