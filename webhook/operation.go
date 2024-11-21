package webhook

import (
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/types"
)

// operation represents an instance of each webhook operation that comes in.  This becomes one object created
// per webhook request so that operations may be parallelized for multiple webhook calls.
type operation struct {
	request  *request
	response *response
}

// NewOperation return a new instance of an operation object.
func NewOperation(w http.ResponseWriter, r *http.Request) (*operation, error) {
	// create the base operation object
	op := &operation{
		response: &response{allowed: true, uid: types.UID(""), writer: w},
	}

	// create the request object
	req, err := newRequest(r)
	if err != nil {
		return op, fmt.Errorf("unable to create request object; %w", err)
	}

	// set some values
	op.response.uid = req.admissionRequest.UID
	op.request = req
	op.response.review = req.admissionReview

	return op, nil
}
