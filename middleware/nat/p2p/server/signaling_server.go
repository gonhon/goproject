// signaling-server/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Type      string          `json:"type"`
	Room      string          `json:"room"`
	SDP       json.RawMessage `json:"sdp,omitempty"`
	Candidate json.RawMessage `json:"candidate,omitempty"`
}

type Room struct {
	clients []*websocket.Conn
}

var rooms = make(map[string]*Room)
var roomsMutex sync.Mutex

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Signaling server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			removeClient(conn)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		switch message.Type {
		case "join":
			handleJoin(conn, message.Room)
		case "offer", "answer", "candidate":
			broadcast(conn, message.Room, msg)
		}
	}
}

func handleJoin(conn *websocket.Conn, roomID string) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	if _, exists := rooms[roomID]; !exists {
		rooms[roomID] = &Room{clients: make([]*websocket.Conn, 0)}
	}

	rooms[roomID].clients = append(rooms[roomID].clients, conn)
}

func broadcast(sender *websocket.Conn, roomID string, message []byte) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		return
	}

	for _, client := range room.clients {
		if client != sender {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Write error:", err)
				removeClient(client)
			}
		}
	}
}

func removeClient(conn *websocket.Conn) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	for roomID, room := range rooms {
		for i, client := range room.clients {
			if client == conn {
				room.clients = append(room.clients[:i], room.clients[i+1:]...)
				if len(room.clients) == 0 {
					delete(rooms, roomID)
				}
				return
			}
		}
	}
}
