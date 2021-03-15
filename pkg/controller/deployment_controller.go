package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

type ReconcileDeployment struct {
	Deployment Reconciler
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// for non grouped deployments and pods
// +kubebuilder:rbac:groups=\ ,resources=deployments,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=\ ,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=\ ,resources=pods,verbs=get;list;watch;create;update;patch

// Implement reconcile.Reconciler so the controller can reconcile objects
var (
	_            reconcile.Reconciler = &ReconcileDeployment{}
	deploylogger                      = log.Log.WithName("ReconcileDeployment")
)

func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	reconciler := r.Deployment
	if reconciler.isIgnoredNamespace(request.NamespacedName.Namespace) {
		deploylogger.Info("This is ignored", request.NamespacedName.Name, request.Namespace)
		return reconciler.successMessage()
	}

	// Fetch the Deployment from the cache
	deployment := &appsv1.Deployment{}
	err := r.Deployment.Client.Get(ctx, request.NamespacedName, deployment)
	if err != nil {
		return reconciler.errorWithoutRequeue(fmt.Errorf("could not fetch deployment: %+v", err))
	}
	// Print the Deployment
	containers := deployment.Spec.Template.Spec.Containers
	for i, container := range containers {
		deploylogger.Info("Reconciling deployment", "containerimage name", container.Name)
		if !strings.HasPrefix(container.Image, r.Deployment.BackupRegistry.RepositoryPath) {
			imagePath, err := r.Deployment.BackupRegistry.CloneImage(container.Image)
			if err != nil {
				return reconciler.requeueMessage(err)
			}
			containers[i].Image = imagePath
		}
	}
	deployment.Spec.Template.Spec.Containers = containers
	//re assign updated containers to deployment
	err = r.Deployment.Client.Update(ctx, deployment)
	if err != nil {
		return reconciler.requeueMessage(fmt.Errorf("could not write Deployment: %+v", err))
	}

	return reconciler.successMessage()

}
