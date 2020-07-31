package controllers

import (
	"github.com/go-logr/logr"
	"github.com/jonboulle/clockwork"

	core "k8s.io/api/core/v1"

	lifecyclerebuycomv1alpha1 "github.com/rebuy-de/kubernetes-pod-restarter/api/v1alpha1"
)

var clock clockwork.Clock

func init() {
	clock = clockwork.NewRealClock()
}

func isAvailable(log logr.Logger, minAvailable int32, maxUnavailable int32, podList *core.PodList) bool {
	unavailable := int32(0)
	available := int32(0)
	log = log.WithValues("MinAvailable", minAvailable, "MaxUnavailable", maxUnavailable, "Unavailable", unavailable)

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

	if maxUnavailable > 0 && unavailable >= maxUnavailable {
		log.Info("Too many Pods unready.")
		return false
	}

	if minAvailable >= available {
		log.Info("Not enough Pods ready.")
		return false
	}

	return true
}

func needsCooldown(log logr.Logger, o *lifecyclerebuycomv1alpha1.PodRestarter) bool {
	var (
		cooldown   = o.Spec.CooldownPeriod.Duration
		lastAction = o.Status.LastAction.Time
		nextAction = lastAction.Add(cooldown)
		now        = clock.Now()
	)

	if !lastAction.IsZero() && cooldown > 0 && nextAction.After(now) {
		log = log.WithValues("NextAction", nextAction, "LastAction", lastAction, "Cooldown", cooldown)
		log.Info("PodRestarter needs cooldown")
		return true
	}

	return false
}

func isOldEnough(log logr.Logger, restarter *lifecyclerebuycomv1alpha1.PodRestarter, pod *core.Pod) bool {
	var (
		maxAge  = restarter.Spec.RestartCriteria.MaxAge.Duration
		created = pod.ObjectMeta.CreationTimestamp.Time
		age     = clock.Since(created)
	)

	log = log.WithValues("Name", pod.ObjectMeta.Name, "Reason", "TooOld", "Age", age, "MaxAge", maxAge, "Created", created)

	if age > maxAge {
		log.Info("Pod is old enough for restart.")
		return true
	}

	return false
}
