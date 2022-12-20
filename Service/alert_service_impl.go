package Service

import (
	"donations_alert/model"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
)

type AlertServiceImpl struct {
	Reader       *kafka.Reader
	DonationChan chan model.ToClient
}

func NewAlertService(reader *kafka.Reader, donationChan chan model.ToClient) AlertService {
	return &AlertServiceImpl{
		Reader:       reader,
		DonationChan: donationChan,
	}
}

func (service *AlertServiceImpl) GetAlertService(w http.ResponseWriter, r *http.Request) {
	for {
		err := service.Reader.SetOffset(kafka.LastOffset)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("sedang membaca data dari kafka")
		message, err := service.Reader.ReadMessage(r.Context())
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("get message from kafka")

		var donation model.Donations
		err = json.Unmarshal(message.Value, &donation)
		if err != nil {
			log.Fatalln(err)
		}

		if donation.IsAnonymous == true {
			donation.From = "Seseorang"
		}

		fmt.Println("mengirim request ke google api")

		byteAudio, object := Synthesize(r.Context(), donation)
		audioUrl := fmt.Sprintf("https://storage.googleapis.com/donation_alert/%s", object)

		fmt.Println("mendapatkan response dari google api")
		toClient := model.ToClient{
			Donation: donation,
			Audio:    byteAudio,
			AudioUrl: audioUrl,
		}

		fmt.Println("mengirim data ke channel")
		service.DonationChan <- toClient
	}
}
