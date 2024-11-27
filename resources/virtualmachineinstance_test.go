package resources

import (
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "kubevirt.io/api/core/v1"
)

func Test_virtualMachineInstance_NeedsValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		vmi  virtualMachineInstance
		want bool
	}{
		{
			name: "hasSysprepVolume: ensure resource with sysprep volume returns true",
			vmi: virtualMachineInstance{
				Spec: corev1.VirtualMachineInstanceSpec{
					Volumes: []corev1.Volume{
						{
							VolumeSource: corev1.VolumeSource{
								Sysprep: &corev1.SysprepSource{},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "hasSysprepVolume: ensure resource without sysprep volume returns false",
			vmi: virtualMachineInstance{
				Spec: corev1.VirtualMachineInstanceSpec{
					Volumes: []corev1.Volume{
						{
							VolumeSource: corev1.VolumeSource{
								HostDisk: &corev1.HostDisk{},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "hasWindowsDriverDiskVolume: ensure resource with a windows driver disk volume returns true",
			vmi: virtualMachineInstance{
				Spec: corev1.VirtualMachineInstanceSpec{
					Volumes: []corev1.Volume{
						{
							VolumeSource: corev1.VolumeSource{
								DataVolume: &corev1.DataVolumeSource{
									Name: "windows-drivers-disk",
								},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "hasWindowsDriverDiskVolume: ensure resource without a windows driver disk volume returns false",
			vmi: virtualMachineInstance{
				Spec: corev1.VirtualMachineInstanceSpec{
					Volumes: []corev1.Volume{
						{
							VolumeSource: corev1.VolumeSource{
								HostDisk: &corev1.HostDisk{},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "hasHyperV: ensure resource with hyperv settings configured returns true",
			vmi: virtualMachineInstance{
				Spec: corev1.VirtualMachineInstanceSpec{
					Domain: corev1.DomainSpec{
						Features: &corev1.Features{
							Hyperv: &corev1.FeatureHyperv{},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "hasHyperV: ensure resource without hyperv settings configured returns false",
			vmi: virtualMachineInstance{
				Spec: corev1.VirtualMachineInstanceSpec{},
			},
			want: false,
		},
		{
			name: "hasWindowsPreference: ensure resource with windows prefix annotation returns true",
			vmi: virtualMachineInstance{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kubevirt.io/cluster-preference-name": "windows.2k19",
					},
				},
			},
			want: true,
		},
		{
			name: "hasWindowsPreference: ensure resource without windows prefix annotation returns true",
			vmi: virtualMachineInstance{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kubevirt.io/cluster-preference-name": "rhel.9",
					},
				},
			},
			want: false,
		},
		{
			name: "hasWindowsPreference: ensure resource without annotation returns false",
			vmi: virtualMachineInstance{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"fake": "annotation",
					},
				},
			},
			want: false,
		},
		{
			name: "hasWindowsPreference: ensure resource with no annotations returns false",
			vmi:  virtualMachineInstance{},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.vmi.NeedsValidation()
			if got := result.NeedsValidation; got != tt.want {
				t.Errorf("virtualMachineInstance.NeedsValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}
