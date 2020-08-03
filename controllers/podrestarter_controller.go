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

package controllers

import (
	"context"
	go_errors "errors"
	"fmt"
	"sort"

	"github.com/go-logr/logr"
	"github.com/labstack/gommon/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lifecyclev1alpha1 "github.com/rebuy-de/kubernetes-pod-restarter/api/v1alpha1"
)

// PodRestarterReconciler reconciles a PodRestarter object
type PodRestarterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=lifecycle.rebuy.com,resources=podrestarters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.rebuy.com,resources=podrestarters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *PodRestarterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("podrestarter", req.NamespacedName)

	podRestarter := &lifecyclev1alpha1.PodRestarter{}
	err := r.Get(ctx, req.NamespacedName, podRestarter)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("PodRestarter resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get PodRestarter")
		return ctrl.Result{}, err
	}
	selector, err := metav1.LabelSelectorAsMap(podRestarter.Spec.Selector)
	if err != nil {
		log.Error(err, "Can't format label selector as map")
		return ctrl.Result{}, err
	}

	if podRestarter.Spec.RestartCriteria.MaxAge == nil {
		err = go_errors.New("maxAge criteria not found")
		log.Error(err, ".spec.restartCriteria.maxAge is required")
		return ctrl.Result{}, err
	}

	if needsCooldown(log, podRestarter) {
		return ctrl.Result{}, nil
	}

	podList, err := r.listPods(ctx, podRestarter.ObjectMeta.Namespace, selector)
	if err != nil {
		log.Error(err, "Failed to list pods for selector")
		return ctrl.Result{}, err
	}

	if !isAvailable(log, podRestarter.Spec.MinAvailable, podRestarter.Spec.MaxUnavailable, podList) {
		err = go_errors.New("no pod available for deletion")
		log.Error(err, "Not enough ready pods available")
		return ctrl.Result{}, err
	}

	sort.Sort(PodsByAge(*podList))

	log.Info(fmt.Sprintf("Found %d matching Pods.", len(podList.Items)))

	// All other Pods are necessarily younger, because we sorted the list.
	pod := podList.Items[0]
	if isOldEnough(log, podRestarter, &pod) {
		podRestarter.Status.LastAction = metav1.Now()
		err := r.Update(ctx, podRestarter)
		if err != nil {
			return ctrl.Result{}, err
		}

		err = r.Delete(ctx, &pod)
		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	log.Info("Did not find any matching Pod.")

	return ctrl.Result{}, nil
}

func (r *PodRestarterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lifecyclev1alpha1.PodRestarter{}).
		Complete(r)
}

func (r *PodRestarterReconciler) listPods(ctx context.Context, namespace string, selector map[string]string) (*corev1.PodList, error) {
	podList := corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(selector),
	}
	if err := r.List(ctx, &podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods")
		return &corev1.PodList{}, err
	}

	return &podList, nil
}
