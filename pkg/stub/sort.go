package stub

import (
	core "k8s.io/api/core/v1"
)

type PodsByAge core.PodList

func (a PodsByAge) Len() int {
	return len(a.Items)
}

func (a PodsByAge) Swap(i, j int) {
	a.Items[i], a.Items[j] = a.Items[j], a.Items[i]
}

func (a PodsByAge) Less(i, j int) bool {
	return a.Items[i].ObjectMeta.CreationTimestamp.Time.Before(a.Items[j].ObjectMeta.CreationTimestamp.Time)
}
