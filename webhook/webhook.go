package webhook

import (
	"context"
	"fmt"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kubevirt.io/client-go/kubecli"

	"github.com/scottd018/rosa-windows-overcommit-webhook/resources"
)

// webhook represents a webhook object.
type webhook struct {
	Context    context.Context
	KubeClient *kubernetes.Clientset
	VirtClient kubecli.KubevirtClient
	NodeFilter resources.NodeFilter
}

// NewWebhook returns a new instance of a webhook object.
func NewWebhook() (*webhook, error) {
	// create the kubernetes client alongside the virtualization client
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config; %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client; %w", err)
	}

	virtClient, err := kubecli.GetKubevirtClientFromRESTConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubevirt client; %w", err)
	}

	// create and run the webhook
	return &webhook{
		Context:    context.Background(),
		KubeClient: kubeClient,
		VirtClient: virtClient,
		NodeFilter: resources.NewNodeFilter(os.Getenv(resources.EnvLabelKey), os.Getenv(resources.EnvLabelValues)),
	}, nil
}

// Validate runs the validation logic for the webhook.
func (wh *webhook) Validate(w http.ResponseWriter, r *http.Request) {
	// create the operation object
	op, err := NewOperation(w, r)
	if err != nil {
		op.respond(err.Error(), true)
	}
	op.log("received validation request")

	op.log(fmt.Sprintf("DEBUG: %+v", op.object))

	// return immediately if we do not need validation
	if !op.object.NeedsValidation() {
		op.respond("skipping validation", true)
		return
	}

	// get the requested capacity from the request
	requested := op.object.SumCPU()

	// get the virtual machine instance list from the cluster and the current used capacity
	vmInstanceList, err := wh.getFilteredVirtualMachineInstances()
	if err != nil {
		op.respond(err.Error(), true)
		return
	}
	used := vmInstanceList.SumCPU()

	// return if we found an instance in the cluster matching this name
	// TODO: this likely needs to be handled differently for an UPDATE request
	for i := 0; i < len(vmInstanceList); i++ {
		if vmInstanceList[i].GetName() == op.object.GetName() && vmInstanceList[i].GetNamespace() == op.object.GetNamespace() {
			op.respond("skipping validation", true)
			return
		}
	}

	// get the node list from the cluster and the total capacity
	nodeList, err := wh.getFilteredNodes()
	if err != nil {
		op.respond(err.Error(), true)
		return
	}
	total := nodeList.SumCPU()

	available := total - used

	op.log(fmt.Sprintf(
		"capacity: total=[%d], requested=[%d], used=[%d], available=[%d]",
		total,
		requested,
		used,
		available,
	))

	// ensure the requested capacity would not exceed the available capacity
	if requested > available {
		msg := fmt.Sprintf("requested capacity: [%d], exceeds available capacity: [%d]; currently used [%d]",
			requested,
			available,
			used,
		)
		op.response.allowed = false
		op.respond(msg, true)

		return
	}

	op.respond("request success", false)
}

const statusOkMessage = `{"msg": "server is healthy"}`

// HealthZ implements a simple health check that returns a 200 ok response.
func (wh *webhook) HealthZ(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := fmt.Fprint(w, statusOkMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// getFilteredNodes returns a list of filtered nodes that exist in the cluster.
func (wh *webhook) getFilteredNodes() (resources.Nodes, error) {
	nodeList, err := wh.KubeClient.CoreV1().Nodes().List(wh.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes; %v", err)
	}

	var nodes resources.Nodes = nodeList.Items

	return nodes.Filter(wh.NodeFilter), nil
}

// getFilteredVirtualMachineInstances returns a list of filtered virtual machine instances that exist in the cluster.
// we need to gather both virtual machines and virtual machine instances in the case that an instance is not yet
// created from a virtual machine object.  then we can merge the two together.
func (wh *webhook) getFilteredVirtualMachineInstances() (resources.VirtualMachineInstances, error) {
	vmInstancesAll, err := wh.VirtClient.VirtualMachineInstance("").List(wh.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list virtual machine instances; %v", err)
	}

	// convert our list to our internal resource and filter
	instancesFiltered := resources.VirtualMachineInstances(vmInstancesAll.Items).Filter(
		&resources.VirtualMachineInstancesFilter{},
	)

	vmsAll, err := wh.VirtClient.VirtualMachine("").List(wh.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list virtual machines; %v", err)
	}

	// convert our list to our internal resource and filter
	vmsFiltered := resources.VirtualMachines(vmsAll.Items).Filter(
		&resources.VirtualMachinesFilter{},
	)

	filtered := append(instancesFiltered, vmsFiltered...)

	// return only instances with unique names and namespaces.  this is to avoid a situation where we have a
	// vm instance created by a vm, but also accounts for someone trying to bypass the overcommit by creating a
	// vm instance directly
	return filtered.Unique(), nil
}
