package scene

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/3auris/snakery/internal/object"
	"github.com/3auris/snakery/pkg/grafio"
)

// Scene holds paints and state of the current game
type Scene struct {
	r *sdl.Renderer
	w *sdl.Window

	drawer grafio.Drawer
	state  object.GameState
	paints map[object.GameState][]object.Paintable
}

// New create new Scene with given parameters
func New(d grafio.Drawer) (*Scene, error) {
	scrn := object.GameScreen{W: d.ScreenWidth(), H: d.ScreenHeight()}

	apple := object.NewApple()
	score := object.NewScore()
	snake := object.NewSnake(apple, score, scrn)
	deadScreen := &object.DeadScreen{Score: score, Screen: scrn}
	menuScreen := &object.WelcomeText{Screen: scrn, Snake: snake}

	return &Scene{
		drawer: d,

		state: object.MenuScreen,
		paints: map[object.GameState][]object.Paintable{
			object.MenuScreen:   {menuScreen},
			object.SnakeRunning: {snake, apple, score},
			object.DeadSnake:    {deadScreen},
		},
	}, nil
}

// Run runs goroutine and updates all paints and listening of events
func (s *Scene) Run(events <-chan sdl.Event) <-chan error {
	errc := make(chan error)

	go func() {
		ticker := time.Tick(55 * time.Millisecond)

		for {
			select {
			case e := <-events:
				for _, paint := range s.paints[s.state] {
					switch p := paint.(type) {
					case object.Handleable:
						p.HandleEvent(e)
					}
				}

				if done := s.handleExit(e); done {
					os.Exit(0)
				}
			case <-ticker:
				s.state = s.update()

				if err := s.paint(); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *Scene) handleExit(event sdl.Event) bool {
	switch ev := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		if ev.State != sdl.PRESSED {
			break
		}

		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_ESCAPE:
			return true
		}
	}
	return false
}

func (s Scene) update() object.GameState {
	for _, paint := range s.paints[s.state] {
		switch p := paint.(type) {
		case object.Updateable:
			state := p.Update()
			if state != s.state {
				return state
			}
		}
	}

	return s.state
}

func (s Scene) paint() error {
	err := s.drawer.Present(func() error {
		for _, paint := range s.paints[s.state] {
			if err := paint.Paint(s.drawer); err != nil {
				return errors.Wrap(err, "failed to paint")
			}
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to present")
	}

	return nil
}
