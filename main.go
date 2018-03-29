package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const wHeight = 600
const wWidth = 800

func main() {
	if err := run(); err != nil {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("could not initialize ttf: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(wWidth, wHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not open window: %v", err)
	}
	defer w.Destroy()

	if err := drawTitle(r, "Flappy lalala!"); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(time.Second * 1)

	s, err := newScene(r)

	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}

	events := make(chan sdl.Event)
	errc := s.run(events, r)

	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

func drawTitle(r *sdl.Renderer, text string) error {
	r.Clear()

	font, err := ttf.OpenFont("res/fonts/font.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer font.Close()

	color := sdl.Color{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}

	s, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return fmt.Errorf("could not render color: %v", err)
	}
	defer s.Free()

	texture, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer texture.Destroy()

	err = r.Copy(texture, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy texure: %v", err)
	}

	r.Present()

	return nil
}
