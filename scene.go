package main

import (
	"fmt"
	"time"

	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes []*pipe
}

func newScene(r *sdl.Renderer) (*Scene, error) {
	bg, err := img.LoadTexture(r, "res/img/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load bg img: %v", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, fmt.Errorf("could not create bird: %v", nil)
	}

	var pipes []*pipe
	pipe, err := newPipe(r, 400, false)
	if err != nil {
		return nil, fmt.Errorf("could not create pipe: %v", nil)
	}
	pipes = append(pipes, pipe)

	for i := 1; i < 4; i++ {
		prevPipe := pipes[i-1]
		pipe, err := newPipe(r, prevPipe.x+200, !prevPipe.inverted)
		if err != nil {
			return nil, fmt.Errorf("could not create pipe: %v", nil)
		}
		pipes = append(pipes, pipe)
	}

	return &Scene{bg: bg, bird: bird, pipes: pipes}, nil
}

func (s *Scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-time.Tick(10 * time.Millisecond):
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "BIRD IS DEAD!")
					time.Sleep(2 * time.Second)
					s.restart()
				}

				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *Scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		log.Printf("quit event: %T", e)
		return true

	case *sdl.MouseButtonEvent:
		if e.State == sdl.RELEASED {
			log.Printf("click event: %T", e)
			s.bird.jump()
		}
	case *sdl.MouseMotionEvent:
	case *sdl.WindowEvent:
	default:
		log.Printf("unknown event: %T", e)
	}

	return false
}

func (s *Scene) update() {
	s.bird.update()
	for i := range s.pipes {
		s.pipes[i].update()
		if s.bird.hitsPipe(s.pipes[i]) {
			s.bird.kill()
		}
	}
}

func (s *Scene) paint(r *sdl.Renderer) error {
	r.Clear()

	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy bg: %v", err)
	}

	s.bird.paint(r)
	for i := range s.pipes {
		s.pipes[i].paint(r)
	}

	r.Present()

	return nil
}

func (s *Scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()

	for i := range s.pipes {
		s.pipes[i].destroy()
	}
}

func (s *Scene) restart() {
	s.bird.restart()

	for i := range s.pipes {
		s.pipes[i].restart()
	}
}
