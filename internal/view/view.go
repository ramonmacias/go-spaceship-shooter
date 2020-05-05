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
	drawFrequency   = 17 * time.Millisecond
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
	log.Println("[DEBUG] Start user interface")
	ui.drawViewPort()
	// ui.drawMap()
	// box = ui.drawActors(box)
	// ui.setupListeners()

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

	box.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		ui.Engine.Mu.RLock()
		defer ui.Engine.Mu.RUnlock()
		style := tcell.StyleDefault.Background(backgroundColor)
		log.Println("[DEBUG] Printing map")
		// Re visit this center stuff
		centerX := width / 2
		centerY := height / 2
		for _, wall := range ui.Engine.GameMap.GetMapElements()[game.MapElementWall] {
			x := centerX + wall.X
			y := centerY + wall.Y
			screen.SetContent(x, y, '█', nil, style.Foreground(wallColor))
		}
		// style := tcell.StyleDefault.Background(backgroundColor)
		// Re visit this center stuff
		// centerX := width / 2
		// centerY := height / 2
		for _, actor := range ui.Engine.Actors {
			x := centerX + actor.Position.X
			y := centerY + actor.Position.Y

			screen.SetContent(x, y, 'A', nil, style.Foreground(playerColor))
		}
		return 0, 0, 0, 0
	})
	box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
		log.Println("[DEBUG] received movement", direction)
		ui.Engine.ActionChan <- &game.MoveAction{
			ActorID:   ui.MainPlayerID,
			Direction: direction,
			CreatedAt: time.Now(),
		}
		return event
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(box, 0, 1, true)
	ui.pages.AddPage("viewport", flex, true, true)
	ui.viewPort = box
}

// drawMap will render the game map into your terminal
func (ui *UserInterface) drawMap() *tview.Box {
	return ui.viewPort.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
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
		// style := tcell.StyleDefault.Background(backgroundColor)
		// Re visit this center stuff
		// centerX := width / 2
		// centerY := height / 2
		for _, actor := range ui.Engine.Actors {
			x := centerX + actor.Position.X
			y := centerY + actor.Position.Y

			screen.SetContent(x, y, 'A', nil, style.Foreground(playerColor))
		}
		return 0, 0, 0, 0
	})
}

// drawActors will render all the actors involved on the game
func (ui *UserInterface) drawActors(box *tview.Box) *tview.Box {
	return box.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
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
		log.Println("[DEBUG] received movement", direction)
		ui.Engine.ActionChan <- &game.MoveAction{
			ActorID:   ui.MainPlayerID,
			Direction: direction,
			CreatedAt: time.Now(),
		}
		return event
	})
}
