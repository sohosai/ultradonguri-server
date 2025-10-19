package responses

import (
	"net/http"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewErrorResponseAndHTTPStatus(err entities.AppError) (*ErrorResponse, int) {
	message := err.Message
	code := err.Kind
	var status int

	switch err.Kind {
	case entities.NoConnectionToOBS:
		status = http.StatusInternalServerError
	case entities.NoConnectionToViewer:
		status = http.StatusInternalServerError
	case entities.InvalidPerformancesJson:
		status = http.StatusInternalServerError
	case entities.InvalidFormat:
		status = http.StatusBadRequest
	case entities.CannotConversion:
		status = http.StatusBadRequest
	case entities.CannotForceMute:
		status = http.StatusBadRequest
	case entities.CannotChangeState:
		status = http.StatusInternalServerError
	}

	return &ErrorResponse{Message: message, Code: code}, status
}
