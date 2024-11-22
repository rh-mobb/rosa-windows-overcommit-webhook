package resources

import (
	corev1 "kubevirt.io/api/core/v1"
)

// VirtualMachineInstances is an object which represents a list of virtual machine instances.
type VirtualMachineInstances []corev1.VirtualMachineInstance

// VirtualMachineInstanceFilter represents a filter based on a set of key value inputs that are used to filter nodes.
type VirtualMachineInstanceFilter struct{}

// Filter filters a Store object and returns a new store with only filtered virtual machine instances.  In the
// instance of this webhook, we only want virtual machine instances that are running a windows operating system.
func (instances VirtualMachineInstances) Filter(filter *VirtualMachineInstanceFilter) VirtualMachineInstances {
	return instances
}

// SumCPU sums up the value of all CPUs in the store.
func (instances VirtualMachineInstances) SumCPU() int {
	var sum int

	if len(instances) == 0 {
		return 0
	}

	for vm := 0; vm < len(instances); vm++ {
		// continue the loop if we have a nil CPU configuration.  this simply defaults to 1 * 1 * 1 so we simply sum
		// the value (1) and continue the loop.
		if instances[vm].Spec.Domain.CPU == nil {
			sum += 1

			continue
		}

		// according to kubevirt docs, vcpu is determined by the value of sockets * cores * threads
		// see https://kubevirt.io/user-guide/compute/dedicated_cpu_resources/#requesting-dedicated-cpu-resources
		sockets := int(instances[vm].Spec.Domain.CPU.Sockets)
		if sockets == 0 {
			sockets = 1
		}

		cores := int(instances[vm].Spec.Domain.CPU.Cores)
		if cores == 0 {
			cores = 1
		}

		threads := int(instances[vm].Spec.Domain.CPU.Threads)
		if threads == 0 {
			threads = 1
		}

		sum += sockets * cores * threads
	}

	return sum
}
