package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// 	他のクライアントに転送するためのメッセージを保存するチャネル
	forward chan []byte
	// 	チャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	//	チャットルームから退室しようとしているクライアントのためのチャネル
	leave chan *client
	//	在室している全てのクライアントが保持される
	clients map[*client]bool
}

// チャットルームを生成して返す
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	// 無限ループ（強制終了されるまで実行し続ける）
	for {
		// いずれかのチャネルにメッセージが届くと、それぞれの処理を行う
		select {
		case client := <-r.join:
			//	参加
			r.clients[client] = true
		case client := <-r.leave:
			//	退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			//	全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					//	メッセージ送信
				default:
					//	送信失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
