package resources

import (
	admissionv1 "k8s.io/api/admission/v1"
)

type WindowsValidationResult struct {
	NeedsValidation bool
	Reason          string
}

// WindowsInstanceValidator is an interface that represents an object containing all methods required to
// validate a windows instance.
type WindowsInstanceValidator interface {
	Extract(*admissionv1.AdmissionRequest) (*virtualMachineInstance, error)
	SumCPU() int
	NeedsValidation() *WindowsValidationResult

	GetName() string
	GetNamespace() string
}

// SupportedResourceTypes returns the supported resources for this webhook.
func SupportedResourceTypes() []string {
	return []string{
		VirtualMachineType,
		VirtualMachineInstanceType,
	}
}
