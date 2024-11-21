package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	kubevirtcorev1 "kubevirt.io/api/core/v1"

	"github.com/scottd018/rosa-windows-overcommit-webhook/vm"
)

type request struct {
	admissionRequest       *admissionv1.AdmissionRequest
	admissionReview        *admissionv1.AdmissionReview
	virtualMachineInstance *kubevirtcorev1.VirtualMachineInstance
}

// newRequest creates a new request object.
func newRequest(r *http.Request) (*request, error) {
	// read in the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body; %w", err)
	}
	defer r.Body.Close()

	log.Println("unmarshaling request in admission review object")
	var admissionReview admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		return nil, fmt.Errorf("failed to unmarshal admission review; %w", err)
	}

	// we only care about create operations for now
	// TODO: handle update in case create was bypassed somehow
	admissionRequest := admissionReview.Request
	if admissionRequest.Operation != admissionv1.Create {
		return nil, fmt.Errorf("unsupported operation [%s]; only create supported", admissionRequest.Operation)
	}

	// get the extractor used for extracting the instance
	var instanceExtractor vm.VirtualMachineInstanceExtractor
	switch admissionRequest.Kind.Kind {
	case vm.VirtualMachineType:
		instanceExtractor = vm.NewVirtualMachine()
	case vm.VirtualMachineInstanceType:
		instanceExtractor = vm.NewVirtualMachineInstance()
	default:
		return nil, fmt.Errorf(
			"unsupported kind [%s]; only [%+v] supported",
			admissionRequest.Kind.Kind,
			vm.SupportedTypes(),
		)
	}

	// extract the instance
	log.Printf("extracting object from request [%s]", instanceExtractor.Type())
	instance, err := instanceExtractor.Extract(admissionRequest)
	if err != nil {
		return nil, fmt.Errorf("failed extracting object from request; %w", err)
	}

	return &request{
		virtualMachineInstance: instance,
		admissionRequest:       admissionRequest,
		admissionReview:        &admissionReview,
	}, nil
}
