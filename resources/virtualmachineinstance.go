package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "kubevirt.io/api/core/v1"
)

const (
	VirtualMachineInstanceType = "VirtualMachineInstance"
)

// VirtualMachineInstance represents a VirtualMachineInstance object.
type virtualMachineInstance corev1.VirtualMachineInstance

// NewVirtualMachineInstance returns a new virtualMachineInstance object.
func NewVirtualMachineInstance() *virtualMachineInstance {
	return &virtualMachineInstance{}
}

// Extract extracts a VirtualMachineInstance object from a VirtualMachine object.
func (vmi virtualMachineInstance) Extract(admissionRequest *admissionv1.AdmissionRequest) (*virtualMachineInstance, error) {
	instance := &virtualMachineInstance{}
	if err := json.Unmarshal(admissionRequest.Object.Raw, &instance); err != nil {
		return nil, fmt.Errorf("failed to decode virtual machine instance object; %w", err)
	}

	return instance, nil
}

// NeedsValidation returns if a virtual machine instance object needs validation or not.
func (vmi virtualMachineInstance) NeedsValidation() *WindowsValidationResult {
	// TODO: correct logic if we ever need to account for both virtual machines and virtual machine instances.  For now
	// we are only counting virtual machine instances.
	// // if we have owner references, see if we are owned by a virtual machine
	// if len(vmi.GetOwnerReferences()) > 0 {
	// 	for _, ref := range vmi.GetOwnerReferences() {
	// 		// we only want to validate virtual machine instances that do not have a windows preference.  this is
	// 		// because provisioning from a preference seems to have a specialized workflow that is difficult
	// 		// to determine for windows.
	// 		// TODO: this logic is likely to need adjusted.
	// 		if ref.Name == VirtualMachineType {
	// 			return &WindowsValidationResult{Reason: "virtual machine instance owned by virtual machine object"}
	// 		}
	// 	}
	// }

	// finally use the windows logic to determine if we need validation
	return vmi.isWindows()
}

// SumCPU sums up the value of all CPUs for the virtual machine instance.
func (vmi virtualMachineInstance) SumCPU() int {
	sockets, cores, threads := 1, 1, 1

	// this simply defaults to 1 * 1 * 1
	if vmi.Spec.Domain.CPU == nil {
		return sockets * cores * threads
	}

	// according to kubevirt docs, vcpu is determined by the value of sockets * cores * threads
	// see https://kubevirt.io/user-guide/compute/dedicated_cpu_resources/#requesting-dedicated-cpu-resources
	sockets = int(vmi.Spec.Domain.CPU.Sockets)
	if sockets == 0 {
		sockets = 1
	}

	cores = int(vmi.Spec.Domain.CPU.Cores)
	if cores == 0 {
		cores = 1
	}

	threads = int(vmi.Spec.Domain.CPU.Threads)
	if threads == 0 {
		threads = 1
	}

	return sockets * cores * threads
}

// isWindows determines if a virtual machine instance object is a windows instance or not.
func (vmi virtualMachineInstance) isWindows() *WindowsValidationResult {
	for _, hasWindowsIdentifier := range []func() *WindowsValidationResult{
		vmi.hasSysprepVolume,
		vmi.hasWindowsDriverDiskVolume,
		vmi.hasHyperV,
		vmi.hasWindowsPreference,
	} {
		result := hasWindowsIdentifier()

		if result.NeedsValidation {
			return result
		}
	}

	return &WindowsValidationResult{Reason: "no validation required"}
}

// hasSysprepVolume returns if the virtualmachineinstance has a sysprep volume or not.  Sysprep volumes are exclusive
// to windows machines.
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error that includes
// sysprep volumes for linux machines. This is guaranteed to work when using out of the box OpenShift templates, but
// may not work with use created templates.
func (vmi virtualMachineInstance) hasSysprepVolume() *WindowsValidationResult {
	for _, volume := range vmi.Spec.Volumes {
		if volume.Sysprep != nil {
			return &WindowsValidationResult{NeedsValidation: true, Reason: "has sysprep volume"}
		}
	}

	return &WindowsValidationResult{Reason: "has no sysprep volume"}
}

// hasWindowsDriverDiskVolume returns if the virtualmachineinstance has a windows driver volume or not.  Windows driver
// volumes are used for adding windows drivers to windows machines, however it is not restrictive that this only may
// be included on windows machines (although unlikely).
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error that includes
// sysprep volumes for linux machines. This is guaranteed to work when using out of the box OpenShift templates, but
// may not work with use created templates.
func (vmi virtualMachineInstance) hasWindowsDriverDiskVolume() *WindowsValidationResult {
	for _, volume := range vmi.Spec.Volumes {
		if volume.DataVolume == nil {
			continue
		}

		if volume.DataVolume.Name == "windows-drivers-disk" {
			return &WindowsValidationResult{NeedsValidation: true, Reason: "has windows-driver-disk-volume"}
		}
	}

	return &WindowsValidationResult{Reason: "has no windows-driver-disk-volume"}
}

// hasHyperV returns if the virtualmachineinstance has Hyper-V settings set.
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error.
// This is guaranteed to work when using out of the box OpenShift templates, but may not work with use created
// templates.
func (vmi virtualMachineInstance) hasHyperV() *WindowsValidationResult {
	if vmi.Spec.Domain.Features == nil {
		return &WindowsValidationResult{Reason: "has nil features"}
	}

	if vmi.Spec.Domain.Features.Hyperv != nil {
		return &WindowsValidationResult{NeedsValidation: true, Reason: "has hyper-v features"}
	}

	return &WindowsValidationResult{Reason: "has no hyper-v features"}
}

// hasWindowsPreference returns if the virtualmachineinstance has a windows preference annotation.
// WARN: it should be noted that this annotation is created when provisioning from instance type.  It is entirely
// possible that users can select their own instance type and bypass this check.
func (vmi virtualMachineInstance) hasWindowsPreference() *WindowsValidationResult {
	annotations := vmi.GetAnnotations()

	if len(annotations) == 0 {
		return &WindowsValidationResult{Reason: "has no annotations"}
	}

	if annotations["vm.kubevirt.io/os"] == "windows" {
		return &WindowsValidationResult{
			NeedsValidation: true,
			Reason:          "has 'vm.kubevirt.io/os' windows annotation",
		}
	}

	if strings.HasPrefix(annotations["kubevirt.io/cluster-preference-name"], "windows") {
		return &WindowsValidationResult{
			NeedsValidation: true,
			Reason:          "has 'kubevirt.io/cluster-preference-name' windows annotation",
		}
	}

	return &WindowsValidationResult{Reason: "has no windows preference"}
}
