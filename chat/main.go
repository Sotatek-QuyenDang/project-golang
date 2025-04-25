/* File: main.go */
package main

import (
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"os"
	"sync"
)

type UserInfo struct {
	Nickname string
	Room     string
}

var mu sync.Mutex
var chatLogFile = "chat.log"

func saveMessage(nickname, room, msg string) {
	mu.Lock()
	defer mu.Unlock()
	f, err := os.OpenFile(chatLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("❌ Failed to open log file:", err)
		return
	}
	defer f.Close()
	logLine := fmt.Sprintf("[%s] %s: %s\n", room, nickname, msg)
	if _, err := f.WriteString(logLine); err != nil {
		log.Println("❌ Failed to write to log file:", err)
	}
}

func main() {
	server := socketio.NewServer(nil)
	if server == nil {
		log.Fatal("❌ Failed to create socket.io server")
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext(UserInfo{})
		log.Printf("✅ Client connected - ID: %s, RemoteAddr: %s\n", s.ID(), s.RemoteAddr())
		return nil
	})

	server.OnEvent("/", "join", func(s socketio.Conn, data string) {
		log.Printf("📥 Received join request - Data: %s\n", data)
		var info UserInfo
		if err := json.Unmarshal([]byte(data), &info); err != nil {
			log.Printf("❌ Failed to unmarshal join data: %v - Raw data: %s\n", err, data)
			return
		}
		if info.Nickname == "" || info.Room == "" {
			log.Printf("❌ Invalid join data - Nickname: '%s', Room: '%s'\n", info.Nickname, info.Room)
			return
		}
		s.SetContext(info)
		s.Join(info.Room)
		log.Printf("🚪 Client %s joined room %s as %s\n", s.ID(), info.Room, info.Nickname)
		s.Emit("system", fmt.Sprintf("You joined room: %s as %s", info.Room, info.Nickname))
		server.BroadcastToRoom("/", info.Room, "system", fmt.Sprintf("%s joined the chat", info.Nickname))
	})

	server.OnEvent("/", "chat", func(s socketio.Conn, msg string) {
		info, ok := s.Context().(UserInfo)
		if !ok {
			log.Printf("❌ Invalid context for client %s\n", s.ID())
			return
		}
		if msg == "" {
			log.Printf("⚠️ Empty message received from %s\n", info.Nickname)
			return
		}
		log.Printf("💬 Message from %s in room %s: %s\n", info.Nickname, info.Room, msg)
		fullMsg := fmt.Sprintf("%s: %s", info.Nickname, msg)
		s.Emit("chat", fullMsg)
		server.BroadcastToRoom("/", info.Room, "chat", fullMsg)
		saveMessage(info.Nickname, info.Room, msg)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		info, ok := s.Context().(UserInfo)
		if !ok {
			log.Printf("❌ Client disconnected without context - ID: %s, Reason: %s\n", s.ID(), reason)
			return
		}
		log.Printf("👋 Client %s (%s) left room %s - Reason: %s\n", s.ID(), info.Nickname, info.Room, reason)
		server.BroadcastToRoom("/", info.Room, "system", fmt.Sprintf("%s left the chat", info.Nickname))
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal("❌ Failed to start socket.io server:", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))

	log.Printf("🚀 Server starting at http://localhost:8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("❌ HTTP server failed:", err)
	}
}
