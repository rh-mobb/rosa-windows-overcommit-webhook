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

// IsWindows determines if a virtual machine instance object is a windows instance or not.
func (vmi virtualMachineInstance) IsWindows() bool {
	if vmi.HasSysprepVolume() {
		return true
	}

	if vmi.HasWindowsDriverDiskVolume() {
		return true
	}

	return false
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

// Type simply returns the type.
func (vmi virtualMachineInstance) Type() string {
	return VirtualMachineInstanceType
}

// HasSysprepVolume returns if the virtualmachineinstance has a sysprep volume or not.  Sysprep volumes are exclusive
// to windows machines.
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error that includes
// sysprep volumes for linux machines.  This is guaranteed to work when using out of the box OpenShift templates.
func (vmi virtualMachineInstance) HasSysprepVolume() bool {
	for _, volume := range vmi.Spec.Volumes {
		if volume.Sysprep != nil {
			return true
		}
	}

	return false
}

// HasWindowsDriverDiskVolume returns if the virtualmachineinstance has a windows driver volume or not.  Windows driver
// volumes are used for adding windows drivers to windows machines, however it is not restrictive that this only may
// be included on windows machines (although unlikely).
// WARN: it should be noted that users who deploy their instances via YAML may have a copy/paste error that includes
// sysprep volumes for linux machines.  This is guaranteed to work when using out of the box OpenShift templates.
func (vmi virtualMachineInstance) HasWindowsDriverDiskVolume() bool {
	for _, volume := range vmi.Spec.Volumes {
		if volume.DataVolume.Name == "windows-drivers-disk" {
			return true
		}
	}

	return false
}
