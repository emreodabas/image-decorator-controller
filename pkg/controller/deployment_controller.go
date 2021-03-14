package controller

import (
	"context"
	"fmt"
	"github.com/emreodabas/image-decorator-controller/pkg/containerimage"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
	"time"
)

type ReconcileDeployment struct {
	Client            client.Client
	RequeueDuration   time.Duration
	IgnoredNamespaces []string
	BackupRegistry    string
}

// !! Kubebuilder will read this lines and generate related resources !!
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch

// Implement reconcile.Reconciler so the controller can reconcile objects
var (
	_      reconcile.Reconciler = &ReconcileDeployment{}
	logger                      = log.Log.WithName("ReconcileDeployment")
)

func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

	if r.isIgnoredNamespace(request.NamespacedName.Namespace) {
		logger.Info("This is ignored", request.NamespacedName.Name, request.Namespace)
		return r.successMessage()
	}

	// Fetch the Deployment from the cache
	deployment := &appsv1.Deployment{}
	err := r.Client.Get(ctx, request.NamespacedName, deployment)
	if errors.IsNotFound(err) {
		logger.Error(nil, "Could not find deployment")
		return r.requeueMessage(err)
	}

	if err != nil {
		return r.requeueMessage(fmt.Errorf("could not fetch deployment: %+v", err))
	}
	// Print the Deployment
	containers := deployment.Spec.Template.Spec.Containers
	for i, container := range containers {
		logger.Info("Reconciling deployment", "containerimage name", container.Name)
		if !strings.HasPrefix(container.Image, r.BackupRegistry) {
			imagePath, err := containerimage.CloneImage(container.Image, r.BackupRegistry)
			if err != nil {
				return r.requeueMessage(err)
			}
			containers[i].Image = imagePath
		}
	}
	deployment.Spec.Template.Spec.Containers = containers
	//re assign updated containers to deployment
	err = r.Client.Update(ctx, deployment)
	if err != nil {
		return r.requeueMessage(fmt.Errorf("could not write Deployment: %+v", err))
	}

	return r.successMessage()

}

func (r *ReconcileDeployment) isIgnoredNamespace(name string) bool {

	for _, ns := range r.IgnoredNamespaces {
		if ns == name {
			return true
		}
	}
	return false
}

func (r *ReconcileDeployment) successMessage() (reconcile.Result, error) {
	return reconcile.Result{Requeue: false}, nil
}

func (r *ReconcileDeployment) requeueMessage(err error) (reconcile.Result, error) {
	return reconcile.Result{Requeue: true, RequeueAfter: r.RequeueDuration}, err
}
