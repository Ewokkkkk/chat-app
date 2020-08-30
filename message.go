package main

// messageは１つのメッセージを表す
type message struct {
	Name      string // ユーザー名
	Message   string // メッセージ
	When      string // タイムスタンプ
	AvatarURL string // アバターのURL
}
