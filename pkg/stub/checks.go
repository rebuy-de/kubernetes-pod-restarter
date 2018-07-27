package stub

import (
	"github.com/sirupsen/logrus"

	core "k8s.io/api/core/v1"
)

func isAvailable(minAvailable int32, maxUnavailable int32, podList *core.PodList) bool {
	unavailable := int32(0)
	available := int32(0)

	for _, pod := range podList.Items {
		ready := true

		for _, container := range pod.Status.InitContainerStatuses {
			if !container.Ready {
				ready = false
			}
		}

		for _, container := range pod.Status.ContainerStatuses {
			if !container.Ready {
				ready = false
			}
		}

		if !ready {
			unavailable++
		} else {
			available++
		}
	}

	logger := logrus.WithFields(logrus.Fields{
		"MinAvailable":   minAvailable,
		"MaxUnavailable": maxUnavailable,
		"Unavailable":    unavailable,
	})

	if unavailable >= maxUnavailable {
		logger.Info("Too much Pods unready.")
		return false
	}

	if minAvailable >= available {
		logger.Info("Not enough Pods ready.")
		return false
	}

	logger.Debug("Enough Pods available.")
	return true

}
