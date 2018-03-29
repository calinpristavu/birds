package main

import (
	"fmt"

	"math/rand"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	x        int32
	h        int32
	w        int32
	speed    int32
	inverted bool
	texture  *sdl.Texture
}

func newPipe(r *sdl.Renderer, offset int32, inverted bool) (*pipe, error) {
	direction := "top"
	if inverted {
		direction = "bottom"
	}
	texture, err := img.LoadTexture(r, fmt.Sprintf("res/img/%sPipe.png", direction))
	if err != nil {
		return nil, fmt.Errorf("could not load texture: %v", err)
	}

	h := rand.Intn(150) + 150

	return &pipe{
		x:        offset,
		h:        int32(h),
		w:        30,
		speed:    5,
		inverted: inverted,
		texture:  texture,
	}, nil
}

func (p *pipe) paint(r *sdl.Renderer) error {
	y := int32(0)
	if p.inverted {
		y = 600 - p.h
	}
	rect := sdl.Rect{
		W: p.w,
		H: p.h,
		X: p.x,
		Y: y,
	}
	err := r.Copy(p.texture, nil, &rect)
	if err != nil {
		return fmt.Errorf("could not copy pipe: %v", err)
	}

	return nil
}

func (p *pipe) restart() {
	p.x += 400
}

func (p *pipe) update() {
	p.x -= p.speed
	if p.x < 0 {
		p.x = 800
	}
}

func (p *pipe) destroy() {
	p.texture.Destroy()
}
