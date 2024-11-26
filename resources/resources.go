package resources

import (
	admissionv1 "k8s.io/api/admission/v1"
)

// WindowsInstanceValidator is an interface that represents an object containing all methods required to
// validate a windows instance.
type WindowsInstanceValidator interface {
	Extract(*admissionv1.AdmissionRequest) (*virtualMachineInstance, error)
	Type() string
	SumCPU() int
	NeedsValidation() bool

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
