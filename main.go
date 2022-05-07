package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"time"
)

func iniWindow() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("玩数学")
}

func main() {
	go PlaySoundWelcome()
	rand.Seed(time.Now().UnixNano())
	iniWindow()
	err := ebiten.RunGame(NewGame())
	if err != nil {
		log.Fatal(err)
	}
}
