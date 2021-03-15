package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

type ReconcileDaemonSet struct {
	DaemonSet Reconciler
}

// !! Kubebuilder will read this lines and generate related resources !!
// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=apps,resources=daemonsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=\ ,resources=daemonsets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=\ ,resources=daemonsets/status,verbs=get;update;patch

// Implement reconcile.Reconciler so the controller can reconcile objects
var (
	_            reconcile.Reconciler = &ReconcileDaemonSet{}
	daemonlogger                      = log.Log.WithName("ReconcileDaemonSet")
)

func (r *ReconcileDaemonSet) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	reconciler := r.DaemonSet
	if reconciler.isIgnoredNamespace(request.NamespacedName.Namespace) {
		daemonlogger.Info("This is ignored", request.NamespacedName.Name, request.Namespace)
		return reconciler.successMessage()
	}

	// Fetch the DaemonSet from the cache
	daemonSet := &appsv1.DaemonSet{}
	err := r.DaemonSet.Client.Get(ctx, request.NamespacedName, daemonSet)
	if err != nil {
		return reconciler.requeueMessage(fmt.Errorf("could not fetch DaemonSet: %+v", err))
	}
	// Print the DaemonSet
	containers := daemonSet.Spec.Template.Spec.Containers
	for i, container := range containers {
		daemonlogger.Info("Reconciling daemonSet", "containerimage name", container.Name)
		if !strings.HasPrefix(container.Image, r.DaemonSet.BackupRegistry.RepositoryPath) {
			imagePath, err := r.DaemonSet.BackupRegistry.CloneImage(container.Image)
			if err != nil {
				return reconciler.requeueMessage(err)
			}
			containers[i].Image = imagePath
		}
	}
	daemonSet.Spec.Template.Spec.Containers = containers
	//re assign updated containers to daemonSet
	err = r.DaemonSet.Client.Update(ctx, daemonSet)
	if err != nil {
		return reconciler.requeueMessage(fmt.Errorf("could not write DaemonSet: %+v", err))
	}
	return reconciler.successMessage()
}
