package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rawdaGastan/urls_checker/internal"
)

func main() {
	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	http.HandleFunc("/site/", func(w http.ResponseWriter, r *http.Request) {

		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		// Continuosly read and write message
		for {
			_, website, err := conn.ReadMessage()
			if err != nil {
				log.Println("read failed:", err)
				break
			}

			service := internal.NewCheckerService(100)
			service.AddSite(string(website))
			service.AddSocket(conn)
			service.Start()
		}
	})

	fmt.Println("server is running at", 3000)
	http.ListenAndServe(":3000", nil)
}
