package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

var upgrader = websocket.Upgrader{} // use default options

var clients = make(map[*websocket.Conn]chan struct{})
var clientsMutex sync.Mutex

func (app *application) startPostgresFetcher(stopChan chan struct{}) {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			log.Println("Postgres fetcher stopping")
			return
		case <-ticker.C:
			devices, err := app.Models.Device.GetAll()
			if err != nil {
				app.logger.Error(err.Error())
				return
			}

			if len(devices) > 0 {
				jsonData, err := json.Marshal(devices)
				if err == nil {
					broadcastToClients(jsonData)
				}
			}
		default:
		}
	}
}

func (app *application) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // allow all origins

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Close error:", err)
		}
	}(conn)

	// Create a stop channel for the client
	stopChan := make(chan struct{})

	// Register the client in the global clients map with the stop channel
	clientsMutex.Lock()
	clients[conn] = stopChan
	clientsMutex.Unlock()

	log.Println("New client connected")

	// Start the goroutine, passing the stop channel
	go app.startPostgresFetcher(stopChan)

	// Listen for client messages and handle disconnection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	// When the client disconnects, close the stop channel to notify the goroutine
	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
	log.Println("Client disconnected")

	// Close the stop channel to notify the goroutine to stop
	close(stopChan)
}

func broadcastToClients(message []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("WebSocket error:", err)
			err := client.Close()
			if err != nil {
				log.Println("Close error:", err)
				return
			}
			delete(clients, client)
		}
	}
}
