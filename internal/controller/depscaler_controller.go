/*
Copyright 2024.

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

package controller

import (
	"context"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	depscalev1 "github.com/yrs147/test-operator/api/v1"
)

// DepScalerReconciler reconciles a DepScaler object
type DepScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=depscale.yrs.scaler,resources=depscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=depscale.yrs.scaler,resources=depscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=depscale.yrs.scaler,resources=depscalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DepScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *DepScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.Log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	log.Info("Reconcile called")
	depscaler := &depscalev1.DepScaler{}
	err := r.Get(ctx, req.NamespacedName, depscaler)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	beginTime := depscaler.Spec.Begin
	endTime := depscaler.Spec.End
	replicas := depscaler.Spec.Replicas

	currentTime := time.Now().UTC().Hour()

	if currentTime >= beginTime && currentTime <= endTime {
		for _, deploy := range depscaler.Spec.Deployments {
			deployment := &v1.Deployment{}
			err := r.Get(ctx, client.ObjectKey{
				Namespace: deploy.Namespace,
				Name:      deploy.Name,
			}, deployment)
			if err != nil {
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			if deployment.Spec.Replicas != &replicas {
				deployment.Spec.Replicas = &replicas
				err := r.Update(ctx, deployment)
				if err != nil {
					depscaler.Status.Status = depscalev1.FAILED
					return ctrl.Result{}, err
				}
				depscaler.Status.Status = depscalev1.SUCCESS
				err = r.Status().Update(ctx, depscaler)
				if err != nil {
					return ctrl.Result{}, err
				}
			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DepScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&depscalev1.DepScaler{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				log.Log.WithValues("Namespace", e.Object.GetNamespace(), "Name", e.Object.GetName()).Info("DepScaler created")
				return true
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				log.Log.WithValues("Namespace", e.ObjectNew.GetNamespace(), "Name", e.ObjectNew.GetName()).Info("DepScaler updated")
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				log.Log.WithValues("Namespace", e.Object.GetNamespace(), "Name", e.Object.GetName()).Info("DepScaler deleted")
				return true
			},
		}).
		Complete(r)
}
