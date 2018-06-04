// from github.com/waterloop/arduino-challenge-dashboard

package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))

	messages := make(chan string)

	http.HandleFunc("/in", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("HTTP /in")
		messages <- "i"
	})
	http.HandleFunc("/out", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("HTTP /out")
		messages <- "o"
	})

	wsconns := make(map[*websocket.Conn]bool, 0)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the http connection to a WebSocket connetion
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// log.Println("err: upgrade:", err)
			return
		}
		defer c.Close()

		// Add it to the set of connections
		// log.Printf("Connection added")
		wsconns[c] = true
		defer delete(wsconns, c)

		// Wait for messages and deal with them
		for {
			_, message, err := c.ReadMessage()
			// Usually means that the client has closed the connection
			if err != nil {
				// log.Println("err: read:", err)
				break
			}
			// log.Printf("ws-recv: path{%s} message{%s}", r.URL.Path, string(message))
			messages <- string(message)
		}
	})

	log.Println("About to listen on :3333")
	go func() { log.Fatal(http.ListenAndServe(":3333", nil)) }()

	for {
		message := <-messages

		// Send the string to each websocket connection
		for c, _ := range wsconns {
			c.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}
