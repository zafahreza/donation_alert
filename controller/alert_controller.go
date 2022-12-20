package controller

import (
	"donations_alert/Service"
	"donations_alert/app"
	"donations_alert/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type AlertController struct {
	Service      Service.AlertService
	DonationChan chan model.ToClient
}

func NewAlertController(service Service.AlertService, donationChan chan model.ToClient) *AlertController {
	return &AlertController{Service: service, DonationChan: donationChan}
}

func (controller *AlertController) GetAlertController(w http.ResponseWriter, r *http.Request) {
	ws := app.NewWebSocket()

	ws.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	wsConn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer wsConn.Close()

	var req Messages
	err = wsConn.ReadJSON(&req)
	if err != nil {
		log.Println(err)

	}

	fmt.Printf("New Message : %s\n", req.Greating)

	go controller.Service.GetAlertService(w, r)

	for {
		newAlert := <-controller.DonationChan
		fmt.Println("mendapatkan data dari channel")
		newJson, err := json.Marshal(newAlert)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("mengirim data ke client")
		err = wsConn.WriteMessage(websocket.TextMessage, newJson)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("message sent")
	}
}
