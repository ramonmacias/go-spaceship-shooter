package view

import (
	"log"
	"time"

	"github.com/gdamore/tcell"
	"github.com/ramonmacias/go-spaceship-shooter/internal/game"
	"github.com/rivo/tview"
)

const (
	backgroundColor = tcell.Color234
	wallColor       = tcell.Color24
	drawFrequency   = 17 * time.Millisecond
)

// UserInterface will keep the basics for render the game on a terminal and listen
// for all the interacionts from the user
type UserInterface struct {
	Engine   *game.Engine
	App      *tview.Application
	ErrChan  chan error
	pages    *tview.Pages
	viewPort *tview.Box
}

// New function will build a new View with the basics intialized
func New(engine *game.Engine) *UserInterface {
	app := tview.NewApplication()
	pages := tview.NewPages()

	app.SetRoot(pages, true)
	return &UserInterface{
		Engine: engine,
		App:    app,
		pages:  pages,
	}
}

// Start will setup the basics and run the game on your terminal
func (ui *UserInterface) Start() {
	ui.drawViewPort()
	ui.drawMap()

	drawTicker := time.NewTicker(drawFrequency)
	stop := make(chan bool)
	go func() {
		for {
			ui.App.Draw()
			<-drawTicker.C
			select {
			case <-stop:
				return
			}
		}
	}()
	go func() {
		err := ui.App.Run()
		if err != nil {
			log.Println("Error starting the user interface", err)
		}
		stop <- true
		drawTicker.Stop()
		select {
		case ui.ErrChan <- err:
		}
	}()
}

// drawViewPort will render the screen where it going to start the game
func (ui *UserInterface) drawViewPort() {
	box := tview.NewBox().
		SetBorder(true).
		SetTitle("Spaceship Shooter").
		SetBackgroundColor(backgroundColor)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(box, 0, 1, true)
	ui.pages.AddPage("viewport", flex, true, true)
	ui.viewPort = box
}

// drawMap will render the game map into your terminal
func (ui *UserInterface) drawMap() {
	ui.viewPort.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
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
