package main

import (
	"github.com/gorilla/websocket"
)

// clientはチャットを行っている一人のユーザーを表す
type client struct {
	//	このクライアントのためのwebsocket. websocketは、サーバーとクライアント双方向から通信ができるプロトコル
	socket *websocket.Conn
	//	メッセージが送られるチャネル
	send chan []byte
	// クライアントが参加しているルーム
	room *room
}

// websocketからデータを読み込む*clientのメソッド
func (c *client) read() {
	for {
		// errがなければメッセージをroomのforwardチャネルに送る
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

// websocketへ書き込む*clientのメソッド
func (c *client) write() {
	for msg := range c.send {
		// sendチャネルからメッセージを受け取り、websocketのWriteMessageメソッドを使って書き出す。
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
