package stub

import (
	"github.com/jonboulle/clockwork"
	"github.com/rebuy-de/kubernetes-pod-restarter/pkg/apis/lifecycle/v1alpha1"
	"github.com/sirupsen/logrus"

	core "k8s.io/api/core/v1"
)

var clock clockwork.Clock

func init() {
	clock = clockwork.NewRealClock()
}

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

func needsCooldown(o *v1alpha1.PodRestarter) bool {
	var (
		cooldown   = o.Spec.CooldownPeriod.Duration
		lastAction = o.Status.LastAction.Time
		nextAction = lastAction.Add(cooldown)
		now        = clock.Now()
	)

	if !lastAction.IsZero() && cooldown > 0 && nextAction.After(now) {
		logrus.WithFields(logrus.Fields{
			"NextAction": nextAction,
			"LastAction": lastAction,
			"Cooldown":   cooldown,
		}).Info("PodRestarter needs cooldown")
		return true
	}

	return false
}
