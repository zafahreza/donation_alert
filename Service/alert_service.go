package Service

import (
	"net/http"
)

type AlertService interface {
	GetAlertService(w http.ResponseWriter, r *http.Request)
}
