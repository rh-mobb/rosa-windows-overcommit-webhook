package vm

import (
	corev1 "kubevirt.io/api/core/v1"
)

type Store []corev1.VirtualMachineInstance

// Filter filters a Store object and returns a new store with only filtered virtual machine instances.  In the
// instance of this webhook, we only want virtual machine instances that are running a windows operating system.
func (store Store) Filter(filter *Filter) Store {
	return store
}

// SumCPU sums up the value of all CPUs in the store.
func (store Store) SumCPU() int {
	var sum int

	if len(store) == 0 {
		return 0
	}

	for vm := 0; vm < len(store); vm++ {
		// according to kubevirt docs, vcpu is determined by the value of sockets * cores * threads
		// see https://kubevirt.io/user-guide/compute/dedicated_cpu_resources/#requesting-dedicated-cpu-resources
		sockets := int(store[vm].Spec.Domain.CPU.Sockets)
		if sockets == 0 {
			sockets = 1
		}

		cores := int(store[vm].Spec.Domain.CPU.Cores)
		if cores == 0 {
			cores = 1
		}

		threads := int(store[vm].Spec.Domain.CPU.Threads)
		if threads == 0 {
			threads = 1
		}

		sum += sockets * cores * threads
	}

	return sum
}
