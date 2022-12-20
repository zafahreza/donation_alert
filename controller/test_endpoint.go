package controller

import (
	"donations_alert/app"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Messages struct {
	Greating string `json:"greating"`
}

func Endpoint(w http.ResponseWriter, r *http.Request) {
	ws := app.NewWebSocket()

	ws.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	wsConn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer wsConn.Close()

	for {
		var req Messages
		err = wsConn.ReadJSON(&req)
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Printf("New Message : %s\n", req.Greating)

		req.Greating = "hello back from server"

		newJson, err := json.Marshal(req)
		if err != nil {
			log.Fatal(err)
		}

		err = wsConn.WriteMessage(websocket.TextMessage, newJson)
		if err != nil {
			log.Println(err)
			break
		}

	}
}
