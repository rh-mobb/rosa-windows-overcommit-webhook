package vm

import (
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	kubevirtcorev1 "kubevirt.io/api/core/v1"
)

const (
	VirtualMachineType = "VirtualMachine"
)

// virtualMachine represents a VirtualMachine object.
type virtualMachine struct{}

// NewVirtualMachine returns a new virtualMachine object.
func NewVirtualMachine() *virtualMachine {
	return &virtualMachine{}
}

// Extract extracts a VirtualMachineInstance object from a VirtualMachine object.
func (vm *virtualMachine) Extract(admissionRequest *admissionv1.AdmissionRequest) (*kubevirtcorev1.VirtualMachineInstance, error) {
	virtualMachine := &kubevirtcorev1.VirtualMachine{}
	if err := json.Unmarshal(admissionRequest.Object.Raw, virtualMachine); err != nil {
		return nil, fmt.Errorf("failed to decode virtual machine object; %w", err)
	}

	// derive the spec from the virtual machine instance spec
	return VirtualMachineInstanceFromVirtualMachine(virtualMachine), nil
}

// Type simply returns the type.
func (vm *virtualMachine) Type() string {
	return VirtualMachineType
}
