package nodes

// import (
// 	"os"
// 	"strings"
// )

// const (
// 	EnvLabelKey    string = "WEBHOOK_NODE_LABEL_KEY"
// 	EnvLabelValues string = "WEBHOOK_NODE_LABEL_VALUES"

// 	DefaultLabelKey    string = "windows"
// 	DefaultLabelValues string = "true"
// )

// // Filter represents a filter based on a set of key value inputs that are used to filter nodes.
// type Filter struct {
// 	// LabelKey represents the node label that the filter is looking for.  It is derived from
// 	// the EnvLabelKey environment variables constant when creating a NodeFilter object from the helper function.
// 	LabelKey string

// 	// LabelValues represents the values associated with teh ImageLabelKey which
// 	// are equivalent to the AMI-images used to provision the nodes.  It is derived from the
// 	// EnvLabelValues environment variables constant when creating a NodeFilter object from the helper function.
// 	// The environment variable should be a comma-separated list and is created as such.
// 	LabelValues []string
// }

// // NewNodeFilter returns a new instance of a NodeFilter object with sane defaults.
// func NewNodeFilter() *Filter {
// 	labelKey, labelValuesString := os.Getenv(EnvLabelKey), os.Getenv(EnvLabelValues)

// 	// set the image label key
// 	if labelKey == "" {
// 		labelKey = DefaultLabelKey
// 	}

// 	// set the image label values
// 	if labelValuesString == "" {
// 		labelValuesString = DefaultLabelValues
// 	}
// 	labelValues := strings.Split(labelValuesString, ",")

// 	return &Filter{
// 		LabelKey:    labelKey,
// 		LabelValues: labelValues,
// 	}
// }
