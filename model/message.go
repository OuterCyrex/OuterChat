package model

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FromId   uint
	TargetId uint `json:"targetId"`
	Type     int  `gorm:"type:tinyint" json:"type"`
	Media    int
	Content  string `json:"content"`
	Pic      string
	Url      string
	Desc     string
	Amount   int
}

type History struct {
	From []Message
	To   []Message
}

// Type
const (
	private = iota + 1
)

// Media
const (
	text = iota + 1
	emoji
	picture
	sound
)

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

var clientMap = make(map[uint]*Node)

var rwMutex sync.RWMutex

func Chat(writer http.ResponseWriter, request http.Request) {
	query := request.URL.Query()
	userId, _ := strconv.Atoi(query.Get("userId"))
	// token := request.Header.Get("token")

	tokenValid := true

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return tokenValid
		},
	}).Upgrade(writer, &request, nil)
	if err != nil {
		fmt.Printf("Http Upgrade to Websocket Failed: %v", err)
		return
	}

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	rwMutex.Lock()
	clientMap[uint(userId)] = node
	rwMutex.Unlock()

	go sendProc(node)
	go recvProc(userId, node)

	SendMsg(uint(userId), []byte("欢迎来到聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Printf("Send Message Failed: %v", err)
				return
			}
		}
	}
}

func recvProc(userId int, node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("Receive Message Failed: %v", err)
			return
		}
		fmt.Println("[Receive-Data]:", string(data))
		dispatch(userId, data)
	}
}

func dispatch(userId int, data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("unmarshal failed: %v", err)
		return
	}
	msg.FromId = uint(userId)
	DB.Create(&msg)
	switch msg.Type {
	case 1:
		SendMsg(msg.TargetId, data)
	default:
		return
	}
}

func SendMsg(targetId uint, msg []byte) {
	rwMutex.RLock()
	node, ok := clientMap[targetId]
	rwMutex.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func GetHistory(userId uint, targetId uint) (History, error) {
	var history History
	err := DB.Model(&Message{}).Where("from_id = ? and target_id = ?", userId, targetId).Find(&history.From).Error
	if err != nil {
		return History{}, err
	}
	err = DB.Model(&Message{}).Where("target_id = ? and from_id = ?", userId, targetId).Find(&history.To).Error
	if err != nil {
		return History{}, err
	}
	return history, nil
}
