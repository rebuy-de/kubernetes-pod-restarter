package stub

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	lifecycle "github.com/rebuy-de/kubernetes-pod-restarter/pkg/apis/lifecycle/v1alpha1"
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

func TestNeedsCooldown(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	clock = clockwork.NewFakeClockAt(now)

	cases := []struct {
		name     string
		cooldown time.Duration
		offset   time.Duration
		want     bool
	}{
		{
			name:     "last_action_just_now",
			cooldown: 30 * time.Second,
			offset:   -5 * time.Second,
			want:     true,
		},
		{
			name:     "last_action_long_ago",
			cooldown: 30 * time.Second,
			offset:   -5 * time.Hour,
			want:     false,
		},
		{
			name:     "last_action_in_the_future",
			cooldown: 30 * time.Second,
			offset:   time.Minute,
			want:     true,
		},
		{
			name:     "last_action_just_now_without_cooldown",
			cooldown: 0,
			offset:   -5 * time.Second,
			want:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			restarter := &lifecycle.PodRestarter{
				Spec: lifecycle.PodRestarterSpec{
					CooldownPeriod: meta.Duration{Duration: tc.cooldown},
				},
				Status: lifecycle.PodRestarterStatus{
					LastAction: meta.NewTime(now.Add(tc.offset)),
				},
			}

			have := needsCooldown(restarter)
			if have != tc.want {
				t.Fail()
			}
		})
	}
}

func TestIsOldEnough(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	clock = clockwork.NewFakeClockAt(now)

	cases := []struct {
		name    string
		created time.Time
		maxAge  time.Duration
		want    bool
	}{
		{
			name:    "created_just_now",
			created: now.Add(-5 * time.Second),
			maxAge:  30 * time.Second,
			want:    false,
		},
		{
			name:    "created_in_the_future",
			created: now.Add(1 * time.Hour),
			maxAge:  30 * time.Second,
			want:    false,
		},
		{
			name:    "created_long_ago",
			created: now.Add(-1 * time.Hour),
			maxAge:  30 * time.Second,
			want:    true,
		},
		{
			name:    "created_just_now_without_maxAge",
			created: now.Add(-5 * time.Second),
			maxAge:  0,
			want:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			restarter := &lifecycle.PodRestarter{
				Spec: lifecycle.PodRestarterSpec{
					RestartCriteria: lifecycle.PodRestarterCriteria{
						MaxAge: &meta.Duration{Duration: tc.maxAge},
					},
				},
			}
			pod := &core.Pod{
				ObjectMeta: meta.ObjectMeta{
					CreationTimestamp: meta.NewTime(tc.created),
				},
			}

			have := isOldEnough(restarter, pod)
			if have != tc.want {
				t.Fail()
			}
		})
	}
}
