// from github.com/waterloop/arduino-challenge-dashboard

package main

import (
	"github.com/tarm/serial"
  "log"
  "time"
  "net/http"
  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func main() {
  wsconns := make(map[*websocket.Conn]bool, 0)

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    // Upgrade the http connection to a WebSocket connetion
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil { log.Println("err: upgrade:", err); return }
    defer c.Close()

    // Add it to the set of connections
    log.Printf("Connection added")
    wsconns[c] = true
    defer delete(wsconns, c)

    // Wait for messages and deal with them
    for {
      _, message, err := c.ReadMessage()
      if err != nil { log.Println("err: read:", err); break }
      log.Printf("ws-recv: path{%s} message{%s}", r.URL.Path, string(message))
      // TODO
    }
  })

  log.Println("About to listen on :3333")
  go func () { log.Fatal(http.ListenAndServe(":3333", nil)) }()

  // Open the UART connection.
	s, err := serial.OpenPort(&serial.Config{Name: "/dev/ttyS3", Baud: 115200})
  if err != nil { log.Fatal(err) }
	time.Sleep(1 * time.Second) // Sometimes Arduinos need a bit of time to reset after connecting.
  s.Flush()

  buf := make([]byte, 256)

  for {
    // Read some bytes (note that UART may fragment packets, so a forloop reading into buf[bytesread:] might be good here.)
    bytesread, err := s.Read(buf)
    if err != nil { log.Fatal(err) }

    log.Printf("uart-recv: message{%s}", string(buf[:bytesread]))

    // Send the string to each websocket connection
    for c, _ := range wsconns {
      c.WriteMessage(websocket.TextMessage, buf[:bytesread])
    }
  }
}