package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// clientはチャットを行っている一人のユーザーを表す
type client struct {
	//	このクライアントのためのwebsocket. websocketは、サーバーとクライアント双方向から通信ができるプロトコル
	socket *websocket.Conn
	//	メッセージが送られるチャネル
	send chan *message
	// クライアントが参加しているルーム
	room *room
	// ユーザーに関する情報を保持する
	userData map[string]interface{}
}

// websocketからデータを読み込む*clientのメソッド
func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now().Format("15:04")
			msg.Name = c.userData["name"].(string)
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
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
