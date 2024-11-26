package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtcorev1 "kubevirt.io/api/core/v1"
)

const (
	VirtualMachineType = "VirtualMachine"
)

// virtualMachine represents a VirtualMachine object.
type virtualMachine kubevirtcorev1.VirtualMachine

// NewVirtualMachine returns a new virtualMachine object.
func NewVirtualMachine() *virtualMachine {
	return &virtualMachine{}
}

// Extract extracts a VirtualMachineInstance object from a VirtualMachine object.
func (vm virtualMachine) Extract(admissionRequest *admissionv1.AdmissionRequest) (*virtualMachineInstance, error) {
	if err := json.Unmarshal(admissionRequest.Object.Raw, &vm); err != nil {
		return nil, fmt.Errorf("failed to decode virtual machine object; %w", err)
	}

	// derive the spec from the virtual machine instance spec
	return vm.VirtualMachineInstance(), nil
}

// NeedsValidation returns if a virtual machine object needs validation or not.
func (vm virtualMachine) NeedsValidation() bool {
	return vm.isWindows()
}

// SumCPU sums up the value of all CPUs for the virtual machine.
func (vm virtualMachine) SumCPU() int {
	return vm.VirtualMachineInstance().SumCPU()
}

// Type simply returns the type.
func (vm virtualMachine) Type() string {
	return VirtualMachineType
}

// VirtualMachineInstance returns the virtual machine instance object from the virtual machine template spec.
func (vm virtualMachine) VirtualMachineInstance() *virtualMachineInstance {
	return &virtualMachineInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       VirtualMachineInstanceType,
			APIVersion: fmt.Sprintf("%s/%s", vm.GroupVersionKind().Group, vm.APIVersion),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vm.Name,
			Namespace: vm.Namespace,
		},
		Spec: vm.Spec.Template.Spec,
	}
}

// isWindows determines if a virtual machine object is a windows instance or not.
func (vm virtualMachine) isWindows() bool {
	if vm.hasWindowsPreference() {
		return true
	}

	return vm.VirtualMachineInstance().isWindows()
}

// hasWindowsPreference returns if the virtualmachineinstance has a windows preference set.  This is for virtual machines
// that end up being provisioned by an instances type versus a template.  This appears to be the most reliable
// way to determine windows but it does not stop a user from mislabeling the instance type.
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error that includes
// sysprep volumes for linux machines.  This is guaranteed to work when using out of the box OpenShift templates.
func (vm virtualMachine) hasWindowsPreference() bool {
	if vm.Spec.Preference == nil {
		return false
	}

	return strings.HasPrefix(vm.Spec.Preference.Name, "windows")
}
