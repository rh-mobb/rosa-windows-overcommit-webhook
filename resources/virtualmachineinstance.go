package resources

import (
	"encoding/json"
	"fmt"

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
func (vmi virtualMachineInstance) NeedsValidation() bool {
	// if we have no owner references, use windows logic to determine if we need validation
	if len(vmi.GetOwnerReferences()) == 0 {
		return vmi.isWindows()
	}

	// we do not need to validate an instance already owned by a virtual machine
	for _, ref := range vmi.GetOwnerReferences() {
		if ref.Name == VirtualMachineType {
			return false
		}
	}

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
func (vmi virtualMachineInstance) isWindows() bool {
	for _, hasWindowsIdentifier := range []func() bool{
		vmi.hasSysprepVolume,
		vmi.hasWindowsDriverDiskVolume,
		vmi.hasHyperV,
	} {
		if hasWindowsIdentifier() {
			return true
		}
	}

	return false
}

// hasSysprepVolume returns if the virtualmachineinstance has a sysprep volume or not.  Sysprep volumes are exclusive
// to windows machines.
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error that includes
// sysprep volumes for linux machines. This is guaranteed to work when using out of the box OpenShift templates, but
// may not work with use created templates.
func (vmi virtualMachineInstance) hasSysprepVolume() bool {
	for _, volume := range vmi.Spec.Volumes {
		if volume.Sysprep != nil {
			return true
		}
	}

	return false
}

// hasWindowsDriverDiskVolume returns if the virtualmachineinstance has a windows driver volume or not.  Windows driver
// volumes are used for adding windows drivers to windows machines, however it is not restrictive that this only may
// be included on windows machines (although unlikely).
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error that includes
// sysprep volumes for linux machines. This is guaranteed to work when using out of the box OpenShift templates, but
// may not work with use created templates.
func (vmi virtualMachineInstance) hasWindowsDriverDiskVolume() bool {
	for _, volume := range vmi.Spec.Volumes {
		if volume.DataVolume == nil {
			continue
		}

		if volume.DataVolume.Name == "windows-drivers-disk" {
			return true
		}
	}

	return false
}

// hasHyperV returns if the virtualmachineinstance has Hyper-V settings set.
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error.
// This is guaranteed to work when using out of the box OpenShift templates, but may not work with use created
// templates.
func (vmi virtualMachineInstance) hasHyperV() bool {
	if vmi.Spec.Domain.Features == nil {
		return false
	}

	return vmi.Spec.Domain.Features.Hyperv != nil
}
