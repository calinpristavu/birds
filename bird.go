package main

import (
	"fmt"

	"sync"

	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const gravity = 0.2
const jumpSpeed = -5

type bird struct {
	mu sync.RWMutex

	time     int
	textures []*sdl.Texture
	y        float64
	speed    float64
	dead     bool
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		image := fmt.Sprintf("res/img/frame-%d.png", i)

		bird, err := img.LoadTexture(r, image)
		if err != nil {
			return nil, fmt.Errorf("could not load textures img: %v", err)
		}

		textures = append(textures, bird)
	}

	return &bird{textures: textures, y: 300}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.time++

	b.y -= b.speed
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	rect := sdl.Rect{
		W: 50,
		H: 43,
		X: 10,
		Y: (600 - int32(b.y)) - 43/2,
	}

	i := b.time / 10 % len(b.textures)

	err := r.Copy(b.textures[i], nil, &rect)
	if err != nil {
		return fmt.Errorf("could not copy birds: %v", err)
	}

	return nil
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.speed = jumpSpeed
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dead
}

func (b *bird) kill() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dead = true
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) hitsPipe(p *pipe) bool {
	if p.x > p.w {
		return false
	}

	safe := true
	if p.inverted && b.y < float64(p.h) {
		safe = false
	}

	if !p.inverted && b.y > float64(600-p.h) {
		safe = false
	}

	log.Printf("DEAD: \n %+v \n %+v", b, p)

	return !safe
}
