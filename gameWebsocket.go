package main

import (
	"cards/cards"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Player struct {
	id              byte
	Conn            *websocket.Conn
	UpdateBoardChan chan bool
	ResultsChan     chan int
}

type Match struct {
	players [2]*Player
	board   *cards.Board
}

func (m *Match) NewPlayer(conn *websocket.Conn) *Player {
	for i, p := range m.players {
		if p == nil {
			p = &Player{
				id:              byte(i),
				Conn:            conn,
				UpdateBoardChan: make(chan bool, 1),
				ResultsChan:     make(chan int, 1),
			}
			m.players[i] = p
			return p
		}
	}
	return nil
}

func (m *Match) DisconnectPlayer(id byte) {
	if m.players[id] != nil && m.players[id].Conn != nil {
		m.players[id].Conn.Close()
	}
	m.players[id] = nil
}

func (m *Match) BroadcastBoard() {
	for _, p := range m.players {
		if p != nil {
			p.Conn.WriteJSON(map[string]any{
				"type": "board",
				"data": TranslateBoard(m.board, p.id),
			})
		}
	}
}

func (m *Match) SendBoard(id byte) {
	p := m.players[id]
	if p != nil {
		p.Conn.WriteJSON(map[string]any{
			"type": "board",
			"data": TranslateBoard(m.board, p.id),
		})
	}
}

func (m *Match) SendError(id byte, err error) {
	p := m.players[id]
	log.Println("sendind error")
	if p != nil {
		p.Conn.WriteJSON(map[string]any{
			"type": "error",
			"data": err.Error(),
		})
	}
}

func (m *Match) UpdatePlayersBoard() {
	for _, p := range m.players {
		if p != nil && len(p.UpdateBoardChan) == 0 {
			p.UpdateBoardChan <- true
		}
	}
}

var match Match = Match{board: b}

func Connect(ctx *gin.Context) {
	// TODO: rework this to find a match in a list of matches (only when the project be about to be finished)
	if match.board == nil {
		match.board = b
	}
	if match.players[0] != nil && match.players[1] != nil {
		ctx.Status(400)
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: allowAnyOrigin}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		ctx.Status(500)
		return
	}
	defer conn.Close()
	p := match.NewPlayer(conn)
	if p == nil {
		ctx.Status(400)
		return
	}
	p.UpdateBoardChan <- true
	log.Println("Player", p.id+1, "connected")
gameloop:
	for {
		select {
		case r := <-p.ResultsChan:
			if r == int(p.id) {
				p.Conn.WriteJSON(map[string]any{
					"type": "result",
					"data": "You win!",
				})
			} else {
				p.Conn.WriteJSON(map[string]any{
					"type": "result",
					"data": "You lose",
				})
			}
			match.DisconnectPlayer(p.id)
			break gameloop
		case <-p.UpdateBoardChan:
			match.SendBoard(p.id)
		}
	actionloop:
		for {
			if match.board.PlayerTurn == p.id {
				if w := <-b.WaitingActionChan; w != -1 {
					p.ResultsChan <- w
					match.players[1-w].ResultsChan <- 1 - w
					break actionloop
				}
				mt, msg, err := conn.ReadMessage()
				if err != nil || mt != websocket.TextMessage {
					break gameloop
				}
				act := &cards.Action{}
				json.Unmarshal(msg, act)
				match.board.ActionChan <- act
				err = <-b.ActionEndChan
				if err != nil {
					match.SendError(p.id, err)
				}
				match.UpdatePlayersBoard()
			}
			break actionloop
		}
	}
	match.DisconnectPlayer(p.id)
	b.WaitingActionChan <- -1
	match.UpdatePlayersBoard()
	log.Println("Player", p.id+1, "disconnected")
}

func allowAnyOrigin(r *http.Request) bool {
	return true
}
