package resources

import (
	"fmt"

	corev1 "kubevirt.io/api/core/v1"
)

// VirtualMachineInstances is an object which represents a list of virtual machine instances.
type VirtualMachineInstances []corev1.VirtualMachineInstance

// VirtualMachineInstancesFilter represents a filter based on a set of key value inputs that are used to filter nodes.
type VirtualMachineInstancesFilter struct{}

// Filter filters a Store object and returns a new store with only filtered virtual machine instances.  In the
// instance of this webhook, we only want virtual machine instances that are running a windows operating system.
func (instances VirtualMachineInstances) Filter(filter *VirtualMachineInstancesFilter) VirtualMachineInstances {
	filtered := VirtualMachineInstances{}

	for i := 0; i < len(instances); i++ {
		var instance virtualMachineInstance = virtualMachineInstance(instances[i])

		if instance.HasSysprepVolume() {
			filtered = append(filtered, corev1.VirtualMachineInstance(instance))
		}

		if instance.HasWindowsDriverDiskVolume() {
			filtered = append(filtered, corev1.VirtualMachineInstance(instance))
		}
	}

	return filtered
}

// Unique returns a set of unique virtual machine instances, designated by the name and namespace.
func (instances VirtualMachineInstances) Unique() VirtualMachineInstances {
	found := map[string]bool{}

	var unique VirtualMachineInstances

	for i := 0; i < len(instances); i++ {
		key := fmt.Sprintf("%s/%s", instances[i].Name, instances[i].Namespace)

		if found[key] {
			continue
		}

		unique = append(unique, instances[i])
		found[key] = true
	}

	return unique
}

// SumCPU sums up the value of all CPUs in the store.
func (instances VirtualMachineInstances) SumCPU() int {
	var sum int

	if len(instances) == 0 {
		return 0
	}

	for vm := 0; vm < len(instances); vm++ {
		sum += virtualMachineInstance(instances[vm]).SumCPU()
	}

	return sum
}
