package resources

import (
	admissionv1 "k8s.io/api/admission/v1"
	kubevirtcorev1 "kubevirt.io/api/core/v1"
)

type VirtualMachineInstanceExtractor interface {
	Extract(*admissionv1.AdmissionRequest) (*kubevirtcorev1.VirtualMachineInstance, error)
	Type() string
}

// SupportedResourceTypes returns the supported resources for this webhook.
func SupportedResourceTypes() []string {
	return []string{
		VirtualMachineType,
		VirtualMachineInstanceType,
	}
}
