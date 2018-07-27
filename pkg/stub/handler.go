package stub

import (
	"context"
	"fmt"
	"time"

	"github.com/rebuy-de/kubernetes-pod-restarter/pkg/apis/lifecycle/v1alpha1"
	"github.com/sirupsen/logrus"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.PodRestarter:
		selector := meta.FormatLabelSelector(o.Spec.Selector)
		logrus.WithFields(logrus.Fields{
			"Selector":  selector,
			"Name":      o.ObjectMeta.Name,
			"Namespace": o.ObjectMeta.Namespace,
		}).Info("Syncing CRD")

		if o.Spec.RestartCriteria.MaxAge == nil {
			// TODO
			return fmt.Errorf(".spec.restartCriteria.maxAge is required")
		}

		if needsCooldown(o) {
			return nil
		}

		podList, err := listPods(o.ObjectMeta.Namespace, selector)
		if err != nil {
			return err
		}

		if insufficientAvailable(o.Spec.MinAvailable, o.Spec.MaxUnavailable, podList) {
			return nil
		}

		logrus.Debugf("Found %d matching Pods.", len(podList.Items))
		for _, pod := range podList.Items {
			created := pod.ObjectMeta.CreationTimestamp.Time
			age := time.Since(created)
			maxAge := o.Spec.RestartCriteria.MaxAge.Duration

			if age > maxAge {
				logrus.WithFields(logrus.Fields{
					"Name":   pod.ObjectMeta.Name,
					"Reason": "TooOld",
					"Age":    age,
				}).Info("Need to restart pod.")

				o.Status.LastAction = meta.Now()
				err := sdk.Update(o)
				if err != nil {
					return err
				}

				err = sdk.Delete(sdk.Object(&pod))
				if err != nil {
					return err
				}

				return nil
			}
		}

		logrus.Info("Did not find any matching Pod.")
	}

	return nil
}

func insufficientAvailable(minAvailable int32, maxUnavailable int32, podList *core.PodList) bool {
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

	if unavailable > maxUnavailable {
		logrus.WithFields(logrus.Fields{
			"MaxUnavailable": maxUnavailable,
			"Unavailable":    unavailable,
		}).Info("Too much Pods are unready.")
		return true
	}

	if minAvailable >= available {
		logrus.WithFields(logrus.Fields{
			"MinAvailable": minAvailable,
			"Available":    available,
		}).Info("Not enough Pods are ready.")
		return true
	}

	return false

}

func needsCooldown(o *v1alpha1.PodRestarter) bool {
	var (
		cooldown   = o.Spec.CooldownPeriod.Duration
		lastAction = o.Status.LastAction.Time
		nextAction = lastAction.Add(cooldown)
	)

	if !lastAction.IsZero() && cooldown > 0 && nextAction.After(time.Now()) {
		logrus.WithFields(logrus.Fields{
			"NextAction": nextAction,
			"LastAction": lastAction,
			"Cooldown":   cooldown,
		}).Info("PodRestarter needs cooldown")
		return true
	}

	return false
}

func listPods(namespace string, selector string) (*core.PodList, error) {
	podList := &core.PodList{
		TypeMeta: meta.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
	}

	options := sdk.WithListOptions(&meta.ListOptions{
		LabelSelector: selector,
	})

	err := sdk.List(namespace, podList, options)

	return podList, err
}
