package vm

// import (
// 	"encoding/json"
// 	"fmt"

// 	admissionv1 "k8s.io/api/admission/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	kubevirtcorev1 "kubevirt.io/api/core/v1"
// )

// const (
// 	VirtualMachineInstanceType = "VirtualMachineInstance"
// )

// // virtualMachineInstance represents a VirtualMachineInstance object.
// type virtualMachineInstance struct{}

// // NewVirtualMachineInstance returns a new virtualMachineInstance object.
// func NewVirtualMachineInstance() *virtualMachineInstance {
// 	return &virtualMachineInstance{}
// }

// // Extract extracts a VirtualMachineInstance object from a VirtualMachine object.
// func (vmi *virtualMachineInstance) Extract(admissionRequest *admissionv1.AdmissionRequest) (*kubevirtcorev1.VirtualMachineInstance, error) {
// 	instance := &kubevirtcorev1.VirtualMachineInstance{}
// 	if err := json.Unmarshal(admissionRequest.Object.Raw, &instance); err != nil {
// 		return nil, fmt.Errorf("failed to decode virtual machine instance object; %w", err)
// 	}

// 	return instance, nil
// }

// // Type simply returns the type.
// func (vmi *virtualMachineInstance) Type() string {
// 	return VirtualMachineInstanceType
// }

// // VirtualMachineInstanceFromVirtualMachine returns a VirutalMachineInstance object given a VirtualMachine object.
// func VirtualMachineInstanceFromVirtualMachine(vm *kubevirtcorev1.VirtualMachine) *kubevirtcorev1.VirtualMachineInstance {
// 	return &kubevirtcorev1.VirtualMachineInstance{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       VirtualMachineInstanceType,
// 			APIVersion: fmt.Sprintf("%s/%s", vm.GroupVersionKind().Group, vm.APIVersion),
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      vm.Name,
// 			Namespace: vm.Namespace,
// 		},
// 		Spec: vm.Spec.Template.Spec,
// 	}
// }
