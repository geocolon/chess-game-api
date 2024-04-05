// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// var clients = make(map[*websocket.Conn]bool) // Connected clients
// var broadcast = make(chan Message)           // Broadcast channel

// // Message struct for WebSocket messages
// type Message struct {
// 	Username string `json:"username"`
// 	Message  string `json:"message"`
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	// Upgrade initial GET request to a WebSocket
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer ws.Close()

// 	// Register new client
// 	clients[ws] = true

// 	for {
// 		var msg Message
// 		// Read message from client
// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			delete(clients, ws)
// 			break
// 		}
// 		// Send received message to broadcast channel
// 		broadcast <- msg
// 	}
// }

// func handleMessages() {

// 	for {
// 		// Get next message from broadcast channel
// 		msg := <-broadcast
// 		// Send message to every connected client
// 		for client := range clients {
// 			err := client.WriteJSON(msg)
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }

// func messenger() {
// 	// Configure WebSocket route
// 	http.HandleFunc("/ws", handleConnections)
// 	// Start listening for incoming chat messages
// 	go handleMessages()
// 	// Start server
// 	log.Println("Server started")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
