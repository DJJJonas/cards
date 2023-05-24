package main

import (
	"cards/cards"

	"github.com/gin-gonic/gin"
)

var b *cards.Board

func main() {
	// ? Game Setup
	b = &cards.Board{
		PlayerTurn:        0,
		TurnCount:         0,
		Players:           [2]*cards.Player{cards.PaladinControlZoth(), cards.PaladinControlZoth()},
		ActionChan:        make(chan *cards.Action),
		WaitingActionChan: make(chan int, 1),
		ActionEndChan:     make(chan error, 1),
	}

	go b.Start()

	// ? Server Setup
	eng := gin.Default()
	eng.Any("/ws/connect", Connect)
	eng.StaticFile("/", "./frontend/dist/index.html")
	eng.Static("/assets", "./frontend/dist/assets/")
	eng.Run(":8123")
}
