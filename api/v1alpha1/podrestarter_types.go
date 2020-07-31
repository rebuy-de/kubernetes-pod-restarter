/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodRestarterSpec defines the desired state of PodRestarter
type PodRestarterSpec struct {
	// Selector is how the target will be selected.
	Selector *metav1.LabelSelector `json:"selector"`

	// CooldownPeriod is the minimal time between to restart actions.
	// +optional
	CooldownPeriod metav1.Duration `json:"cooldownPeriod,omitempty"`

	// MaxUnavailable is the maximum amount of Pods which are allowed to be unavailable among the selected pods.
	// +optional
	MaxUnavailable int32 `json:"maxUnavailable"`

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

// PodRestarterStatus defines the observed state of PodRestarter
type PodRestarterStatus struct {
	// +optional
	LastAction metav1.Time `json:"lastAction,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PodRestarter is the Schema for the podrestarters API
type PodRestarter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodRestarterSpec   `json:"spec,omitempty"`
	Status PodRestarterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PodRestarterList contains a list of PodRestarter
type PodRestarterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodRestarter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodRestarter{}, &PodRestarterList{})
}
