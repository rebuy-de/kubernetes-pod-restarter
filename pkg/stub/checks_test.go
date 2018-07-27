package stub

import (
	"testing"

	core "k8s.io/api/core/v1"
)

func TestIsAvailable(t *testing.T) {
	podList := &core.PodList{
		Items: []core.Pod{
			core.Pod{
				Status: core.PodStatus{
					InitContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
					},
					ContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
					},
				},
			},
			core.Pod{
				Status: core.PodStatus{
					InitContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
						core.ContainerStatus{
							Ready: true,
						},
					},
					ContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
					},
				},
			},
			core.Pod{
				Status: core.PodStatus{
					InitContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
					},
					ContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
						core.ContainerStatus{
							Ready: true,
						},
					},
				},
			},
			core.Pod{
				Status: core.PodStatus{
					InitContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: false,
						},
					},
					ContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
					},
				},
			},
			core.Pod{
				Status: core.PodStatus{
					InitContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
					},
					ContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: false,
						},
					},
				},
			},
			core.Pod{
				Status: core.PodStatus{
					InitContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
					},
					ContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
						core.ContainerStatus{
							Ready: false,
						},
					},
				},
			},
			core.Pod{
				Status: core.PodStatus{
					InitContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
						core.ContainerStatus{
							Ready: false,
						},
					},
					ContainerStatuses: []core.ContainerStatus{
						core.ContainerStatus{
							Ready: true,
						},
					},
				},
			},
		},
	}

	cases := []struct {
		name           string
		minAvailable   int32
		maxUnavailable int32
		want           bool
	}{
		{
			name:           "max_values",
			minAvailable:   0,
			maxUnavailable: 7, // 7 pods in total
			want:           true,
		},
		{
			name:           "lowest_max_unavailable",
			minAvailable:   0,
			maxUnavailable: 5, // 4 unavailable pods + 1 pod that gets restarted
			want:           true,
		},
		{
			name:           "highest_min_available",
			minAvailable:   2, // 3 available - 1 pod that gets restarted
			maxUnavailable: 7,
			want:           true,
		},
		{
			name:           "highest_min_available_and_lowest_max_unavailable",
			minAvailable:   2,
			maxUnavailable: 5,
			want:           true,
		},
		{
			name:           "lowest_max_unavailable_minus_one",
			minAvailable:   0,
			maxUnavailable: 4, // lowest_max_unavailable - 1
			want:           false,
		},
		{
			name:           "highest_min_available_plus_one",
			minAvailable:   3, // highest_min_available + 1
			maxUnavailable: 7,
			want:           false,
		},
	}

	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			have := isAvailable(tc.minAvailable, tc.maxUnavailable, podList)
			if have != tc.want {
				t.Fatalf("Case %d failed.", i)
			}
		})
	}
}
