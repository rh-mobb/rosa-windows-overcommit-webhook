package resources

import (
	corev1 "kubevirt.io/api/core/v1"
)

// VirtualMachines is an object which represents a list of virtual machines.
type VirtualMachines []corev1.VirtualMachine

// VirtualMachinesFilter represents a filter based on a set of key value inputs that are used to filter nodes.
type VirtualMachinesFilter struct {
	VirtualMachineInstances VirtualMachineInstances
}

// Filter filters a Store object and returns a new store with only filtered virtual machine instances.  In the
// instance of this webhook, we only want virtual machine instances that are running a windows operating system.
func (vms VirtualMachines) Filter(filter *VirtualMachinesFilter) VirtualMachineInstances {
	filtered := VirtualMachineInstances{}

	for i := 0; i < len(vms); i++ {
		var vm virtualMachine = virtualMachine(vms[i])

		if vm.isWindows() {
			filtered = append(filtered, corev1.VirtualMachineInstance(*vm.VirtualMachineInstance()))
		}
	}

	return filtered
}
