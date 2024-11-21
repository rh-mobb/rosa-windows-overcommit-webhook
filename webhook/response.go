package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// response is an object which represents the response from an individual webhook operation.
type response struct {
	allowed bool
	uid     types.UID
	writer  http.ResponseWriter
	review  *admissionv1.AdmissionReview
}

// send sends a response.
func (r *response) send(message string, logToStdout bool) {
	if logToStdout {
		log.Printf("returning with message: [%s]", message)
	}

	r.review.Response = &admissionv1.AdmissionResponse{
		Allowed: r.allowed,
		UID:     r.uid,
		Result: &metav1.Status{
			Message: message,
			Code:    http.StatusOK,
		},
	}

	responseBody, err := json.Marshal(r.review)
	if err != nil {
		log.Println(fmt.Errorf("unable to marshal response body; %w", err))
	}

	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.Write(responseBody)
}
