package nodes

// import (
// 	corev1 "k8s.io/api/core/v1"
// )

// type Store []corev1.Node

// // Filter filters a Store object and returns a new store with only filtered nodes.
// // Filter returns the filtered nodes given a filter client.
// func (store Store) Filter(filter *Filter) Store {
// 	filtered := Store{}

// 	for node := 0; node < len(store); node++ {
// 		// continue if we have no filter key
// 		if store[node].GetLabels()[filter.LabelKey] == "" {
// 			continue
// 		}

// 		// store the node if the filter matches
// 		for value := 0; value < len(filter.LabelValues); value++ {
// 			if filter.LabelValues[value] == store[node].GetLabels()[filter.LabelKey] {
// 				filtered = append(filtered, store[node])
// 			}
// 		}
// 	}

// 	return filtered
// }

// // SumCPU sums up the value of all CPUs in the store.
// func (store Store) SumCPU() int {
// 	var sum int

// 	for node := 0; node < len(store); node++ {
// 		sum += int(store[node].Status.Capacity.Cpu().Value())
// 	}

// 	return sum
// }
