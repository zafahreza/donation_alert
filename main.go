// Command quickstart generates an audio file with the content "Hello, World!".
package main

import (
	"donations_alert/Service"
	"donations_alert/app"
	"donations_alert/controller"
	"donations_alert/model"
	"github.com/caddyserver/certmagic"
	"log"
	"net/http"
)

func main() {

	router := http.NewServeMux()
	donationChan := make(chan model.ToClient)
	brokerReader := app.BrokerConsumer()

	alertService := Service.NewAlertService(brokerReader, donationChan)
	alertController := controller.NewAlertController(alertService, donationChan)

	router.HandleFunc("/ws/test", alertController.GetAlertController)

	err := certmagic.HTTPS([]string{"localhost"}, router)
	if err != nil {
		log.Fatalln(err)
	}

}
