package stub

import (
	"context"
	"fmt"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rebuy-de/kubernetes-pod-restarter/pkg/apis/lifecycle/v1alpha1"
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

		if !isAvailable(o.Spec.MinAvailable, o.Spec.MaxUnavailable, podList) {
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
