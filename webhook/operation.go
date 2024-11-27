package webhook

import (
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/types"

	"github.com/scottd018/rosa-windows-overcommit-webhook/resources"
)

// operation represents an instance of each webhook operation that comes in.  This becomes one object created
// per webhook request so that operations may be parallelized for multiple webhook calls.
type operation struct {
	request  *request
	response *response
	object   resources.WindowsInstanceValidator
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

	// get the extractor used for extracting the instance
	// TODO: correct logic if we ever need to account for both virtual machines and virtual machine instances.  For now
	// we are only counting virtual machine instances.
	var validator resources.WindowsInstanceValidator
	switch req.admissionRequest.Kind.Kind {
	// case resources.VirtualMachineType:
	// 	validator = resources.NewVirtualMachine()
	case resources.VirtualMachineInstanceType:
		validator = resources.NewVirtualMachineInstance()
	default:
		return nil, fmt.Errorf(
			"unsupported kind [%s]; only [%+v] supported",
			req.admissionRequest.Kind.Kind,
			resources.SupportedResourceTypes(),
		)
	}

	// extract the instance
	instance, err := validator.Extract(req.admissionRequest)
	if err != nil {
		return nil, fmt.Errorf("failed extracting object from request; %w", err)
	}

	// set some values
	op.response.uid = req.admissionRequest.UID
	op.request = req
	op.response.review = req.admissionReview
	op.object = instance

	return op, nil
}
