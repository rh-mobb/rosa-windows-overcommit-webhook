package resources

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	EnvLabelKey    string = "WEBHOOK_NODE_LABEL_KEY"
	EnvLabelValues string = "WEBHOOK_NODE_LABEL_VALUES"

	DefaultLabelKey    string = "image_type"
	DefaultLabelValues string = "windows"
)

type Nodes []corev1.Node

// NodeFilter is an interface to represent a node filter.
type NodeFilter interface {
	LabelKey() string
	LabelValues() []string
}

// nodeFilter represents a filter based on a set of key value inputs that are used to filter nodes.
type nodeFilter struct {
	// labelKey represents the node label that the filter is looking for.  It is derived from
	// the EnvLabelKey environment variables constant when creating a NodeFilter object from the helper function.
	labelKey string

	// labelValues represents the values associated with teh ImageLabelKey which
	// are equivalent to the AMI-images used to provision the nodes.  It is derived from the
	// EnvLabelValues environment variables constant when creating a NodeFilter object from the helper function.
	// The environment variable should be a comma-separated list and is created as such.
	labelValues []string
}

// NewNodeFilter returns a new instance of a NodeFilter object with sane defaults.
func NewNodeFilter(labelKey string, labelValuesString string) *nodeFilter {
	// set the image label key and default if missing
	if labelKey == "" {
		labelKey = DefaultLabelKey
	}

	// set the image label values and default if missing
	if labelValuesString == "" {
		labelValuesString = DefaultLabelValues
	}

	return &nodeFilter{
		labelKey:    labelKey,
		labelValues: strings.Split(labelValuesString, ","),
	}
}

// Filter filters a list of nodes and returns a new list of nodes with only filtered nodes.
// Filter returns the filtered nodes given a filter client.
func (nodes Nodes) Filter(filter NodeFilter) Nodes {
	filtered := Nodes{}

	for node := 0; node < len(nodes); node++ {
		// continue if we have no filter key
		if nodes[node].GetLabels()[filter.LabelKey()] == "" {
			continue
		}

		// store the node if the filter matches
		for value := 0; value < len(filter.LabelValues()); value++ {
			if filter.LabelValues()[value] == nodes[node].GetLabels()[filter.LabelKey()] {
				filtered = append(filtered, nodes[node])
			}
		}
	}

	return filtered
}

// SumCPU sums up the value of all CPUs in the node list.
func (nodes Nodes) SumCPU() int {
	var sum int

	if len(nodes) == 0 {
		return sum
	}

	for node := 0; node < len(nodes); node++ {
		sum += int(nodes[node].Status.Capacity.Cpu().Value())
	}

	return sum
}

// LabelKey returns the label key for the node filter.  It is used to satisfy the NodeFilter interface.
func (nodeFilter *nodeFilter) LabelKey() string {
	return nodeFilter.labelKey
}

// LabelValues returns the label values for the node filter.  It is used to satisfy the NodeFilter interface.
func (nodeFilter *nodeFilter) LabelValues() []string {
	return nodeFilter.labelValues
}
