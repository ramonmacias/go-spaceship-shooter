package view

import (
	"github.com/gdamore/tcell"
	"github.com/ramonmacias/go-spaceship-shooter/internal/game"
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
