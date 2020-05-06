package view

import (
	"github.com/gdamore/tcell"
	"github.com/ramonmacias/go-spaceship-shooter/internal/game"
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
