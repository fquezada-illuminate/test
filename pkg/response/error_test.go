package response

import (
	"net/http"
	"testing"
)

func TestNewErrorResponse(t *testing.T) {
	inputSC := http.StatusBadRequest
	inputMsg := "Bad request"

	output := NewErrorResponse(inputSC, inputMsg)

	expected := ErrorResponse{
		Error: errorObj{
			Code:    inputSC,
			Message: inputMsg,
		},
	}

	if output != expected {
		t.Errorf("expected %v, got %v", expected, output)
	}
}
