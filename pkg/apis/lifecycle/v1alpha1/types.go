package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PodRestarterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []PodRestarter `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PodRestarter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PodRestarterSpec   `json:"spec"`
	Status            PodRestarterStatus `json:"status,omitempty"`
}

type PodRestarterSpec struct {
	// Selector is how the target will be selected.
	Selector *metav1.LabelSelector `json:"selector"`

	// CooldownPeriod is the minimal time between to restart actions.
	// +optional
	CooldownPeriod metav1.Duration `json:"cooldownPeriod,omitempty"`

	// MaxUnavailable is the maximum amount of Pods which are allowed to be unavailable among the selected pods.
	// +optional
	MaxUnavailable int32 `json:"maxUnavailable,omitempty"`

	// MaxUnavailable is the maximum amount of Pods which are allowed to be unavailable among the selected pods.
	// +optional
	MinAvailable int32 `json:"minAvailable,omitempty"`

	// RestartCriteria describes what Pods should get restarted.
	// +optional
	RestartCriteria PodRestarterCriteria `json:"restartCriteria,omitempty"`
}

type PodRestarterCriteria struct {
	// MaxAge desribes what age a Pod must have at least to get restarted.
	// +optional
	MaxAge *metav1.Duration `json:"maxAge,omitempty"`

	// MaxMemoryRequestRatio desribes what the ratio between memory usage and requests a Pod must have at least to get restarted.
	// +optional
	//MaxMemoryRequestRatio float32 `json:"maxMemoryRequestRatio,omitempty"`

	// MaxMemoryLimitRatio desribes what the ratio between memory usage and limits a Pod must have at least to get restarted.
	// +optional
	//MaxMemoryLimitRatio float32 `json:"maxMemoryLimitRatio,omitempty"`
}

type PodRestarterStatus struct {
	// +optional
	LastAction metav1.Time `json:"lastAction,omitempty"`
}
