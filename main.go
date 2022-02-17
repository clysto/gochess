package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(760*3/4, 840*3/4)
	ebiten.SetWindowTitle("中国象棋")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
