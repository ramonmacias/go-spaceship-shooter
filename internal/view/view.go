package view

import (
	"log"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gofrs/uuid"
	"github.com/ramonmacias/go-spaceship-shooter/internal/game"
	"github.com/rivo/tview"
)

const (
	backgroundColor = tcell.Color234
	wallColor       = tcell.Color24
	playerColor     = tcell.ColorBlue
	laserColor      = tcell.ColorRed
	drawFrequency   = 17 * time.Millisecond
	textColor       = tcell.ColorWhite
)

// UserInterface will keep the basics for render the game on a terminal and listen
// for all the interacionts from the user
type UserInterface struct {
	Engine       *game.Engine
	App          *tview.Application
	ErrChan      chan error
	pages        *tview.Pages
	viewPort     *tview.Box
	MainPlayerID uuid.UUID
}

// New function will build a new View with the basics intialized
func New(engine *game.Engine) *UserInterface {
	app := tview.NewApplication()
	pages := tview.NewPages()
	ui := &UserInterface{
		Engine:  engine,
		App:     app,
		pages:   pages,
		ErrChan: make(chan error),
	}
	ui.drawViewPort()
	ui.draw(
		ui.drawActors(),
		ui.drawMap(),
		ui.drawLasers(),
	)
	ui.setupListeners()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			select {
			case ui.ErrChan <- nil:
			default:
			}
		}
		return event
	})
	app.SetRoot(pages, true)
	return ui
}

// Start will setup the basics and run the game on your terminal
func (ui *UserInterface) Start() {
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

	helpText := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("← → ↑ ↓ move - wasd shoot - p score - esc close - ctrl+q quit").
		SetTextColor(textColor)
	helpText.SetBackgroundColor(backgroundColor)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(box, 0, 1, true).
		AddItem(helpText, 1, 1, false)
	ui.pages.AddPage("viewport", flex, true, true)
	ui.viewPort = box
}

// setupListeners will take care of all the inputs we receive from the user
// and apply the related actions
func (ui *UserInterface) setupListeners() {
	ui.viewPort.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var direction game.Direction
		switch event.Key() {
		case tcell.KeyUp:
			direction = game.DirectionUp
		case tcell.KeyDown:
			direction = game.DirectionDown
		case tcell.KeyRight:
			direction = game.DirectionRight
		case tcell.KeyLeft:
			direction = game.DirectionLeft
		}
		ui.Engine.ActionChan <- &game.MoveAction{
			ActorID:   ui.MainPlayerID,
			Direction: direction,
			CreatedAt: time.Now(),
		}

		var laserDirection game.Direction
		switch event.Rune() {
		case 'w':
			laserDirection = game.DirectionUp
		case 'd':
			laserDirection = game.DirectionRight
		case 's':
			laserDirection = game.DirectionDown
		case 'a':
			laserDirection = game.DirectionLeft
		}
		laserID := uuid.Must(uuid.NewV4())
		ui.Engine.Lasers.Store(laserID, game.Laser{
			ID:       laserID,
			Position: ui.Engine.Actors[ui.MainPlayerID].Position,
		})
		ui.Engine.ActionChan <- &game.LaserAction{
			LaserID:   laserID,
			Direction: laserDirection,
			CreatedAt: time.Now(),
		}
		return event
	})
}
